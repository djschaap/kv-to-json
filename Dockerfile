FROM golang:latest as builder
ARG COMMIT_HASH=container
ARG VER=0.0.4
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN \
  BUILD_DT=`date +%FT%T%z` \
  && echo "BUILD build_dt=${BUILD_DT}" \
  && echo "BUILD commit_hash=${COMMIT_HASH}" \
  && echo "BUILD version=${VER}" \
  && CGO_ENABLED=0 GOOS=linux go build -ldflags \
    "-X main.buildDt=${BUILD_DT} -X main.commit=${COMMIT_HASH} -X main.version=${VER}" \
    -o cli cmd/cli/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/.env /
COPY --from=builder /app/cli /
CMD ["/cli","8080"]
