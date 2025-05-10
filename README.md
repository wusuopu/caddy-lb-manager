# caddy-lb-manager
A simple webui to manage caddy lb rules.


## Install
init database:

```
docker run -it --rm -v data_volume:/data wusuopu/caddy-lb-manager init_db
```

start service:

```
docker run -it -d -v /opt/volumes/lb/caddy-lb-manager:/data wusuopu/caddy-lb-manager
```

open http://host:8080
