# ----------- Build Stage -----------
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o omnidata main.go

# ----------- Run Stage -----------
FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/omnidata ./omnidata
COPY README.md ./README.md
# Optionally copy config, scripts, etc.
# COPY . /app
ENTRYPOINT ["/app/omnidata"]
