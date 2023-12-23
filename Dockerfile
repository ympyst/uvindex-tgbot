# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /uvindexbot

EXPOSE 3000

CMD ["/uvindexbot"]