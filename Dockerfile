# syntax=docker/dockerfile:1

FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod ./

COPY *.go ./

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /build_res

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY --from=build-stage /build_res /service

USER nonroot:nonroot

ENTRYPOINT ["/service"]