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

// The parameters for the `get_weather` tool. The LLM will use the description
// and jsonschema tags to generate valid requests for the tool.
type GetWeatherParams struct {
	Location        string `json:"location" jsonschema:"city name, e.g. Copenhagen, DK"`
	FromDate        string `json:"from_date" jsonschema:"start date to query, YYYY-MM-DD"`
	ToDate          string `json:"to_date" jsonschema:"end date to query, YYYY-MM-DD"`
	TemperatureUnit string `json:"temperature_unit" jsonschema:"C or F"`
}

// The result of the `get_weather` tool. The LLM will use the description
// and jsonschema tags to understand the structure of the result.
type GetWeatherResult struct {
	Location         string        `json:"location" jsonschema:"city name, e.g. Copenhagen, DK"`
	TemperatureUnits string        `json:"temperature_units" jsonschema:"C or F"`
	Days             []WeatherInfo `json:"days" jsonschema:"weather info for each day in the date range"`
}

// WeatherInfo represents the weather information for a single day in the
// GetWeatherResult. The LLM will use the description and jsonschema tags to understand the structure of the result.
type WeatherInfo struct {
	Date          string `json:"date" jsonschema:"date of the weather forecast"`
	DaytimeTemp   string `json:"hitemp" jsonschema:"daytime temperature"`
	NighttimeTemp string `json:"lowtemp" jsonschema:"nighttime temperature"`
	Conditions    string `json:"conditions" jsonschema:"weather conditions"`
}

// dateStrToDate converts a date string in the format "YYYY-MM-DD" to a *date.Date.
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

// formatDate converts a *date.Date to a string in the format "YYYY-MM-DD".
func formatDate(d *date.Date) string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// parseTemperatureUnit converts a string representation of a temperature unit ("C" or "F") to the corresponding weatherv1.TemperatureUnit enum value.
func parseTemperatureUnit(unit string) weatherv1.TemperatureUnit {
	if strings.EqualFold(unit, "C") {
		return weatherv1.TemperatureUnit_CELCIUS
	}
	return weatherv1.TemperatureUnit_FAHRENHEIT
}

// getWeather is the implementation of the `get_weather` tool.
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
		return nil, nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}
	defer conn.Close()

	client := weatherv1.NewWeatherServiceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	from_date, err := dateStrToDate(params.FromDate)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse from_date: %w", err)
	}

	to_date, err := dateStrToDate(params.ToDate)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse to_date: %w", err)
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
		return nil, nil, fmt.Errorf("failed to call GetWeather: %w", err)
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

	// Register a specific `get_weather` tool with the server, which will be callable by clients.
	// The description and schema of it's `params` argument are used by the LLM to generate valid requests for the tool.
	// See the toplevel .mcp.json, that's where the LLM will find the tool
	// and how to start the server.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_weather",
		Description: "Get temperature and conditions for a location and date",
	}, getWeather)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal().Err(err).Msg("cannot run server")
	}
}
