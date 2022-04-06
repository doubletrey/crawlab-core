package interfaces

import (
	"github.com/doubletrey/crawlab-db/generic"
)

type ResultService interface {
	Insert(records ...interface{}) (err error)
	List(query generic.ListQuery, opts *generic.ListOptions) (results []Result, err error)
	Count(query generic.ListQuery) (n int, err error)
}
