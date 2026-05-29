package grpchandler

import (
	"context"
	"count-emission-service/internal/domain"
	"count-emission-service/internal/model/emission"
	pb "count-emission-service/proto/generated"
	"count-emission-service/proto/dailytotal"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type EmissionGRPCServer struct {
	pb.UnimplementedEmissionServer
	EmissioonUseCase  domain.EmissionUseCase
	PreferenceUseCase domain.PreferenceUseCase
}

func NewEmissionGRPCServer(EmissionUseCase domain.EmissionUseCase, PreferenceUseCase domain.PreferenceUseCase) *EmissionGRPCServer {
	return &EmissionGRPCServer{
		EmissioonUseCase:  EmissionUseCase,
		PreferenceUseCase: PreferenceUseCase,
	}
}

func userIDFromContext(ctx context.Context) (int32, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("user-id")) == 0 {
		return 0, status.Error(codes.Unauthenticated, "user-id not found in metadata")
	}
	id, err := strconv.ParseInt(md.Get("user-id")[0], 10, 32)
	if err != nil {
		return 0, status.Error(codes.InvalidArgument, "invalid user-id in metadata")
	}
	return int32(id), nil
}

func (s *EmissionGRPCServer) CreateUserEmission(ctx context.Context, req *pb.EmissionBody) (*pb.PostRespon, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	payload := &emission.EmissionBody{
		UserId:      userID,
		VehicleType: req.VehicleType,
		FuelType:    req.FuelType,
		DistanceKm:  req.DistanceKm,
	}
	if err := s.EmissioonUseCase.CreateUserEmission(ctx, payload); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create emission: %v", err)
	}
	return &pb.PostRespon{Message: "Emission has been created"}, nil
}

func (s *EmissionGRPCServer) GetUserDailyEmission(ctx context.Context, req *pb.Empty) (*pb.UserDailyEmission, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userEmission, err := s.EmissioonUseCase.GetUserDailyEmission(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get daily emission: %v", err)
	}
	return &pb.UserDailyEmission{
		UserId:             userEmission.UserId,
		Date:               userEmission.Date,
		TotalEmissionKgCo2: userEmission.TotalEmissionKgCo2,
	}, nil
}

func (s *EmissionGRPCServer) GetUserMonthlyEmission(ctx context.Context, req *pb.Empty) (*pb.UserMonthlyEmission, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userMonthlyEmission, err := s.EmissioonUseCase.GetUserMonthlyEmission(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get monthly emission: %v", err)
	}

	var dailyEmissions []*pb.UserDailyEmission
	for _, d := range userMonthlyEmission.DailyEmissions {
		dailyEmissions = append(dailyEmissions, &pb.UserDailyEmission{
			UserId:             d.UserId,
			Date:               d.Date,
			TotalEmissionKgCo2: d.TotalEmissionKgCo2,
		})
	}
	return &pb.UserMonthlyEmission{
		UserId:                    userMonthlyEmission.UserId,
		DailyEmissions:            dailyEmissions,
		TotalEmissionMonthlyKgCo2: userMonthlyEmission.TotalEmissionMonthlyKgCo2,
	}, nil
}

func (s *EmissionGRPCServer) GetUserYearlyEmission(ctx context.Context, req *pb.Empty) (*pb.UserYearlyEmission, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userYearlyEmission, err := s.EmissioonUseCase.GetUserYearlyEmission(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get yearly emission: %v", err)
	}

	var monthlyEmissions []*pb.UserMonthlyEmissionDetail
	for _, m := range userYearlyEmission.MonthlyEmissions {
		monthlyEmissions = append(monthlyEmissions, &pb.UserMonthlyEmissionDetail{
			Month:              m.Month,
			TotalEmissionKgCo2: m.TotalEmissionKgCo2,
		})
	}
	return &pb.UserYearlyEmission{
		UserId:                   userYearlyEmission.UserId,
		MonthlyEmissions:         monthlyEmissions,
		TotalEmissionYearlyKgCo2: userYearlyEmission.TotalEmissionYearlyKgCo2,
	}, nil
}

func (s *EmissionGRPCServer) GetUserPreferences(ctx context.Context, req *pb.Empty) (*pb.UserPreferences, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	pref, err := s.PreferenceUseCase.GetUserPreferences(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get preferences: %v", err)
	}

	resp := &pb.UserPreferences{
		UserId:      pref.UserId,
		CountryCode: pref.CountryCode,
	}
	if pref.CustomDailyLimitKgCo2 != nil {
		resp.CustomDailyLimitKgCo2 = *pref.CustomDailyLimitKgCo2
	}
	return resp, nil
}

func (s *EmissionGRPCServer) GetDailyTotal(ctx context.Context, req *dailytotal.DailyTotalRequest) (*dailytotal.DailyTotalResponse, error) {
	total, count, err := s.EmissioonUseCase.GetDailyTotal(ctx, req.UserId, req.Date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetDailyTotal: %v", err)
	}
	return &dailytotal.DailyTotalResponse{
		DailyTotalKg: float32(total),
		TripCount:    count,
	}, nil
}

func (s *EmissionGRPCServer) SetUserPreferences(ctx context.Context, req *pb.SetUserPreferencesBody) (*pb.UserPreferences, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var customLimit *float64
	if req.CustomDailyLimitKgCo2 > 0 {
		v := req.CustomDailyLimitKgCo2
		customLimit = &v
	}

	pref, err := s.PreferenceUseCase.SetUserPreferences(ctx, userID, req.CountryCode, customLimit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set preferences: %v", err)
	}

	resp := &pb.UserPreferences{
		UserId:      pref.UserId,
		CountryCode: pref.CountryCode,
	}
	if pref.CustomDailyLimitKgCo2 != nil {
		resp.CustomDailyLimitKgCo2 = *pref.CustomDailyLimitKgCo2
	}
	return resp, nil
}
