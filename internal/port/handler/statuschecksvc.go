package handler

import (
	"context"

	topav1 "github.com/pranayhere/touringparty/gen/go/v1"
	"github.com/pranayhere/touringparty/pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StatusCheckService struct {
	topav1.UnimplementedStatusCheckServiceServer
}

func NewStatusCheckService() *StatusCheckService {
	return &StatusCheckService{}
}

func (s *StatusCheckService) LivenessCheck(ctx context.Context, _ *emptypb.Empty) (*topav1.StatusCheckResponse, error) {
	logger.Ctx(ctx).Debugw("topa liveness check", "status", "Healthy")
	return &topav1.StatusCheckResponse{
		ServingStatus: topav1.StatusCheckResponse_SERVING_STATUS_SERVING,
	}, nil
}
