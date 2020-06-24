FROM golang
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN \
  BUILD_DT=`date +%FT%T%z` \
  && COMMIT=container \
  && VER=0.0.4 \
  && go build -o main -ldflags \
    "-X main.buildDt=${BUILD_DT} -X main.commit=${COMMIT} -X main.version=${VERSION}" \
    cmd/cli/main.go
CMD ["/app/main"]
