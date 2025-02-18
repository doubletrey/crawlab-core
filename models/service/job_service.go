package service

import (
	"github.com/doubletrey/crawlab-core/errors"
	"github.com/doubletrey/crawlab-core/interfaces"
	models2 "github.com/doubletrey/crawlab-core/models/models"
	"github.com/doubletrey/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeJob(d interface{}, err error) (res *models2.Job, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Job)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetJobById(id primitive.ObjectID) (res *models2.Job, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdJob).GetById(id)
	return convertTypeJob(d, err)
}

func (svc *Service) GetJob(query bson.M, opts *mongo.FindOptions) (res *models2.Job, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdJob).Get(query, opts)
	return convertTypeJob(d, err)
}

func (svc *Service) GetJobList(query bson.M, opts *mongo.FindOptions) (res []models2.Job, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdJob, query, opts, &res)
	return res, err
}
