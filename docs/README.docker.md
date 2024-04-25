FLibGoLite
===

## Build docker container

```
docker build -t flibgolite .
```

## Download ready-to-use docker container from dockerhub

TODO

## Run docker container

```
docker run -d --name=flibgolite -v /srv/flibgolite:/var/flibgolite -p 8085:8085 flibgolite:latest
```

# Use docker container

Put your books into `/srv/flibgolite/books/new` to import new books into your collection.

Run `docker logs flibgolite` to read the logs.

The SQLite database is on the path `/srv/flibgolite/dbdata/flibgolite.db`. Backup it on a ragular basis.
