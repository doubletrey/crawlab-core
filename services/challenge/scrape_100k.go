package challenge

import (
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/globalsign/mgo/bson"
)

type Scrape100kService struct {
	UserId bson.ObjectId
}

func (s *Scrape100kService) Check() (bool, error) {
	query := bson.M{
		"user_id": s.UserId,
		"result_count": bson.M{
			"$gte": 100000,
		},
	}
	list, err := model.GetTaskList(query, 0, 1, "-_id")
	if err != nil {
		return false, err
	}
	return len(list) > 0, nil
}
