# Етап 1: збірка
FROM golang:latest AS builder

WORKDIR /app

# Копіюємо модулі та завантажуємо залежності
COPY go.mod go.sum ./
RUN go mod download

# Копіюємо весь код
COPY . .

# Будуємо бінарний файл
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Етап 2: деплой
FROM alpine:latest

WORKDIR /root/

# Копіюємо бінарний файл
COPY --from=builder /app/main .

# Копіюємо конфігураційний файл та файл .env
COPY configs/config.yml /root/configs/config.yml
COPY .env /root/.env

# Встановлюємо сертифікати (якщо потрібні для HTTPS-запитів)
RUN apk --no-cache add ca-certificates

# Виставляємо порт
EXPOSE 9999

# Запускаємо додаток
CMD ["./main"]
