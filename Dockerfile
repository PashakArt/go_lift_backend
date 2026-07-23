FROM golang:alpine AS builder

ARG SERVICE_PATH

WORKDIR /app

ENV GOTOOLCHAIN=auto

COPY . .

WORKDIR /app/${SERVICE_PATH}

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/app .

EXPOSE 8080 50051

CMD ["./app"]