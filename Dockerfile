FROM golang:1.12 AS build
ARG BUILD

WORKDIR /go/src/github.com/ehazlett/atlas
COPY . /go/src/github.com/ehazlett/atlas
RUN make

FROM alpine:latest
COPY --from=build /go/src/github.com/ehazlett/atlas/bin/* /bin/
ENTRYPOINT ["/bin/atlas"]
