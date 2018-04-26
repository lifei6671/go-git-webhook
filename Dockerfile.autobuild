FROM golang:1.9 AS build

ENV PROJECT /go/src/github.com/lifei6671/go-git-webhook

RUN mkdir -p $PROJECT

WORKDIR ${PROJECT}

# ADD https://api.github.com/repos/lifei6671/go-git-webhook/compare/master...HEAD /dev/null
RUN git clone https://github.com/lifei6671/go-git-webhook.git .

RUN  curl https://glide.sh/get | sh \
  && glide install \
  && go get github.com/beego/bee \
  && CGO_ENABLED=0 CGO_ENABLED=0 bee pack -o ./bin

FROM alpine:3.6

ENV PROJECT /go/src/github.com/lifei6671/go-git-webhook

WORKDIR /go-git-webhook

RUN apk add --no-cache ca-certificates \
  && mkdir logs \
  && touch logs/log.log

COPY --from=build $PROJECT/bin/* ./

RUN tar -zxvf go-git-webhook.tar.gz \
    && rm -rf go-git-webhook.tar.gz \
    && ln -s /go-git-webhook/go-git-webhook /usr/bin/

