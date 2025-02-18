package service

import (
	"github.com/doubletrey/crawlab-core/errors"
	"github.com/doubletrey/crawlab-core/interfaces"
	models2 "github.com/doubletrey/crawlab-core/models/models"
	"github.com/doubletrey/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeTaskStat(d interface{}, err error) (res *models2.TaskStat, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.TaskStat)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetTaskStatById(id primitive.ObjectID) (res *models2.TaskStat, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdTaskStat).GetById(id)
	return convertTypeTaskStat(d, err)
}

func (svc *Service) GetTaskStat(query bson.M, opts *mongo.FindOptions) (res *models2.TaskStat, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdTaskStat).Get(query, opts)
	return convertTypeTaskStat(d, err)
}

func (svc *Service) GetTaskStatList(query bson.M, opts *mongo.FindOptions) (res []models2.TaskStat, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdTaskStat, query, opts, &res)
	return res, err
}
