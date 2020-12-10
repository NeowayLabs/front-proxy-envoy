FROM golang:1.15.6-alpine

WORKDIR /app

COPY ./service.go .

RUN CGO_ENABLED=0 go build -o main

EXPOSE 8080

CMD ["./main"]