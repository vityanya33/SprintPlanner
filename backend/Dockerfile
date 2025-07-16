FROM golang:1.24.2-alpine

#установка рабочей директории
WORKDIR /app
#копирование зависимостей голанга
COPY go.mod go.sum ./
#загрузка зависимостей
RUN go mod download
#копирование остального кода
COPY . .
#сбор бинарника(компиляция)
RUN go build -o server
#открываю порт
EXPOSE 3000
#запускаю сервер
CMD ["./server"]