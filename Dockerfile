FROM golang:1.12 AS build
ARG BUILD

WORKDIR /go/src/github.com/stellarproject/atlas
COPY . /go/src/github.com/stellarproject/atlas
RUN make

FROM alpine:latest
COPY --from=build /go/src/github.com/stellarproject/atlas/bin/* /bin/
ENTRYPOINT ["/bin/atlas"]
