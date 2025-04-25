FROM golang:1.24.2-alpine


FROM caddy:2.10-alpine


ENV PORT=8080 \
    DATABASE_TYPE=sqlite \
    DATABASE_DSN=/app/data/production.db \
    GO_ENV=production \
    BASIC_AUTH_USER= \
    BASIC_AUTH_PASSWORD= \
    WEBUI_BASE_URL= \
    CADDY_BIN_PATH= \
    CADDY_DATA_PATH= \
    CADDY_RELOAD_CMD= \
    CADDY_TLS_EMAIL=
