FROM alpine:latest as builder
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add ca-certificates && update-ca-certificates
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "10001" \
    "appuser"

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY external-dns-webhook-libdns /external-dns-webhook-libdns
USER appuser:appuser
ENTRYPOINT ["/external-dns-webhook-libdns"]