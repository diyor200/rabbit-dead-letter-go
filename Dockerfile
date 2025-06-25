FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN  GOOS=linux go build -o app main.go

FROM alpine:latest

COPY --from=builder  app .

CMD [ "./app" ]