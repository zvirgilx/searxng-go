#!/bin/sh
set -e

if [ "$1" = "api" ]; then
    exec /app/searxng-go "api" "$@"
fi

if [ "$1" = 'web' ]; then
    if [ ! -f /app/.output/server/index.mjs ]; then
        cd /app/web
        npm run dev
    else
        exec node "/app/.output/server/index.mjs"
    fi
fi
