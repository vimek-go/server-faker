FROM golang:1.22 AS builder

RUN cat /etc/*-release
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY plugins/ plugins/
COPY main.go .

# Build server-faker
RUN go build -o server-faker main.go
RUN apt-get update

# Build the logger plugin
WORKDIR /app/plugins/logger
RUN go mod download
RUN go build -buildmode=plugin -o endpoint_logger.so endpoint_logger.go

# Build the protobuf plugin
WORKDIR /app/plugins/protobuf
RUN mkdir -p generated
RUN go mod download
RUN go build -buildmode=plugin -o protobuf_response.so protobuf_response.go


FROM debian:bookworm-slim
COPY examples/ ./examples/
# Copy the server-faker executable
COPY --from=builder /app/server-faker server-faker

# RUN chmod +x server-faker
RUN ls -lh
# Copy the plugins in ../plugins directory
RUN mkdir -p /plugins
COPY --from=builder /app/plugins/protobuf/protobuf_response.so /plugins/protobuf_response.so
COPY --from=builder /app/plugins/logger/endpoint_logger.so /plugins/endpoint_logger.so

# Command to run the executable
CMD ["./server-faker", "run", "--file=examples/static.json"]