FROM google/golang:1.4

WORKDIR /gopath/src/github.com/artemave/conways-go

RUN echo "deb http://ftp.us.debian.org/debian wheezy-backports main" >> /etc/apt/sources.list
RUN apt-get update && apt-get install -y --no-install-recommends nodejs-legacy
RUN curl -L --insecure https://www.npmjs.org/install.sh | bash

ADD ./package.json /gopath/src/github.com/artemave/conways-go/
ADD ./gulpfile.js /gopath/src/github.com/artemave/conways-go/
ADD ./public /gopath/src/github.com/artemave/conways-go/
RUN npm install --production

ADD . /gopath/src/github.com/artemave/conways-go/
RUN go get github.com/tools/godep
RUN godep restore
RUN godep go install

CMD []
ENV PORT 9999
ENTRYPOINT ["/gopath/bin/conways-go"]
