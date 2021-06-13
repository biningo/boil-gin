FROM golang:1.16-alpine3.13 as builder
ENV GOPROXY=https://goproxy.io
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 go build -a -ldflags "-s -w" -o boil /build/

FROM scratch
COPY --from=builder /build/boil /
ENTRYPOINT ["/boil"]