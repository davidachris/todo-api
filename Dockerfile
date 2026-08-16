FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -ldflags '-s -w' -o server main.go

FROM scratch


COPY --from=builder /app/server /app/server

ENV AWS_DISABLED=true \
    GIN_MODE=release \
    PORT=8080

EXPOSE 8080
CMD ["/app/server"]
