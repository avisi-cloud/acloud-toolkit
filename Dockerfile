FROM golang:1.24 as builder

WORKDIR /data/
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN make build

# ----
FROM registry.avisi.cloud/library/centos:8
COPY --from=builder /data/bin/acloud-toolkit /usr/local/bin/acloud-toolkit

USER root
# Install envsubst en kubectl
COPY files/kubernetes.repo /etc/yum.repos.d/
RUN yum install -y gettext kubectl && yum clean all

ENV HELM_VERSION="v3.1.2"
ENV HELM_SHA256="e6be589df85076108c33e12e60cfb85dcd82c5d756a6f6ebc8de0ee505c9fd4c"

# Install helm client
RUN curl -L -O https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz \
    && sha256sum helm-${HELM_VERSION}-linux-amd64.tar.gz \
    && echo "${HELM_SHA256} helm-${HELM_VERSION}-linux-amd64.tar.gz" | sha256sum -c \
    && tar -zxvf helm-${HELM_VERSION}-linux-amd64.tar.gz \
    && mv linux-amd64/helm /usr/local/bin/helm \
    && rm -rf helm-${HELM_VERSION}-linux-amd64.tar.gz linux-amd64
