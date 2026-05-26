package grpchandler

import (
	"context"
	"count-emission-service/internal/domain"
	"count-emission-service/internal/model/emission"
	pb "count-emission-service/proto/generated"
)

type EmissionGRPCServer struct {
	pb.UnimplementedEmissionServer
	EmissioonUseCase domain.EmissionUseCase
}

func NewEmissionGRPCServer(EmissionUseCase domain.EmissionUseCase) *EmissionGRPCServer {
	return &EmissionGRPCServer{EmissioonUseCase: EmissionUseCase}
}

func (s *EmissionGRPCServer) CreateUserEmission(ctx context.Context, req *pb.EmissionBody) (*pb.PostRespon, error) {
	var payload emission.EmissionBody
	payload.UserId = req.UserId
	payload.VehicleType = req.VehicleType
	payload.FuelType = req.FuelType
	payload.DistanceKm = req.DistanceKm
	err := s.EmissioonUseCase.CreateUserEmission(ctx, &payload)
	if err != nil {
		return &pb.PostRespon{
			Message: err.Error(),
		}, nil
	}
	return &pb.PostRespon{
		Message: "Emission has been created",
	}, nil
}

func (s *EmissionGRPCServer) GetUserDailyEmission(ctx context.Context, req *pb.Empty) (*pb.UserDailyEmission, error) {
	userId := int32(1)
	userEmission, err := s.EmissioonUseCase.GetUserDailyEmission(ctx, userId)
	if err != nil {
		return nil, err
	}
	respon := &pb.UserDailyEmission{
		UserId:             userEmission.UserId,
		Date:               userEmission.Date,
		TotalEmissionKgCo2: userEmission.TotalEmissionKgCo2,
	}
	return respon, nil
}

func (s *EmissionGRPCServer) GetUserMonthlyEmission(ctx context.Context, req *pb.Empty) (*pb.UserMonthlyEmission, error) {
	userId := int32(1)
	userMonthlyEmission, err := s.EmissioonUseCase.GetUserMonthlyEmission(ctx, userId)
	if err != nil {
		return nil, err
	}
	var dailyEmissions []*pb.UserDailyEmission

	for _, dailyEmission := range userMonthlyEmission.DailyEmissions {
		dailyEmissions = append(dailyEmissions, &pb.UserDailyEmission{
			UserId:             dailyEmission.UserId,
			Date:               dailyEmission.Date,
			TotalEmissionKgCo2: dailyEmission.TotalEmissionKgCo2,
		})
	}
	respon := &pb.UserMonthlyEmission{
		UserId:                    userMonthlyEmission.UserId,
		DailyEmissions:            dailyEmissions,
		TotalEmissionMonthlyKgCo2: userMonthlyEmission.TotalEmissionMonthlyKgCo2,
	}
	return respon, nil
}

func (s *EmissionGRPCServer) GetUserYearlyEmission(ctx context.Context, req *pb.Empty) (*pb.UserYearlyEmission, error) {
	userId := int32(1)
	userYearlyEmission, err := s.EmissioonUseCase.GetUserYearlyEmission(ctx, userId)
	if err != nil {
		return nil, err
	}
	var monthlyEmissions []*pb.UserMonthlyEmissionDetail

	for _, monthlyEmission := range userYearlyEmission.MonthlyEmissions {
		monthlyEmissions = append(monthlyEmissions, &pb.UserMonthlyEmissionDetail{
			Month:              monthlyEmission.Month,
			TotalEmissionKgCo2: monthlyEmission.TotalEmissionKgCo2,
		})
	}
	respon := &pb.UserYearlyEmission{
		UserId:                   userYearlyEmission.UserId,
		MonthlyEmissions:         monthlyEmissions,
		TotalEmissionYearlyKgCo2: userYearlyEmission.TotalEmissionYearlyKgCo2,
	}

	return respon, nil
}
