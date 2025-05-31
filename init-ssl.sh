#!/bin/bash

# Create required directories
mkdir -p data/certbot/conf
mkdir -p data/certbot/www

# Create temporary nginx config for certbot
cat > frontend/nginx.conf.tmp << 'EOL'
server {
    listen 80;
    listen [::]:80;
    server_name webman.stettzy.com;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 404;
    }
}
EOL

# Use temporary config
mv frontend/nginx.conf frontend/nginx.conf.bak
mv frontend/nginx.conf.tmp frontend/nginx.conf

# Start nginx with temporary config
docker compose -f docker-compose.prod.yml up -d frontend

# Wait for nginx to start
sleep 5

# Get the initial certificate
docker compose -f docker-compose.prod.yml run --rm certbot certonly --webroot \
  --webroot-path=/var/www/certbot \
  --email stettzy@gmail.com \
  --agree-tos \
  --no-eff-email \
  --staging \
  -d webman.stettzy.com

# If staging was successful, get the real certificate
docker compose -f docker-compose.prod.yml run --rm certbot certonly --webroot \
  --webroot-path=/var/www/certbot \
  --email stettzy@gmail.com \
  --agree-tos \
  --no-eff-email \
  -d webman.stettzy.com

# Restore original nginx config
mv frontend/nginx.conf.bak frontend/nginx.conf

# Start all services
docker compose -f docker-compose.prod.yml up -d

# Reload nginx to apply new configuration
echo "Reloading nginx configuration..."
docker compose -f docker-compose.prod.yml exec frontend nginx -s reload

echo "SSL setup complete!" 