FROM golang:1.13

RUN go get -u github.com/google/uuid
RUN go get -u github.com/jinzhu/gorm
RUN go get -u github.com/jinzhu/gorm/dialects/mysql
RUN go get -u github.com/sk409/goconst
RUN go get -u github.com/sk409/gofile
RUN go get -u github.com/sk409/gosession
RUN go get -u golang.org/x/crypto/bcrypt

RUN apt-get update && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

ENV DOCKER_VERSION 18.09.1
RUN curl -fL -o docker.tgz "https://download.docker.com/linux/static/test/x86_64/docker-$DOCKER_VERSION.tgz" && \
    tar --strip-component=1 -xvaf docker.tgz -C /usr/bin
RUN curl -L https://github.com/docker/compose/releases/download/1.24.1/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose && chmod +x /usr/local/bin/docker-compose


WORKDIR /go/src/github.com/sk409/sechack365