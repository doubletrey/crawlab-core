package entity

import (
	"encoding/json"
	"github.com/crawlab-team/go-trace"
	"github.com/doubletrey/crawlab-core/interfaces"
)

type GrpcDelegateMessage struct {
	ModelId interfaces.ModelId             `json:"id"`
	Method  interfaces.ModelDelegateMethod `json:"m"`
	Data    []byte                         `json:"d"`
}

func (msg *GrpcDelegateMessage) GetModelId() interfaces.ModelId {
	return msg.ModelId
}

func (msg *GrpcDelegateMessage) GetMethod() interfaces.ModelDelegateMethod {
	return msg.Method
}

func (msg *GrpcDelegateMessage) GetData() []byte {
	return msg.Data
}

func (msg *GrpcDelegateMessage) ToBytes() (data []byte) {
	data, err := json.Marshal(*msg)
	if err != nil {
		_ = trace.TraceError(err)
		return data
	}
	return data
}
