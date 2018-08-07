FROM golang AS builder
WORKDIR /go/src/github.com/jamesqin-cn/ldap_auth/
ADD . .
RUN go get -v && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine
MAINTAINER QinWuquan <jamesqin@vip.qq.com>
COPY --from=builder /go/src/github.com/jamesqin-cn/ldap_auth/app /bin/
ENTRYPOINT ["app"]
EXPOSE 9066
