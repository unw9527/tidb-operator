FROM alpine:3.10

ENV KUBECTL_VERSION=v1.24.17
ENV HELM_VERSION=v3.11.0

RUN apk update && apk add --no-cache ca-certificates curl git openssl bash mysql-client
RUN curl https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
    -o /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl && \
    curl https://get.helm.sh/helm-v${HELM_VERSION}-linux-amd64.tar.gz \
    -o helm-${HELM_VERSION}-linux-amd64.tar.gz && \
    tar -zxvf helm-${HELM_VERSION}-linux-amd64.tar.gz && \
    mv linux-amd64/helm /usr/local/bin/helm && \
    rm -rf linux-amd64 && \
    rm helm-${HELM_VERSION}-linux-amd64.tar.gz

ADD bin/blockwriter /usr/local/bin/
ADD bin/stability-test /usr/local/bin/stability-test
