############################
# STEP 1 build executable binary
############################
# golang alpine 1.11.5
#FROM golang:1.11.4-alpine3.8 as builder
FROM golang:1.11.5-alpine3.8 as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk add --update --no-cache make bash git openssh-client ca-certificates build-base musl-dev curl wget tzdata \
    && update-ca-certificates \
    && go get -u github.com/mailru/easyjson/easyjson

ENV GO111MODULE=on

# Create appuser
RUN adduser -D -g '' appuser

WORKDIR /usr/src/app
COPY . .

# Fetch dependencies.

# Using go mod.
RUN go mod download

# Generate easyjson files
RUN go generate ./dto/*

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/exchange-rates-server

############################
# STEP 2 build a small image
############################
FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=builder /go/bin/exchange-rates-server /go/bin/exchange-rates-server

# Use an unprivileged user.
USER appuser

# Run binary.
ENTRYPOINT ["/go/bin/exchange-rates-server", "--base", "USD", "--", "GBP", "EUR", "CAD", "RUB", "CHF", "BTC", "ETH", "ETC", "ALF", "AGS", "AMC", "APEX", "ARCH", "ARI", "BCX", "BET", "BLK"]