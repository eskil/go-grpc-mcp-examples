package main

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc"
	weatherv1 "weather/v1"
)

type GetWeatherParams struct {
	Location string `json:"location" jsonschema:"city name, e.g. Copenhagen, DK"`
	Date     string `json:"date" jsonschema:"date to query, YYYY-MM-DD"`
}

type WeatherInfo struct {
	Date          string `json:"date"`
	DaytimeTemp   string `json:"hitemp"`
	NighttimeTemp string `json:"hitemp"`
	Conditions    string `json:"conditions"`
}

type GetWeatherResult struct {
	Location         string        `json:"location"`
	TemperatureUnits string        `json:"temperature_units"`
	Days             []WeatherInfo `json:"days"`
}

func dateStrToDate(dateStr string) (*date.Date, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	return &date.Date{
		Year:  int32(t.Year()),
		Month: int32(t.Month()),
		Day:   int32(t.Day()),
	}, nil
}

func getWeather(
	ctx context.Context,
	mcpreq *mcp.CallToolRequest,
	params GetWeatherParams,
) (
	*mcp.CallToolResult,
	any,
	error,
) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to server")
	}
	defer conn.Close()

	client := weatherv1.NewWeatherClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rpcreq := &weatherv1.GetWeatherRequest{
		Location: params.Location,
		DateRange: &weatherv1.DateRange{
			Begin: dateStrToDate(params.Date),
			End:   dateStrToDate(params.Date),
		},
		Units: weatherv1.TemperatureUnit.FAHRENHEIT,
	}

	rpcres := client.GetWeather(ctx, rpcreq)
	if err != nil {
		log.Fatal().Err(err).Msg("calling GetWeather")
		return nil, GetWeatherResult{}, err
	}
	// TODO: dial your gRPC server, call GetWeather, translate the response
	// into a *mcp.CallToolResult (text content, typically)

	wi := rcpres.response[0]
	return nil, &GetWeatherResult{
		Location:         mcpreq.Location,
		TemperatureUnits: "F",
		Days: []WeatherInfo{
			{
				Date:          wi.Date,
				DaytimeTemp:   wi.HiTemperature,
				NighttimeTemp: wi.LowTemperature,
				Conditions:    wi.Conditions,
			},
		},
	}, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{Name: "weather", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_weather",
		Description: "Get temperature and conditions for a location and date",
	}, getWeather)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal().Err(err).Msg("cannot run server")
	}
}
