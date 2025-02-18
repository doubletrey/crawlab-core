package client

import (
	"github.com/doubletrey/crawlab-core/errors"
	"github.com/doubletrey/crawlab-core/interfaces"
	"github.com/doubletrey/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiderServiceDelegate struct {
	interfaces.GrpcClientModelBaseService
}

func (svc *SpiderServiceDelegate) GetSpiderById(id primitive.ObjectID) (s interfaces.Spider, err error) {
	res, err := svc.GetById(id)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.Spider)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *SpiderServiceDelegate) GetSpider(query bson.M, opts *mongo.FindOptions) (s interfaces.Spider, err error) {
	res, err := svc.Get(query, opts)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.Spider)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *SpiderServiceDelegate) GetSpiderList(query bson.M, opts *mongo.FindOptions) (res []interfaces.Spider, err error) {
	list, err := svc.GetList(query, opts)
	if err != nil {
		return nil, err
	}
	for _, item := range list.Values() {
		s, ok := item.(interfaces.Spider)
		if !ok {
			return nil, errors.ErrorModelInvalidType
		}
		res = append(res, s)
	}
	return res, nil
}

func NewSpiderServiceDelegate(opts ...ModelBaseServiceDelegateOption) (svc2 interfaces.GrpcClientModelSpiderService, err error) {
	// apply options
	opts = append(opts, WithBaseServiceModelId(interfaces.ModelIdSpider))

	// base service
	baseSvc, err := NewBaseServiceDelegate(opts...)
	if err != nil {
		return nil, err
	}

	// service
	svc := &SpiderServiceDelegate{baseSvc}

	return svc, nil
}

func ProvideSpiderServiceDelegate(path string, opts ...ModelBaseServiceDelegateOption) func() (svc interfaces.GrpcClientModelSpiderService, err error) {
	if path != "" {
		opts = append(opts, WithBaseServiceConfigPath(path))
	}
	return func() (svc interfaces.GrpcClientModelSpiderService, err error) {
		return NewSpiderServiceDelegate(opts...)
	}
}
