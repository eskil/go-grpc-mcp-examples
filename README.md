![Static Badge](https://img.shields.io/badge/%F0%9F%A4%96%20AI%20content-30%25-red)
![Static Badge](https://img.shields.io/badge/%F0%9F%A7%A0%20Human%20content-70%25-red)
![Static Badge](https://img.shields.io/badge/ircv3-no-blue)

# gRPC MCP Example

## What?

This is a sample MCP backed by a gRPC service. This just demonstrates
the setup. It's backed by static scraped weather data so it does not
provide real functionality.

It's mostly handcoded as a learning exercise.

## Why?

MCPs are extremely useful for connecting agents to internal data (customer data/setup/history; production configuration etc) that isn't readily available through existing LLM tools.

It's common to have existing gRPC services to access relevant data. Fronting this with a MCP server that specifies tools to allow access allows tools like `claude` to access your data reusing existing gRPC services.

See also https://lethain.com/systems-mcp/ and other articles.

## How?

We define a dummy grpc server for weather data. It uses hardcoded data for three months in 2026 for Copenhagen, DK.

We then define a MCP tool named `get_weather` with a description of `temperature and conditions for a location and date`

```go
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_weather",
		Description: "Get temperature and conditions for a location and date",
	}, getWeather)
```

This basically tells the LLM how it can get relevant info.

`getWeather` then just uses a grpc client function;

```go
func getWeather(
	ctx context.Context,
	mcpreq *mcp.CallToolRequest,
	params GetWeatherParams,
) (
	*mcp.CallToolResult,
	any,
	error,
) {
...
```

Where the `params` struct tags inform the LLM how to use it;

```go
type GetWeatherParams struct {
	Location        string `json:"location" jsonschema:"city name, e.g. Copenhagen, DK"`
	FromDate        string `json:"from_date" jsonschema:"start date to query, YYYY-MM-DD"`
    ...
}
```

The response structure also uses tags and jsonschema to describe the result.

```go
type WeatherInfo struct {
	Date          string `json:"date" jsonschema:"date of the weather forecast"`
	Conditions    string `json:"conditions" jsonschema:"weather conditions"`
	DaytimeTemp   string `json:"hitemp" jsonschema:"daytime temperature"`
    ...
}

type GetWeatherResult struct {
	Location         string        `json:"location" jsonschema:"city name, e.g. Copenhagen, DK"`
	Days             []WeatherInfo `json:"days" jsonschema:"weather info for each day in the date range"`
}
```

See the implementation of `getWeather` and the structures for the full example.

See also https://www.redpanda.com/blog/turn-grpc-api-into-mcp-server for a detailed example.

## Setup & run

### Devcontainer

The project includes a devcontainer suitable for either VSCode/Cursor/devcontainer-cli.

* Run docker desktop
* Open this folder in VSCode/Cursor and follow it's steps
* Or run `devcontainer up --id-label name=go-grpc-mcp-dev --workspace-folder .`

_(Odd that it doesn't just work as `devcontainer up` and defaults to looking for the `.devcontainer` and picks the name from that)_.

### Local (bare-metal) MacOS

Install go & protobuf.

```shell
brew install go protobuf protoc-gen-go protoc-gen-go-grpc grpcurl
```

### Build && run

On bare-metal or in a devcontainer terminal;

```shell
# Build everything
make

# Successful output
make -C weather
make[1]: Entering directory '/workspaces/go-grpc-mcp-examples/weather'
protoc --go_out=. --go_opt=module=weather --go-grpc_out=. --go-grpc_opt=module=weather -Iproto proto/weather.proto
go build -o bin/grpc-server ./cmd/grpc-server
go: downloading google.golang.org/genproto v0.0.0-20260715232425-e75dac1f907d
...
go build -o bin/mcp-server ./cmd/mcp-server
go: downloading github.com/modelcontextprotocol/go-sdk v1.6.1
...
make[1]: Leaving directory '/workspaces/go-grpc-mcp-examples/weather'
```

```shell
# Run grpc server
cd weather
make run-grpc`

# Successful output
{"level":"info","time":"2026-07-18T19:01:42Z","message":"weather gRPC server listening on :50051"}```
```

The devcontainer should now be forewarding the port.

### Test

```shell
# Sample grpcurl
cd weather
make curl

# Successful output
grpcurl -plaintext -proto proto/weather.proto -import-path proto \
		-d '{"location":"Copenhagen, DK","dateRange":{"begin":{"year":2026,"month":6,"day":22},"end":{"year":2026,"month":6,"day":28}},"temperatureUnit":"CELCIUS"}' \
		localhost:50051 weather.v1.WeatherService/GetWeather
{
  "response": [
    {
      "date": {
        "year": 2026,
        ...
      "hiTemperature": 27.2,
      "lowTemperature": 22.2,
      "conditions": "partly cloudy"
    }
  ]
}
```

### Using

Run a `claude` session. It should pick up the `.mcp.json`. Run `/mcp` to list the `weather` tool and wether it's connected.

```
     Project MCPs (<source/path/to>/go-grpc-mcp-examples/.mcp.json)                    ❯ weather · ✔ connected · 1 tool
```

```
❯ Which month of may, june and july had the warmest night in Copenhagen in 2026?

  Called weather

June had the warmest night — the low on the night of June 28, 2026 hit 22°C, the highest overnight low across all three months (compared to a peak of 15°C in May and 19°C in July).
```

```
❯ Give me the average day and night time temperature for the last week of June in Copenhagen, in celcius but fahrenheit for night time.

Using June 24–30 (the last full week of June 2026) in Copenhagen:

- Average daytime high: 26.0°C
- Average nighttime low: 65.2°F (18.4°C)
```

```
❯ Describe the weather of June 2026 in Copenhagen.

June 2026 in Copenhagen ran from a cool, unsettled start to a hot finish:

- First half (June 1–17): Cool and rain-prone — highs mostly 14–20°C, lows 10–16°C. Rain fell on 8 of these 17 days (notably June 4–5, 7, 9, 11–14), with only one sunny day (June 10).
- Second half (June 18–30): A steady warming trend into a heatwave. Highs climbed from 21°C to a monthly peak of 32°C on June 27, staying mostly sunny/partly cloudy with no rain at all after June 17. Nights warmed too, peaking at a very mild 22°C low on June 28 — the warmest night of the month.

Monthly stats:
- Average high: 20.0°C (range: 14–32°C)
- Average low: 14.7°C (range: 10–22°C)
- Conditions mix: 12 days partly cloudy, 8 rainy, 6 sunny, 4 cloudy

Overall: a damp, mild first two weeks giving way to a dry, increasingly hot back half culminating in a late-month heatwave (June 24–29).
```