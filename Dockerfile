FROM golang:1.13-alpine as builder

# Install our build tools

RUN apk add --update git make bash ca-certificates

WORKDIR /go/src/github.com/daspawnw/prometheus-aws-discovery
COPY . ./
RUN make bin/linux

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/daspawnw/prometheus-aws-discovery/bin/linux/prometheus-aws-discovery /prometheus-aws-discovery

ENTRYPOINT ["/prometheus-aws-discovery"]
