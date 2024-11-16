FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o bin/server ./cmd/server

FROM alpine

COPY --from=builder /app/bin/server /server

CMD [ "/server" ]
