FROM golang:1.22.5-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o myapp .

FROM alpine:latest

WORKDIR /root/
COPY --from=build /app/myapp .

CMD ["./myapp"]

