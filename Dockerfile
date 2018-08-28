FROM alpine:3.8
RUN apk update && apk add ca-certificates
EXPOSE 8080
ENTRYPOINT ["/golang-http"]
COPY ./bin/ /
