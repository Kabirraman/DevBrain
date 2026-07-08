FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o devbrain ./cmd/main.go

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=build /app/devbrain .

EXPOSE 8080
CMD ["./devbrain"]
