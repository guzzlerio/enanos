FROM golang:latest


ENV ENANOS_VERBOSE false
ENV ENANOS_HOST 0.0.0.0
ENV ENANOS_MIN_SLEEP 1s
ENV ENANOS_MAX_SLEEP 10s
ENV ENANOS_RANDOM_SLEEP true
ENV ENANOS_MIN_SIZE 1KB
ENV ENANOS_MAX_SIZE 10KB
ENV ENANOS_RANDOM_SIZE true
ENV ENANOS_DEAD_TIME 5s

ADD . /go/src/github.com/reaandrew/enanos
RUN go get github.com/reaandrew/enanos
RUN go install github.com/reaandrew/enanos
ENTRYPOINT /go/bin/enanos
EXPOSE 8000
