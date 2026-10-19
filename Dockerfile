# Etapa 1: compilación
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/api .

# Etapa 2: imagen final (solo el binario)
FROM alpine:3.20
WORKDIR /app
RUN adduser -D -u 1000 app
COPY --from=builder /app/api /usr/local/bin/api
COPY static ./static
USER app
EXPOSE 8080
CMD ["api"]
