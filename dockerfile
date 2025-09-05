FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server cmd/server/main.go

FROM alpine:3.21.3

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/static ./static/
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./server"]


# FROM golang:1.23 AS builder

# WORKDIR /app

# COPY go.mod go.sum ./

# COPY . .

# RUN go mod download

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server cmd/server/main.go


# FROM alpine:3.21.3

# WORKDIR /app

# COPY --from=builder /app/server .
# COPY --from=builder /app/static ./static/

# COPY .env .env

# EXPOSE 8080

# CMD ["./server"]