FROM golang:1.9-alpine

WORKDIR /go/src/github.com/dannym87/go-notes-app

RUN apk --no-cache update \
    && apk --no-cache add git build-base

COPY . /go/src/github.com/dannym87/go-notes-app
COPY ./bin/run.sh /opt/bin/run.sh

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]