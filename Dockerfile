# 1. Этап сборки (Build)
FROM golang:1.22-alpine AS builder

ARG SERVICE_PATH

WORKDIR /app

# Копируем весь workspace
COPY . .

# Переходим в папку конкретного сервиса для сборки
WORKDIR /app/${SERVICE_PATH}

# Собираем бинарник с учетом go.work из родительской папки
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app ./cmd/main.go

# 2. Этап запуска (Production)
FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/app .

EXPOSE 8080 50051

CMD ["./app"]