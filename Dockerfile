FROM golang:1.19.5-alpine AS builder

ADD . /app

WORKDIR /app

RUN go build -o dist/server ./cmd/server/...

FROM alpine

COPY --from=builder /app/dist/server /bin/server

CMD [ "server" ]
