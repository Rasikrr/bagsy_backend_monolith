FROM golang:1.25.1-alpine AS builder

RUN apk add --no-cache make

WORKDIR /app

COPY . .

RUN go mod download

RUN make build


FROM golang:1.25.1-alpine AS runner

WORKDIR /app
COPY --from=builder /app/bin ./bin
COPY --from=builder /app/config ./config

CMD ["./bin/app"]
