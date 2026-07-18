package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	weatherv1 "weather/v1"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// The address of the gRPC server that provides weather information.
// In a production setup, this would be a config value (k8s configmap, env var, etc.),
// and point to a load balanced endpoint of the weather service.
const serverAddr = "127.0.0.1:50051"

type GetWeatherParams struct {
	Location        string `json:"location" jsonschema:"city name, e.g. Copenhagen, DK"`
	FromDate        string `json:"from_date" jsonschema:"start date to query, YYYY-MM-DD"`
	ToDate          string `json:"to_date" jsonschema:"end date to query, YYYY-MM-DD"`
	TemperatureUnit string `json:"temperature_unit" jsonschema:"C or F"`
}

type WeatherInfo struct {
	Date          string `json:"date"`
	DaytimeTemp   string `json:"hitemp"`
	NighttimeTemp string `json:"lowtemp"`
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

func formatDate(d *date.Date) string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func parseTemperatureUnit(unit string) weatherv1.TemperatureUnit {
	if strings.EqualFold(unit, "C") {
		return weatherv1.TemperatureUnit_CELCIUS
	}
	return weatherv1.TemperatureUnit_FAHRENHEIT
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
		return nil, nil, err
	}
	defer conn.Close()

	client := weatherv1.NewWeatherServiceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	from_date, err := dateStrToDate(params.FromDate)
	if err != nil {
		return nil, nil, err
	}

	to_date, err := dateStrToDate(params.ToDate)
	if err != nil {
		return nil, nil, err
	}

	temperatureUnit := parseTemperatureUnit(params.TemperatureUnit)

	rpcreq := &weatherv1.GetWeatherRequest{
		Location: params.Location,
		DateRange: &weatherv1.DateRange{
			Begin: from_date,
			End:   to_date,
		},
		TemperatureUnit: temperatureUnit,
	}

	rpcres, err := client.GetWeather(ctx, rpcreq)
	if err != nil {
		log.Error().Err(err).Msg("error calling GetWeather")
		return nil, nil, err
	}

	days := make([]WeatherInfo, 0, len(rpcres.Response))
	for _, wi := range rpcres.Response {
		days = append(days, WeatherInfo{
			Date:          formatDate(wi.Date),
			DaytimeTemp:   fmt.Sprintf("%.0f", wi.GetHiTemperature()),
			NighttimeTemp: fmt.Sprintf("%.0f", wi.GetLowTemperature()),
			Conditions:    wi.Conditions,
		})
	}

	temperatureUnits := "F"
	if temperatureUnit == weatherv1.TemperatureUnit_CELCIUS {
		temperatureUnits = "C"
	}

	return nil, &GetWeatherResult{
		Location:         params.Location,
		TemperatureUnits: temperatureUnits,
		Days:             days,
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
