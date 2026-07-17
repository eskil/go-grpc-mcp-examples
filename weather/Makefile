.PHONY: genproto curl run

all: genproto

genproto: proto/weather.proto
	protoc --go_out=. --go_opt=module=weather --go-grpc_out=. --go-grpc_opt=module=weather -Iproto proto/weather.proto

run-grpc: all
	go run ./cmd/grpcserver

run-mcp: all
	go run ./cmd/mcpserver

curl:
	grpcurl -plaintext -proto proto/weather.proto -import-path proto \
		-d '{"location":"Copenhagen, DK","dateRange":{"begin":{"year":2026,"month":6,"day":22},"end":{"year":2026,"month":6,"day":28}},"units":"CELCIUS"}' \
		localhost:50051 weather.v1.WeatherService/GetWeather

