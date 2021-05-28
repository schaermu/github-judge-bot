FROM alpine as builder
RUN apk add -U --no-cache ca-certificates

FROM scratch
ENTRYPOINT ["/github-judge"]
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY github-judge /