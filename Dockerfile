#docker build -t go-payment .
#docker run -dit --name go-payment -p 5000:5000 go-payment

FROM golang:1.21 As builder

RUN apt-get update && apt-get install bash && apt-get install -y --no-install-recommends ca-certificates

WORKDIR /app
COPY . .

WORKDIR /app/cmd
RUN go build -o go-payment -ldflags '-linkmode external -w -extldflags "-static"'

FROM alpine

WORKDIR /app
COPY --from=builder /app/cmd/go-payment .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app/go-payment"]