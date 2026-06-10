FROM golang:1.26-alpine AS base

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

FROM base AS development
RUN go install github.com/air-verse/air@latest
COPY . .
CMD [ "air","-c",".air.toml"]

FROM base AS builder
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd

FROM alpine:3.19 AS production

RUN apk --no-cache add ca-certificates wget \
    && addgroup -g 10001 -S appgroup \
    && adduser -u 10001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server .

USER appuser

ENV SERVER_PORT=8080

EXPOSE 8080

CMD ["./server"]