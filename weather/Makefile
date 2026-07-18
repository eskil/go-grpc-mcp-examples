.PHONY: genproto curl run

all: genproto bin/grpc-server bin/mcp-server

genproto: proto/weather.proto
	protoc --go_out=. --go_opt=module=weather --go-grpc_out=. --go-grpc_opt=module=weather -Iproto proto/weather.proto

run-grpc: all
	go run ./cmd/grpc-server

run-mcp: all
	go run ./cmd/mcp-server

bin/grpc-server: cmd/grpc-server/*.go
	go build -o bin/grpc-server ./cmd/grpc-server

bin/mcp-server: cmd/mcp-server/*.go
	go build -o bin/mcp-server ./cmd/mcp-server

curl:
	grpcurl -plaintext -proto proto/weather.proto -import-path proto \
		-d '{"location":"Copenhagen, DK","dateRange":{"begin":{"year":2026,"month":6,"day":22},"end":{"year":2026,"month":6,"day":28}},"temperatureUnit":"CELCIUS"}' \
		localhost:50051 weather.v1.WeatherService/GetWeather

clean:
	rm -f bin/grpc-server bin/mcp-server
	rm -rf v1/