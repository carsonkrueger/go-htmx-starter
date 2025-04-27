#!/bin/bash

cd /app

# DB-START
make migrate-internal
make jet-all-internal
# DB-END
make build

./bin/main -internal=true web
