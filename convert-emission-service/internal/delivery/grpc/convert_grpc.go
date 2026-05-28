package grpchandler

import (
	"context"
	"convert-emission-service/internal/domain"
	pb "convert-emission-service/proto/generated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConvertGRPCServer struct {
	pb.UnimplementedConvertServer
	usecase domain.ConversionUsecase
}

func NewConvertGRPCServer(uc domain.ConversionUsecase) *ConvertGRPCServer {
	return &ConvertGRPCServer{usecase: uc}
}

func (s *ConvertGRPCServer) ConvertToIDR(ctx context.Context, req *pb.ConvertToIDRRequest) (*pb.ConvertToIDRResponse, error) {
	if req.EmissionKgCo2 < 0 {
		return nil, status.Error(codes.InvalidArgument, "emission_kg_co2 must be non-negative")
	}

	pricePerTonUsd, exchangeRate, totalIdr, err := s.usecase.ConvertToIDR(ctx, req.EmissionKgCo2)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert emission: %v", err)
	}

	return &pb.ConvertToIDRResponse{
		PricePerTonUsd:    pricePerTonUsd,
		ExchangeRateUsdIdr: exchangeRate,
		TotalIdr:          totalIdr,
	}, nil
}
