FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY Backend/go.mod Backend/go.sum ./
RUN go mod download

COPY Backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/server .

ENV SERVER_PORT=8080

EXPOSE 8080

CMD ["./server"]