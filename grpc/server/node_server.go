package server

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-grpc"
	"github.com/doubletrey/crawlab-core/constants"
	"github.com/doubletrey/crawlab-core/entity"
	"github.com/doubletrey/crawlab-core/errors"
	"github.com/doubletrey/crawlab-core/interfaces"
	"github.com/doubletrey/crawlab-core/models/delegate"
	"github.com/doubletrey/crawlab-core/models/models"
	"github.com/doubletrey/crawlab-core/models/service"
	"github.com/doubletrey/crawlab-core/node/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
)

type NodeServer struct {
	grpc.UnimplementedNodeServiceServer

	// dependencies
	modelSvc service.ModelService
	cfgSvc   interfaces.NodeConfigService

	// internals
	server interfaces.GrpcServer
}

// Register from handler/worker to master
func (svr NodeServer) Register(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	// unmarshall data
	var nodeInfo entity.NodeInfo
	if req.Data != nil {
		if err := json.Unmarshal(req.Data, &nodeInfo); err != nil {
			return HandleError(err)
		}

		if nodeInfo.IsMaster {
			// error: cannot register master node
			return HandleError(errors.ErrorGrpcNotAllowed)
		}
	}

	// node key
	var nodeKey string
	if req.NodeKey != "" {
		nodeKey = req.NodeKey
	} else {
		nodeKey = nodeInfo.Key
	}
	if nodeKey == "" {
		return HandleError(errors.ErrorModelMissingRequiredData)
	}

	// find in db
	node, err := svr.modelSvc.GetNodeByKey(nodeKey, nil)
	if err == nil {
		if node.IsMaster {
			// error: cannot register master node
			return HandleError(errors.ErrorGrpcNotAllowed)
		} else {
			// register existing
			node.Status = constants.NodeStatusRegistered
			node.Active = true
			nodeD := delegate.NewModelNodeDelegate(node)
			if err := nodeD.Save(); err != nil {
				return HandleError(err)
			}
			var ok bool
			node, ok = nodeD.GetModel().(*models.Node)
			if !ok {
				return HandleError(errors.ErrorGrpcInvalidType)
			}
			log.Infof("[NodeServer] updated worker[%s] in db. id: %s", nodeKey, nodeD.GetModel().GetId().Hex())
		}
	} else if err == mongo.ErrNoDocuments {
		// register new
		node = &models.Node{
			Key:         nodeKey,
			Name:        nodeInfo.Name,
			Ip:          nodeInfo.Ip,
			Hostname:    nodeInfo.Hostname,
			Description: nodeInfo.Description,
			MaxRunners:  nodeInfo.MaxRunners,
			Status:      constants.NodeStatusRegistered,
			Active:      true,
			Enabled:     true,
		}
		if node.Name == "" {
			node.Name = nodeKey
		}
		nodeD := delegate.NewModelDelegate(node)
		if err := nodeD.Add(); err != nil {
			return HandleError(err)
		}
		var ok bool
		node, ok = nodeD.GetModel().(*models.Node)
		if !ok {
			return HandleError(errors.ErrorGrpcInvalidType)
		}
		log.Infof("[NodeServer] added worker[%s] in db. id: %s", nodeKey, nodeD.GetModel().GetId().Hex())
	} else {
		// error
		return HandleError(err)
	}

	log.Infof("[NodeServer] master registered worker[%s]", req.GetNodeKey())

	return HandleSuccessWithData(node)
}

// SendHeartbeat from worker to master
func (svr NodeServer) SendHeartbeat(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	// find in db
	node, err := svr.modelSvc.GetNodeByKey(req.NodeKey, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return HandleError(errors.ErrorNodeNotExists)
		}
		return HandleError(err)
	}

	// validate status
	if node.Status == constants.NodeStatusUnregistered {
		return HandleError(errors.ErrorNodeUnregistered)
	}

	// update status
	nodeD := delegate.NewModelNodeDelegate(node)
	if err := nodeD.UpdateStatusOnline(); err != nil {
		return HandleError(err)
	}

	return HandleSuccessWithData(node)
}

// Ping from worker to master
func (svr NodeServer) Ping(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return HandleSuccess()
}

func (svr NodeServer) Subscribe(request *grpc.Request, stream grpc.NodeService_SubscribeServer) (err error) {
	log.Infof("[NodeServer] master received subscribe request from node[%s]", request.NodeKey)

	// finished channel
	finished := make(chan bool)

	// set subscribe
	svr.server.SetSubscribe("node:"+request.NodeKey, &entity.GrpcSubscribe{
		Stream:   stream,
		Finished: finished,
	})
	ctx := stream.Context()

	log.Infof("[NodeServer] master subscribed node[%s]", request.NodeKey)

	// Keep this scope alive because once this scope exits - the stream is closed
	for {
		select {
		case <-finished:
			log.Infof("[NodeServer] closing stream for node[%s]", request.NodeKey)
			return nil
		case <-ctx.Done():
			log.Infof("[NodeServer] node[%s] has disconnected", request.NodeKey)
			return nil
		}
	}
}

func (svr NodeServer) Unsubscribe(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	sub, err := svr.server.GetSubscribe("node:" + req.NodeKey)
	if err != nil {
		return nil, errors.ErrorGrpcSubscribeNotExists
	}
	select {
	case sub.GetFinished() <- true:
		log.Infof("unsubscribed node[%s]", req.NodeKey)
	default:
		// Default case is to avoid blocking in case client has already unsubscribed
	}
	svr.server.DeleteSubscribe(req.NodeKey)
	return &grpc.Response{
		Code:    grpc.ResponseCode_OK,
		Message: "unsubscribed successfully",
	}, nil
}

func NewNodeServer(opts ...NodeServerOption) (res *NodeServer, err error) {
	// node server
	svr := &NodeServer{}

	// apply options
	for _, opt := range opts {
		opt(svr)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, err
	}
	if err := c.Provide(config.ProvideConfigService(svr.server.GetConfigPath())); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(modelSvc service.ModelService, cfgSvc interfaces.NodeConfigService) {
		svr.modelSvc = modelSvc
		svr.cfgSvc = cfgSvc
	}); err != nil {
		return nil, err
	}

	return svr, nil
}

func ProvideNodeServer(server interfaces.GrpcServer, opts ...NodeServerOption) func() (res *NodeServer, err error) {
	return func() (*NodeServer, error) {
		opts = append(opts, WithServerNodeServerService(server))
		return NewNodeServer(opts...)
	}
}
