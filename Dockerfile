FROM --platform=$BUILDPLATFORM golang:1.19-alpine as builder

ARG TARGETOS
ARG TARGETARCH
ARG RELEASE_VERSION=development

# Install our build tools
RUN apk add --update git make bash ca-certificates

WORKDIR /go/src/github.com/daspawnw/prometheus-aws-discovery
COPY . ./
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-X 'main.Version=${RELEASE_VERSION}'" -o bin/prometheus-aws-discovery ./cmd/prometheus-aws-discovery/...

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/daspawnw/prometheus-aws-discovery/bin/prometheus-aws-discovery /prometheus-aws-discovery

ENTRYPOINT ["/prometheus-aws-discovery"]
