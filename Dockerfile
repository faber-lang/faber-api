FROM golang:1.12-alpine

ENV GOPATH=/go

RUN apk add -u curl ca-certificates git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/faber-api/
COPY . ./

RUN dep ensure -v
RUN CGO_ENABLED=0 go build

FROM scratch

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/src/faber-api/faber-api /usr/local/bin/faber-api

ENV GIN_MODE=release

EXPOSE 8080
ENV PORT=8080

WORKDIR /tmp

ENTRYPOINT ["/usr/local/bin/faber-api"]