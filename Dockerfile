FROM golang:1.22-alpine AS builder
USER root
WORKDIR /app

COPY . /app

RUN apk update && \
    apk add --no-cache git tzdata

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /build ./cmd/shortener/main.go

FROM alpine:latest
USER root
WORKDIR /app/build

COPY --from=builder /build ./main
COPY internal/app/migrations migrations

RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENTRYPOINT [ "./main"]
CMD ["-d="]