FROM node:16 AS build

WORKDIR /app

COPY package*.json package-lock.json ./
RUN npm install
RUN npm audit fix --force
RUN npm install --save-dev webpack@latest
RUN npm install web-vitals
RUN npm install react@latest react-dom@latest
RUN npm install --save-dev @babel/plugin-proposal-private-property-in-object


COPY src ./src
COPY public ./public

RUN npm run build \
RUN ls -al /app/build

FROM nginx:alpine

COPY --from=build /app/build /usr/share/nginx/html
RUN chmod -R 755 /usr/share/nginx/html


COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
