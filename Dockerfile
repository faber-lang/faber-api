FROM golang:1.12-alpine

ENV GOPATH=/go

WORKDIR /go/src/faber-api/
COPY . ./

RUN apk add -u curl ca-certificates
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN CGO_ENABLED=0 go build

FROM scratch

COPY --from=0 /go/src/faber-api/faber-api /usr/local/bin/faber-api

ENV GIN_MODE=release

EXPOSE 80
ENV PORT=80

CMD ["/usr/local/bin/faber-api"]