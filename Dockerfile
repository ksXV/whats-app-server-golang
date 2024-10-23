FROM golang:1.22.7

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY Makefile ./

RUN make build

CMD ["./bin/server"]
