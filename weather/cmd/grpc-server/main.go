package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"

	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc"
	weatherv1 "weather/v1"
)

func formatDate(d *date.Date) string {
	return fmt.Sprintf("%04d/%02d/%02d", d.Year, d.Month, d.Day)
}

func dateRangeDict(dr *weatherv1.DateRange) *zerolog.Event {
	return zerolog.Dict().
		Str("begin", formatDate(dr.Begin)).
		Str("end", formatDate(dr.End))
}

type weatherServer struct {
	weatherv1.UnimplementedWeatherServiceServer
}

func (s *weatherServer) GetWeather(ctx context.Context, req *weatherv1.GetWeatherRequest) (*weatherv1.GetWeatherResponse, error) {
	// TODO: your dummy logic here, keyed off req.Location and req.Range
	log.Info().
		Dict("dateRange", dateRangeDict(req.DateRange)).
		Str("location", req.Location).
		Msg("getWeather")
	info, err := GetDataForDateRange(
		req.Location,
		req.TemperatureUnit,
		req.DateRange,
	)
	if err != nil {
		return nil, err
	}
	return &weatherv1.GetWeatherResponse{Response: info}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal().Err(err).Msg("listen")
	}

	grpcServer := grpc.NewServer()
	weatherv1.RegisterWeatherServiceServer(grpcServer, &weatherServer{})

	log.Info().Msg("weather gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("serve")
	}
}
