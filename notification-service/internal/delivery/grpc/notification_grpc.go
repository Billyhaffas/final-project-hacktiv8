package grpchandler

import (
	"context"

	"notification-service/internal/domain"
	pb "notification-service/proto/generated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationGRPCServer struct {
	pb.UnimplementedNotificationServer
	usecase domain.NotificationUsecase
}

func NewNotificationGRPCServer(uc domain.NotificationUsecase) *NotificationGRPCServer {
	return &NotificationGRPCServer{usecase: uc}
}

func (s *NotificationGRPCServer) CheckDailyAlert(ctx context.Context, req *pb.DailyAlertRequest) (*pb.DailyAlertResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	isExceeded, dailyTotal, dailyLimit, thresholdSource, msg, err := s.usecase.CheckDailyAlert(ctx, req.UserId, req.Date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check daily alert: %v", err)
	}

	return &pb.DailyAlertResponse{
		IsExceeded:      isExceeded,
		DailyTotalKg:    float32(dailyTotal),
		DailyLimitKg:    float32(dailyLimit),
		ThresholdSource: thresholdSource,
		Message:         msg,
	}, nil
}
