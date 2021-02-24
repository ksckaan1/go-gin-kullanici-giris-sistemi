FROM golang:1.15.0 AS builder

WORKDIR /app
COPY . . 

RUN go mod download 
ENV CGO_ENABLED=0
RUN go build -o main . 

FROM alpine:3.13

COPY --from=builder /app . 

CMD ["./main"]