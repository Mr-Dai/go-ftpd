FROM golang:1.10.1-alpine

# Install Git
RUN apk add --update git
RUN rm -rf /var/cache/apk/*

# Install go-ftpd
WORKDIR /go/src/github.com/Mr-Dai/go-ftpd
ADD . .
RUN go get -v ./...
RUN go install -v github.com/Mr-Dai/go-ftpd

VOLUME /ftpd/auth.db
VOLUME /ftpd/data

EXPOSE 21

CMD ["go-ftpd", "-a", "/ftpd/auth.db", "run", "-d", "/ftpd/data"]
