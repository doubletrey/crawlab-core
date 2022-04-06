package server

import (
	"encoding/json"
	"github.com/apex/log"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"github.com/doubletrey/crawlab-core/entity"
	"github.com/doubletrey/crawlab-core/errors"
	"github.com/doubletrey/crawlab-core/interfaces"
	"github.com/doubletrey/crawlab-core/models/service"
	"github.com/doubletrey/crawlab-core/node/config"
	"github.com/doubletrey/crawlab-core/task/stats"
	"github.com/doubletrey/crawlab-core/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"io"
	"strings"
)

type TaskServer struct {
	grpc.UnimplementedTaskServiceServer

	// dependencies
	modelSvc service.ModelService
	cfgSvc   interfaces.NodeConfigService
	statsSvc interfaces.TaskStatsService

	// internals
	server interfaces.GrpcServer
}

func (svr TaskServer) Subscribe(stream grpc.TaskService_SubscribeServer) (err error) {
	for {
		msg, err := stream.Recv()
		utils.LogDebug(msg.String())
		if err == io.EOF {
			return nil
		}
		if err != nil {
			if strings.HasSuffix(err.Error(), "context canceled") {
				return nil
			}
			trace.PrintError(err)
			continue
		}
		switch msg.Code {
		case grpc.StreamMessageCode_INSERT_DATA:
			err = svr.handleInsertData(msg)
		case grpc.StreamMessageCode_INSERT_LOGS:
			err = svr.handleInsertLogs(msg)
		default:
			err = errors.ErrorGrpcInvalidCode
			log.Errorf("invalid stream message code: %d", msg.Code)
			continue
		}
		if err != nil {
			log.Errorf("grpc error[%d]: %v", msg.Code, err)
		}
	}
}

func (svr TaskServer) handleInsertData(msg *grpc.StreamMessage) (err error) {
	data, err := svr.deserialize(msg)
	if err != nil {
		return err
	}
	var records []interface{}
	for _, d := range data.Records {
		res, ok := d["_tid"]
		if ok {
			switch res.(type) {
			case string:
				id, err := primitive.ObjectIDFromHex(res.(string))
				if err == nil {
					d["_tid"] = id
				}
			}
		}
		records = append(records, d)
	}
	return svr.statsSvc.InsertData(data.TaskId, records...)
}

func (svr TaskServer) handleInsertLogs(msg *grpc.StreamMessage) (err error) {
	data, err := svr.deserialize(msg)
	if err != nil {
		return err
	}
	return svr.statsSvc.InsertLogs(data.TaskId, data.Logs...)
}

func (svr TaskServer) deserialize(msg *grpc.StreamMessage) (data entity.StreamMessageTaskData, err error) {
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return data, trace.TraceError(err)
	}
	if data.TaskId.IsZero() {
		return data, trace.TraceError(errors.ErrorGrpcInvalidType)
	}
	return data, nil
}

func NewTaskServer(opts ...TaskServerOption) (res *TaskServer, err error) {
	// task server
	svr := &TaskServer{}

	// apply options
	for _, opt := range opts {
		opt(svr)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, err
	}
	if err := c.Provide(stats.ProvideGetTaskStatsService(svr.server.GetConfigPath())); err != nil {
		return nil, err
	}
	if err := c.Provide(config.ProvideConfigService(svr.server.GetConfigPath())); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
		statsSvc interfaces.TaskStatsService,
		cfgSvc interfaces.NodeConfigService,
	) {
		svr.modelSvc = modelSvc
		svr.statsSvc = statsSvc
		svr.cfgSvc = cfgSvc
	}); err != nil {
		return nil, err
	}

	return svr, nil
}

func ProvideTaskServer(server interfaces.GrpcServer, opts ...TaskServerOption) func() (res *TaskServer, err error) {
	return func() (*TaskServer, error) {
		opts = append(opts, WithServerTaskServerService(server))
		return NewTaskServer(opts...)
	}
}
