package server

import (
	"context"

	"github.com/iantal/btd/internal/files"
	"github.com/iantal/btd/internal/service"
	"github.com/iantal/btd/internal/util"
	protos "github.com/iantal/btd/protos/btd"
	"github.com/sirupsen/logrus"
)

type BuildDetector struct {
	log *util.StandardLogger
	ds  *service.Detector
}

func NewBuildDetector(log *util.StandardLogger, basePath, rkHost string, store files.Storage) *BuildDetector {
	return &BuildDetector{log, service.NewDetector(log, basePath, rkHost, store)}
}

func (b *BuildDetector) GetBuildTools(ctx context.Context, rr *protos.BuildToolRequest) (*protos.BuildToolResponse, error) {
	b.log.WithFields(logrus.Fields{
		"projectID": rr.ProjectID,
		"commit": rr.CommitHash,
	}).Info("Handle request for project")

	breakdownResult, err := b.ds.Detect(rr.GetProjectID(), rr.GetCommitHash())
	if err != nil {
		return nil, err
	}

	return &protos.BuildToolResponse{BuildTools: breakdownResult}, nil
}
