FROM golang:1.15-alpine
WORKDIR /go-ics
RUN apk update && apk upgrade && apk add --no-cache bash git openssh curl
COPY . /go-ics
RUN go list -e $(go list -m all)

# workarround for SP-291. See https://github.com/oxequa/realize/issues/253
RUN go get gopkg.in/urfave/cli.v2@master
RUN go get github.com/tockins/realize

CMD ["./scripts/run-dev.sh"]