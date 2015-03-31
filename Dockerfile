FROM google/golang:1.4

WORKDIR /gopath/src/github.com/artemave/conways-go

RUN echo "deb http://ftp.us.debian.org/debian wheezy-backports main" >> /etc/apt/sources.list
RUN apt-get update && apt-get install -y --no-install-recommends nodejs-legacy
RUN curl -L --insecure https://www.npmjs.org/install.sh | bash

ADD ./package.json /gopath/src/github.com/artemave/conways-go/
RUN npm install --production

ADD . /gopath/src/github.com/artemave/conways-go/

RUN go get github.com/tools/godep
RUN godep restore
RUN godep go build

ENV PORT 9999
CMD ["./conways-go"]
