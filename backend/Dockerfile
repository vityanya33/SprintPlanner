FROM golang:1.24.2-alpine

WORKDIR /app

# Копируем зависимости сначала для кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем ВСЕ исходные файлы
COPY . .

# Собираем приложение
RUN go build -o server .

EXPOSE 3000

CMD ["./server"]