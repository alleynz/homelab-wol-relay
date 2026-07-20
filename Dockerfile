FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/wol-relay ./cmd/wol-relay

FROM alpine:3.20
LABEL org.opencontainers.image.title="wol-relay" \
      org.opencontainers.image.description="Wake-on-LAN broadcast-to-unicast relay" \
      org.opencontainers.image.licenses="MIT"
RUN addgroup -S app && adduser -S -G app app
COPY --from=build /out/wol-relay /usr/local/bin/wol-relay
USER app
ENTRYPOINT ["/usr/local/bin/wol-relay"]
CMD ["-config", "/config/config.json"]
