package models

import (
	"github.com/doubletrey/crawlab-core/interfaces"
	"github.com/doubletrey/crawlab-core/utils/binders"
)

func GetModelColName(id interfaces.ModelId) (colName string) {
	return binders.NewColNameBinder(id).MustBindString()
}
