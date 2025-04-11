#!/bin/sh
set -e

echo "Running migrations..."
url-shortener-svc migrate up

echo "Starting service..."
url-shortener-svc run service
