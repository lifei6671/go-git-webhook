FROM golang:1.8.1-alpine


ADD . /go/src/github.com/lifei6671/go-git-webhook


WORKDIR /go/src/github.com/lifei6671/go-git-webhook

RUN chmod +x start.sh

RUN  go build -ldflags "-w" && \
    rm -rf commands controllers models modules routers tasks vendor .gitignore .travis.yml Dockerfile gide.yaml LICENSE main.go README.md

CMD ["./start.sh"]