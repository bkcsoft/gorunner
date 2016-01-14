FROM golang

EXPOSE 8090
RUN mkdir -p /app
WORKDIR $GOPATH/src/github.com/jakecoffman/gorunner

ENV SHELL bash

RUN apt-get update
RUN apt-get install jq
RUN go get github.com/tools/godep

COPY Godeps/ $GOPATH/src/github.com/jakecoffman/gorunner/Godeps
RUN godep restore

ADD . $GOPATH/src/github.com/jakecoffman/gorunner/

RUN go build -v && \
		cp gorunner /app/gorunner
COPY web /app/web/

WORKDIR /app/
VOLUME ["/app/db", "/app/data"]

CMD ["./gorunner"]
