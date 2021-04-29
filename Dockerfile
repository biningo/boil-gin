FROM golang:alpine as builder 
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN mkdir /build  
ADD . /build/
WORKDIR /build  
RUN go build -o boil-gin main.go
FROM alpine 
COPY --from=builder /build/boil-gin /app/
COPY --from=builder /build/config.yaml /app/
EXPOSE 9090
WORKDIR /app
CMD ["./boil-gin"]
