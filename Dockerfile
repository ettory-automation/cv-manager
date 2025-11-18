# =====================
# STAGE 1 — BUILD
# =====================
FROM golang:1.25.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o cv-manager server.go


# =====================
# STAGE 2 — RUNTIME
# =====================
FROM alpine:3.22.2

SHELL ["/bin/sh", "-c"]

WORKDIR /app

RUN apk add --no-cache curl ca-certificates

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /app/cv-manager /app/cv-manager

RUN chmod 500 /app/cv-manager && chown app:app /app/cv-manager

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD curl -fsS http://localhost:8080/healthz || exit 1

USER app:app

EXPOSE 8080

ENTRYPOINT [ "/app/cv-manager" ]