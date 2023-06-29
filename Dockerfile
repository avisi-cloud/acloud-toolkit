FROM golang:1.20 as builder

WORKDIR /data/
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN make build

# ----
FROM registry.avisi.cloud/library/ubuntu:22.04
COPY --from=builder /data/bin/acloud-toolkit /usr/local/bin/acloud-toolkit

USER root
