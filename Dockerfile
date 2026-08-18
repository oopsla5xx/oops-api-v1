FROM golang:1.26.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w \
      -X github.com/oopsla5xx/oops-api-v1/internal/shared/version.Version=${VERSION} \
      -X github.com/oopsla5xx/oops-api-v1/internal/shared/version.Commit=${COMMIT} \
      -X github.com/oopsla5xx/oops-api-v1/internal/shared/version.BuildTime=${BUILD_TIME}" \
    -o /app/bin/api ./cmd/api

# ---

FROM alpine:3.20

RUN apk add --no-cache ca-certificates wget && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/bin/api ./api

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/api/v1/health || exit 1

ENTRYPOINT ["./api"]
