#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    echo "Loading environment variables from .env file..."
    source .env
else
    echo "No .env file found."
fi

if [ "$ENV" == "development" ]; then
    echo "Environment is set to development. Running with air..."
    air -c .air.toml
elif [ "$ENV" == "production" ]; then
    echo "Environment is set to production. Building with go build..."
    ./asyncq
else
    echo "Environment variable ENV is not set to development or production. Defaulting to development"
    air -c .air.toml
fi
