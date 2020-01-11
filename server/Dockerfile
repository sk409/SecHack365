FROM golang:1.13

RUN go get -u github.com/sk409/goconst
RUN go get -u github.com/sk409/gosession
RUN go get -u github.com/jinzhu/gorm
RUN go get -u github.com/jinzhu/gorm/dialects/mysql
RUN go get -u golang.org/x/crypto/bcrypt

WORKDIR /go/src/github.com/sk409/sechack365