#!/bin/bash

cd /app

make migrate INTERNAL=true

./bin/main -internal=true web
