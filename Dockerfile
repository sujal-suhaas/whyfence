FROM golang:1.23.0-bookworm

RUN apt-get update \
    && apt-get install -y --no-install-recommends git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN go test ./... -count=1 \
    && go build -trimpath -o /usr/local/bin/whyfence ./cmd/whyfence

CMD ["whyfence", "--help"]
