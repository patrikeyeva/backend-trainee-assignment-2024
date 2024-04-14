FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build -o bin/avito-banner ./cmd/server

FROM alpine

COPY --from=builder /app/bin/avito-banner /avito-banner

CMD [ "/avito-banner" ]
