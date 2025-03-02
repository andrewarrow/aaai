#!/bin/bash

# Create the db directory structure if it doesn't exist
mkdir -p db

# Touch the database file if it doesn't exist (Flyway needs the file to exist)
if [ ! -f "./db/vibecoders.db" ]; then
  touch ./db/vibecoders.db
  echo "Created empty database file at ./db/vibecoders.db"
fi

# Run Flyway migrations
echo "Running database migrations..."
npx flyway -configFiles=db/flyway.conf migrate

echo "Database initialization complete"