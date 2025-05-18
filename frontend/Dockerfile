FROM node:lts-alpine

RUN npm install -g http-server

WORKDIR /app/frontend

COPY package.json .

RUN npm install

COPY . .

RUN npm run build

CMD npx http-server --port ${PORT:-8081} dist
