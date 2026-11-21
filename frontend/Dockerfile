FROM oven/bun:1.4-alpine AS build

WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .
RUN bun run build

FROM nginx:1.30.7-alpine@sha256:caf4e1e032a20f9c3c5e28fc7f25ffd768ad8677b4975c54193075997b3a97e3

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
