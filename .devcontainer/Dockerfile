FROM mcr.microsoft.com/devcontainers/base:bookworm

# protoc (protobuf-compiler) for generating gRPC code from .proto files
RUN apt-get update && apt-get install -y --no-install-recommends \
        protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*
