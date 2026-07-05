FROM alpine:3.21
RUN apk --no-cache add ca-certificates curl grpc-health-probe

COPY bin/authz /app

ENTRYPOINT ["/app"]
