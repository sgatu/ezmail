ARG GOVERS=1.24

FROM golang:${GOVERS}

WORKDIR /opt/src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN rm -rf ./dist/executor
RUN mkdir ./dist/executor
RUN go build -v -o ./dist/executor/consumer ./cmd/consumer/main.go
RUN chmod +x ./dist/executor/consumer

CMD ["./dist/executor/consumer"]

