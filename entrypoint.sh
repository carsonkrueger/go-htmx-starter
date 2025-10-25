#!/bin/bash

cd /app

make migrate-internal

./bin/main -internal=true web
