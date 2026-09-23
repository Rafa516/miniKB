# --- Build stage ---
FROM node:22-alpine AS build
WORKDIR /src

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# VITE_API_URL é opcional: se não for passado, o frontend descobre o
# endereço do backend sozinho a partir do hostname usado para abrir a
# página (ver frontend/src/services/api.js) — o que funciona bem quando
# frontend e backend estão no mesmo host (Docker Compose local). Em
# serviços como o Railway, onde cada um fica num domínio diferente, é
# preciso passar essa variável explicitamente no build:
#   docker build --build-arg VITE_API_URL=https://seu-backend.up.railway.app frontend/
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

RUN npm run build

# --- Run stage ---
FROM nginx:1.27-alpine

# vite.config.js sets base: '/miniKB/', so the built assets expect to be
# served from that sub-path. Serving them there (instead of editing the
# source config) keeps the app buildable unchanged for both GitHub Pages
# and Docker.
COPY --from=build /src/dist /usr/share/nginx/html/miniKB

# Template (não um .conf fixo): a imagem oficial do Nginx substitui
# ${PORT} sozinha ao iniciar o container (envsubst), a partir da
# variável de ambiente PORT — necessário para rodar no Railway.
COPY nginx.conf.template /etc/nginx/templates/default.conf.template
ENV PORT=80

EXPOSE 80
