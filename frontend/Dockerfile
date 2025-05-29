FROM node:20-alpine AS base

WORKDIR /app

# Development stage
FROM base AS development
EXPOSE 3000
CMD ["npm", "run", "dev"]

# Production stage
FROM base AS production

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Install serve to run the built app
RUN npm install -g serve

# Expose port
EXPOSE 3000

# Serve the built application
CMD ["serve", "-s", "dist", "-l", "3000"]
