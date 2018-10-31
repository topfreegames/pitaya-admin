package remote

import (
	"context"
	"encoding/json"

	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya-admin/protos"
	"github.com/topfreegames/pitaya/component"
)

// Doc is the documentation remote
type Doc struct {
	component.Base
}

// NewDoc returns an instance of doc remote
func NewDoc() *Doc {
	return &Doc{}
}

// Docs route returns auto documentation for Tennis
func (d *Doc) Docs(ctx context.Context, flag *protos.DocMsg) (*protos.Doc, error) {
	doc, err := pitaya.Documentation(flag.GetProtos)

	if err != nil {
		return nil, err
	}
	documentation, err := json.Marshal(doc)

	if err != nil {
		return nil, err
	}

	return &protos.Doc{Doc: string(documentation)}, nil
}
