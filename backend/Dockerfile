FROM golang:1.24-alpine

WORKDIR /app

COPY main.go .

RUN go mod init backend && go build -o backend main.go

EXPOSE 9001
EXPOSE 9002

CMD ["./backend"]
