FROM golang:1.10.4-alpine3.7 AS builder
ENV GOPATH /go
ADD . $GOPATH/src/github.com/yipeeio/kubeiql
ADD https://storage.googleapis.com/kubernetes-release/release/v1.6.4/bin/linux/amd64/kubectl /usr/local/bin/kubectl
ADD https://github.com/kubernetes-sigs/aws-iam-authenticator/releases/download/0.4.0-alpha.1/aws-iam-authenticator_0.4.0-alpha.1_linux_amd64 /usr/local/bin/heptio-authenticator-aws
WORKDIR $GOPATH/src/github.com/yipeeio/kubeiql
RUN apk add --no-cache git \
                       curl \
                       ca-certificates && \
    chmod +x /usr/local/bin/kubectl && \
    sh -x gobuild.sh

FROM alpine:3.7
USER 496
WORKDIR /usr/local/bin/
COPY --from=builder --chown=496:496 /go/src/github.com/yipeeio/kubeiql/kubeiql ./
#ADD --chown=496:496 https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VSN:-v1.11.2}/bin/linux/amd64/kubectl /usr/local/bin/kubectl
#RUN chmod 777 /usr/local/bin/kubectl
EXPOSE 8128

CMD ["./kubeiql"]
