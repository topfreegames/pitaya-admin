package remote

import (
	"context"
	"errors"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/protos"
	"github.com/topfreegames/pitaya/component"
)

// Proto is the protobuf remote
type Proto struct {
	component.Base
}

// NewProto returns an instance of proto remote
func NewProto() *Proto {
	return &Proto{}
}

// Proto returns the descriptor of a given proto
func (p *Proto) Proto(ctx context.Context, message *protos.ProtoName) (*protos.ProtoDescriptor, error) {
	descriptor, err := pitaya.Descriptor(message.GetName())

	if err != nil {
		return nil, errors.New("failed to get proto descriptor")
	}

	return &protos.ProtoDescriptor{
		Desc: descriptor,
	}, nil
}
