FROM golang:1.24.2-alpine AS buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/bohdan-vykhovanets/url-shortener-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/url-shortener-svc /go/src/github.com/bohdan-vykhovanets/url-shortener-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/url-shortener-svc /usr/local/bin/url-shortener-svc
RUN apk add --no-cache ca-certificates

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
