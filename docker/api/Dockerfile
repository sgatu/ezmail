ARG GOVERS=1.24

FROM golang:${GOVERS}

WORKDIR /opt/src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN rm -rf ./dist/api
RUN mkdir ./dist/api
RUN go build -v -o ./dist/api/ezmail ./cmd/api/main.go
RUN chmod +x ./dist/api/ezmail

CMD ["./dist/api/ezmail"]

