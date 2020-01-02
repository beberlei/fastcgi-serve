FROM golang:1.12-alpine AS builder

RUN apk add --no-cache --no-progress git

WORKDIR /go/src/fastcgi-serve
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...


#############
# final image
#############
FROM alpine:3.10 AS final
EXPOSE 8080
ENTRYPOINT [ "fastcgi-serve" ]
COPY --from=builder /go/bin /usr/local/bin
