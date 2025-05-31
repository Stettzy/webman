#!/bin/bash

# Build frontend
echo "Building frontend..."
cd frontend
npm install
npm run build

# Build and start containers
echo "Starting containers..."
cd ..
docker compose -f docker-compose.prod.yml up --build -d 