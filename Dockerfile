FROM golang:1.23.0-bookworm

RUN apt-get update \
    && apt-get install -y --no-install-recommends git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd

RUN go test ./...

CMD ["go", "test", "./..."]
