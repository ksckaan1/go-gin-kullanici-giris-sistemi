FROM golang:1.15.0-alpine 

WORKDIR /go/src/app
COPY . .

RUN go mod download 
RUN go build .

CMD ["./go-gin-kullanici-giris-sistemi"]