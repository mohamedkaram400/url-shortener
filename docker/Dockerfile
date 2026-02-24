# ---------- Build stage ----------
FROM golang:1.25.6-alpine AS builder

ENV CGO_ENABLED=0

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

# ---------- Runtime stage ----------
FROM alpine:3.21

RUN apk add --no-cache curl

WORKDIR /app

COPY --from=builder /build/main .
COPY --from=builder /build/db ./db 

EXPOSE 9000

CMD [ "./main" ]