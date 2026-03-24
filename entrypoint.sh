#!/bin/sh

if [ "$1" = "server" ]; then
    echo "Starting API server..."
    exec /app/server
elif [ "$1" = "cli" ]; then
    if [ -z "$2" ]; then
        echo "Usage: docker run --rm crypto-parser cli <coin>"
        echo "Example: docker run --rm crypto-parser cli btc"
        exit 1
    fi
    echo "🔍 Running CLI for coin: $2"
    exec /app/cli "$2"
else
    echo "Crypto Parser - Docker Edition"
    echo ""
    echo "Usage:"
    echo "  docker run --rm -p 8080:8080 crypto-parser server"
    echo "    Start the API server"
    echo ""
    echo "  docker run --rm crypto-parser cli <coin>"
    echo "    Run CLI client (needs server running separately)"
    echo ""
    echo "Examples:"
    echo "  docker run --rm -p 8080:8080 crypto-parser server"
    echo "  docker run --rm crypto-parser cli btc"
    exit 1
fi