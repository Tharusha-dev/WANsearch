FROM golang:1.22.1-alpine

ENV GIN_MODE release

WORKDIR /app

RUN apk add --no-cache gcc g++ git openssh-client

COPY go.mod go.sum ./


RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .

EXPOSE 8080

CMD ["./main"]
