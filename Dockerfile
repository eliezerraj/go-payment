#docker build -t go-payfee .
#docker run -dit --name go-payfee -p 5000:5000 go-payfee

FROM golang:1.21 As builder

WORKDIR /app
COPY . .

WORKDIR /app/cmd
RUN go build -o go-payment -ldflags '-linkmode external -w -extldflags "-static"'

FROM alpine

WORKDIR /app
COPY --from=builder /app/cmd/go-payment .

CMD ["/app/go-payment"]