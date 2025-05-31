#!/bin/bash

# Create required directories
mkdir -p data/certbot/conf
mkdir -p data/certbot/www

# Start nginx
docker compose -f docker-compose.prod.yml up -d frontend

# Wait for nginx to start
sleep 5

# Stop nginx
docker compose -f docker-compose.prod.yml stop frontend

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

# Start all services
docker compose -f docker-compose.prod.yml up -d

# Update nginx configuration to include SSL
echo "Updating nginx configuration..."
cat > frontend/nginx.conf << 'EOL'
server {
    listen 80;
    listen [::]:80;
    server_name webman.stettzy.com;

    # Redirect all HTTP traffic to HTTPS
    location / {
        return 301 https://$host$request_uri;
    }

    # For certbot validation
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }
}

server {
    listen 443 ssl;
    listen [::]:443 ssl;
    server_name webman.stettzy.com;

    # SSL configuration
    ssl_certificate /etc/letsencrypt/live/webman.stettzy.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/webman.stettzy.com/privkey.pem;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:50m;
    ssl_session_tickets off;

    # Modern configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    # HSTS (uncomment if you are sure)
    # add_header Strict-Transport-Security "max-age=63072000" always;

    root /usr/share/nginx/html;
    index index.html;

    # Handle Single Page Application routing
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Proxy API requests to backend
    location /api/ {
        proxy_pass http://backend:9090/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
EOL

# Reload nginx to apply new configuration
echo "Reloading nginx configuration..."
docker compose -f docker-compose.prod.yml exec frontend nginx -s reload

echo "SSL setup complete!" 