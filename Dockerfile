# Этап 1: Сборка бинарного файла (Build stage)
FROM golang:1.24.5-alpine AS builder

# Устанавливаем ca-certificates и git для предотвращения x509 ошибок TLS и успешного скачивания модулей
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Настройка GOPROXY для обхода TLS-ошибок сертификатов
ENV GOPROXY=https://golang.org,direct

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачивание зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка оптимизированного бинарного файла приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api/main.go

# Этап 2: Финальный образ
FROM alpine:latest

WORKDIR /app

# Копирование системных сертификатов из этапа сборки для безопасных внешних запросов
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Копирование миграций СУБД
COPY --from=builder /app/migrations ./migrations

# Копирование скомпилированного исполняемого файла
COPY --from=builder /app/api .

# Автоматическое создание директории под асинхронные логи
RUN mkdir -p data

# Экспонирование сетевого порта приложения
EXPOSE 8080

# Запуск приложения
CMD ["./api"]
