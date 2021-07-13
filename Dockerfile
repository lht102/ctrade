FROM golang:1.16-buster AS test-env
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.41.1
WORKDIR /go/src/ctrade
ADD . /go/src/ctrade
RUN make test-all

FROM golang:1.16-buster AS build-env
WORKDIR /go/src/ctrade
ADD . /go/src/ctrade
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s' \
    -o /go/bin/ctraded ./cmd/ctraded

FROM gcr.io/distroless/static AS runtime
COPY --from=build-env /go/bin/ctraded /go/bin/ctraded
USER nonroot
ENTRYPOINT ["/go/bin/ctraded"]
