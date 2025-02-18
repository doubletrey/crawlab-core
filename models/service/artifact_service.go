package service

import (
	"github.com/doubletrey/crawlab-core/errors"
	"github.com/doubletrey/crawlab-core/interfaces"
	models2 "github.com/doubletrey/crawlab-core/models/models"
	"github.com/doubletrey/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeArtifact(d interface{}, err error) (res *models2.Artifact, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Artifact)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetArtifactById(id primitive.ObjectID) (res *models2.Artifact, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdArtifact).GetById(id)
	return convertTypeArtifact(d, err)
}

func (svc *Service) GetArtifact(query bson.M, opts *mongo.FindOptions) (res *models2.Artifact, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdArtifact).Get(query, opts)
	return convertTypeArtifact(d, err)
}

func (svc *Service) GetArtifactList(query bson.M, opts *mongo.FindOptions) (res []models2.Artifact, err error) {
	err = svc.getListSerializeTarget(interfaces.ModelIdArtifact, query, opts, &res)
	return res, err
}
