FROM golang:1.22.5-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o myapp .

FROM alpine:latest

WORKDIR /root/
COPY --from=build /app/myapp .

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/America/New_York /etc/localtime && \
    echo "America/New_York" > /etc/timezone

ENV TZ=America/New_York

CMD ["./myapp"]

