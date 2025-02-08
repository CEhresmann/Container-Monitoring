FROM node:latest AS build

WORKDIR /app

COPY package*.json package-lock.json ./
RUN npm install
RUN npm install web-vitals
RUN npm install react@latest react-dom@latest


COPY src ./
COPY . .

ENV NODE_OPTIONS=--openssl-legacy-provider
RUN npm run build

FROM nginx:alpine

COPY --from=build /app/build /usr/share/nginx/html

COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
