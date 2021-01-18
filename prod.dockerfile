# Multistaged build production golang service
FROM golang:1.15-alpine as base

FROM base AS ci

RUN apk update && apk upgrade && apk add --no-cache git
RUN mkdir /build
ADD . /build/
WORKDIR /build

# Build prod
FROM ci AS build-env

RUN go mod download

RUN sudo apt-get install g++-arm-linux-gnueabi
RUN sudo apt-get install gcc-arm-linux-gnueabi

RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabi-gcc \
    go build -a -installsuffix \
    cgo -ldflags '-extldflags "-static"' -o server ./cmd/ics/

FROM alpine AS prod
RUN apk --no-cache add ca-certificates

COPY --from=build-env build/server ./ics/

CMD ["./ics/server"]