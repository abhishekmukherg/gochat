FROM golang:1.4

COPY main.go /go/src/app/main.go
COPY src /go/src/github.com/linkinpark342/gchat
WORKDIR /go/src/app
RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]
