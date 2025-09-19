FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN make build


FROM golang:1.25.1-alpine AS runner

WORKDIR /app
COPY --from=builder /app/bin ./bin
COPY --from=builder /app/configs ./configs

CMD ["./bin/app"]
