FROM golang:1.24.2-alpine as builder

COPY ./src /app
WORKDIR /app/

ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=auto

RUN go build -ldflags "-s -w" -o dist/goose cmd/goose.go && \
    go build -tags=jsoniter -ldflags "-s -w" -o dist/app . && \
    mv -v run.sh dist/


FROM caddy:2.10-alpine

RUN apk add --no-cache tini

COPY --from=builder /app/dist/ /app
WORKDIR /app/

VOLUME ["/data"]

EXPOSE 80 443 8080

ENV PORT=8080 \
    DATABASE_TYPE=sqlite \
    DATABASE_DSN=/data/webui/production.db \
    GO_ENV=production \
    WEBUI_BASIC_AUTH_USER= \
    WEBUI_BASIC_AUTH_PASSWORD= \
    WEBUI_BASE_URL= \
    WEBUI_ENABLE_CLUSTER_MODE= \
    CADDY_BIN_PATH=caddy \
    CADDY_DATA_PATH=/data/caddy \
    CADDY_CONFIG_PATH=/data/Caddyfile \
    CADDY_RELOAD_CMD= \
    CADDY_TLS_EMAIL=admin@example.com

ENTRYPOINT ["/sbin/tini", "--", "sh", "run.sh"]

CMD ["start_all"]
