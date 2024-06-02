FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o pokedex .

FROM alpine:3
WORKDIR /root/
COPY --from=builder /app/pokedex .
ENV HTTP_PORT=8080
EXPOSE ${HTTP_PORT}
CMD ["./pokedex"]
