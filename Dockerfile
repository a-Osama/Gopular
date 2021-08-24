FROM golang:1.17-alpine AS build

WORKDIR /app
COPY . ./

RUN go mod tidy
RUN go build -o /gopular

FROM alpine:latest

WORKDIR /
COPY --from=build /gopular /gopular

ENTRYPOINT ["/gopular", "popular"]