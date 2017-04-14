FROM golang:1.8.1-alpine

ADD . /usr/src/github.com/lifei6671/go-git-webhook


WORKDIR /usr/src/github.com/lifei6671/go-git-webhook

RUN chmod +x start.sh

RUN go get -v -u ./... && \
    go build -ldflags "-w"

CMD ["./start.sh"]