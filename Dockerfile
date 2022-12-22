FROM golang:1.18

# Copy source code
RUN mkdir -p /usr/src/julo-test
COPY . /usr/src/julo-test
WORKDIR /usr/src/julo-test

RUN go install github.com/githubnemo/CompileDaemon@latest
CMD CompileDaemon -log-prefix=false -build="go build" -command="./julo-test"

EXPOSE 5000