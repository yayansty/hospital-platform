FROM golang:1.25-bookworm

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod ./

RUN go mod download

COPY . .

EXPOSE 8080

CMD ["go", "run", "main.go"]
