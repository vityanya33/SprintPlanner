# === STAGE 1: Build the app ===
FROM node:20-alpine AS build-stage

WORKDIR /app

# Устанавливаем зависимости
COPY package*.json ./
RUN npm install

# Копируем остальные исходники и собираем приложение
COPY . .
RUN npm run build

# === STAGE 2: Serve the app with nginx ===
FROM nginx:alpine AS production-stage

# Копируем собранные файлы из предыдущего этапа
COPY --from=build-stage /app/dist /usr/share/nginx/html

# Копируем кастомный конфиг nginx (если нужно)
# COPY nginx.conf /etc/nginx/nginx.conf

# Порт (по умолчанию nginx использует 80)
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]