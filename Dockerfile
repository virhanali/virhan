FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main .

ENV POSTGRES_HOST_AUTH_METHOD=md5

EXPOSE 80

CMD ["./main"]