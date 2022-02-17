FROM golang:1.17 AS builder
WORKDIR /opt
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR /opt/consumer
RUN go build
WORKDIR /opt/producer
RUN go build

FROM gcr.io/distroless/cc-debian11 AS runtime
COPY --from=builder /opt/consumer /usr/local/bin/consumer
COPY --from=builder /opt/producer /usr/local/bin/producer
