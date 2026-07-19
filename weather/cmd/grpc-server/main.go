package main

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc"

	weatherv1 "weather/v1"
)

// The address of the gRPC server that provides weather information.
// In a production setup, this would be a config value (k8s configmap, env var, etc.),
// and point to a load balanced endpoint of the weather service.
const serverAddr = ":50051"

type weatherServer struct {
	weatherv1.UnimplementedWeatherServiceServer
}

// dateStrToDate converts a date string in the format "YYYY-MM-DD" to a *date.Date.
func formatDate(d *date.Date) string {
	return fmt.Sprintf("%04d/%02d/%02d", d.Year, d.Month, d.Day)
}

// dateRangeDict returns a zerolog.Dict with the begin and end dates of the given DateRange.
func dateRangeDict(dr *weatherv1.DateRange) *zerolog.Event {
	return zerolog.Dict().
		Str("begin", formatDate(dr.Begin)).
		Str("end", formatDate(dr.End))
}

func (s *weatherServer) GetWeather(ctx context.Context, req *weatherv1.GetWeatherRequest) (*weatherv1.GetWeatherResponse, error) {
	log.Info().
		Dict("date_range", dateRangeDict(req.DateRange)).
		Str("location", req.Location).
		Str("temp_unit", req.TemperatureUnit.String()).
		Msg("getWeather")

	// The actual logic of the GetWeather RPC would go here. For this sample,
	// we will just return data from a static sample dataset. In a real application,
	// this would likely involve querying a database or calling an external API.
	info, err := GetDataForDateRange(
		req.Location,
		req.TemperatureUnit,
		req.DateRange,
	)
	if err != nil {
		return nil, fmt.Errorf("GetDataForDateRange failed: %w", err)
	}

	return &weatherv1.GetWeatherResponse{Response: info}, nil
}

func main() {
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("listen")
	}

	grpcServer := grpc.NewServer()
	weatherv1.RegisterWeatherServiceServer(grpcServer, &weatherServer{})

	log.Info().Msg("weather gRPC server listening on " + serverAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("serve")
	}
}
