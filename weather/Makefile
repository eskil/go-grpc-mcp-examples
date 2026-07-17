.PHONY: genproto

all: genproto

genproto: proto/weather.proto
	protoc --go_out=. --go-grpc_out=. -Iproto proto/weather.proto
