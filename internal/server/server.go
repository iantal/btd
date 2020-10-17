package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/iantal/btd/internal/files"
	"github.com/iantal/btd/internal/service"
	protos "github.com/iantal/btd/protos/btd"
)

type BuildDetector struct {
	log hclog.Logger
	ds  *service.Detector
}

func NewBuildDetector(log hclog.Logger, basePath, rkHost string, store files.Storage) *BuildDetector {
	return &BuildDetector{log, service.NewDetector(log, basePath, rkHost, store)}
}

func (b *BuildDetector) GetBuildTools(ctx context.Context, rr *protos.BuildToolRequest) (*protos.BuildToolResponse, error) {
	b.log.Info("Handle request for project", "projectID", rr.GetProjectID())

	breakdownResult, err := b.ds.Detect(rr.GetProjectID(), rr.GetCommitHash())
	if err != nil {
		return nil, err
	}

	return &protos.BuildToolResponse{BuildTools: breakdownResult}, nil
}
