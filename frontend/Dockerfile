# --- Build stage ---
FROM node:22-alpine AS build
WORKDIR /src

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# --- Run stage ---
FROM nginx:1.27-alpine

# vite.config.js sets base: '/miniKB/', so the built assets expect to be
# served from that sub-path. Serving them there (instead of editing the
# source config) keeps the app buildable unchanged for both GitHub Pages
# and Docker.
COPY --from=build /src/dist /usr/share/nginx/html/miniKB
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
