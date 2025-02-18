package service

import (
	"github.com/doubletrey/crawlab-core/errors"
	"github.com/doubletrey/crawlab-core/interfaces"
	models2 "github.com/doubletrey/crawlab-core/models/models"
	"github.com/doubletrey/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeSchedule(d interface{}, err error) (res *models2.Schedule, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Schedule)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetScheduleById(id primitive.ObjectID) (res *models2.Schedule, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdSchedule).GetById(id)
	return convertTypeSchedule(d, err)
}

func (svc *Service) GetSchedule(query bson.M, opts *mongo.FindOptions) (res *models2.Schedule, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdSchedule).Get(query, opts)
	return convertTypeSchedule(d, err)
}

func (svc *Service) GetScheduleList(query bson.M, opts *mongo.FindOptions) (res []models2.Schedule, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdSchedule, query, opts, &res)
	return res, err
}
