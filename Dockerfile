FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN \
  BUILD_DT=`date +%FT%T%z` \
  && COMMIT=container \
  && VER=0.0.3 \
  && go build -ldflags \
    "-X main.build_dt=${BUILD_DT} -X main.commit=${COMMIT} -X main.version=${VERSION}" \
    -o main cmd/cli/cli.go
CMD ["/app/main"]
