#!/bin/sh

start_lb () {
  config=""
  if [[ -e "/data/Caddyfile" ]]; then
    config="--adapter caddyfile --config /data/Caddyfile"
  fi

  if [[ "$1" == "true" ]]; then
    # run in background
    echo caddy start $config
    caddy start $config
  else
    # run in foreground
    echo caddy run $config
    caddy run $config
  fi
}

start_server () {
  ./app
}

start_all () {
  start_lb true
  start_server
}

init_db () {
  if [[ "$DATABASE_TYPE" == "sqlite" ]]; then
    ./goose db:create
  fi
  ./goose up
}

if [ "$1" = 'start_all' ]; then
  start_all
elif [ "$1" = 'start_lb' ]; then
  start_lb
elif [ "$1" = 'start_server' ]; then
  start_server
elif [ "$1" = 'init_db' ]; then
  init_db
else
  echo $@
  exec $@
fi
