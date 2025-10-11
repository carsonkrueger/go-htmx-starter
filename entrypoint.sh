#!/bin/bash

cd /app

make migrate-internal
make jet-all-internal
make build

./bin/main -internal=true web
