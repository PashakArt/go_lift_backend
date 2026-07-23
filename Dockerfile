# 1. Этап сборки (Build) - берем актуальный Go alpine
FROM golang:alpine AS builder

ARG SERVICE_PATH

WORKDIR /app

# Включаем автоматическое скачивание нужной версии Go, если в go.mod прописана новее
ENV GOTOOLCHAIN=auto

# Копируем весь workspace
COPY . .

# Переходим в папку конкретного сервиса
WORKDIR /app/${SERVICE_PATH}

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app ./cmd/main.go

# 2. Этап запуска (Production)
FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/app .

EXPOSE 8080 50051

CMD ["./app"]