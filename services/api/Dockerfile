FROM golang:1.25-bookworm

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        gcc \
        g++ \
        make \
        unzip \
        libaio1 \
    && rm -rf /var/lib/apt/lists/*


# Oracle Instant Client 19.32
COPY instantclient-basic-linux.x64-19.32.0.0.0dbru.zip /tmp/

RUN mkdir -p /opt/oracle \
    && unzip -q /tmp/instantclient-basic-linux.x64-19.32.0.0.0dbru.zip -d /opt/oracle \
    && echo "/opt/oracle/instantclient_19_32" > /etc/ld.so.conf.d/oracle-instantclient.conf \
    && ldconfig \
    && rm -f /tmp/instantclient-basic-linux.x64-19.32.0.0.0dbru.zip

ENV LD_LIBRARY_PATH=/opt/oracle/instantclient_19_32

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8080

CMD ["go", "run", "main.go"]
