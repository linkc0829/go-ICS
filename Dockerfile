FROM golang:1.15-alpine
WORKDIR /go-ics
RUN apk update && apk upgrade && apk add --no-cache bash git openssh curl
COPY . /go-ics/
RUN go mod download

# workarround for SP-291. See https://github.com/oxequa/realize/issues/253
RUN go get github.com/urfave/cli/v2
RUN go get github.com/oxequa/realize

CMD ["./scripts/run-dev.sh"]