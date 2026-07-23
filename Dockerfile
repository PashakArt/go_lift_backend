# 1. Этап сборки (Build)
FROM golang:1.22-alpine AS builder

# Аргумент, какой именно сервис собираем
ARG SERVICE_PATH

WORKDIR /app

# Копируем весь Go workspace (включая api, services, go.work)
COPY . .

# Собираем бинарник нужного сервиса
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./${SERVICE_PATH}/cmd/main.go

# 2. Этап запуска (Production)
FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/app .

EXPOSE 8080 50051

CMD ["./app"]