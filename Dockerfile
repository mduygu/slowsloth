FROM golang:1.22.3 AS build-stage

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /slowsloth

FROM alpine:latest AS build-release-stage
WORKDIR /
COPY --from=build-stage /slowsloth /slowsloth
ENTRYPOINT ["/slowsloth"]
