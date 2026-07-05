FROM alpine:3.21
RUN apk --no-cache add ca-certificates curl

COPY bin/account /app

ENTRYPOINT ["/app"]
