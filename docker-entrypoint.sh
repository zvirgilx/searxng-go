#!/bin/sh
set -e

if [ "$1" = "api" ]; then
    exec /app/searxng-go "api" "$@"
fi

if [ "$1" = 'web' ]; then
    exec node "/app/.output/server/index.mjs"
fi
