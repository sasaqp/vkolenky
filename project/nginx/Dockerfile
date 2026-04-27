FROM nginx:alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY frontend-telegram /usr/share/nginx/html/telegram
COPY frontend-mobile /usr/share/nginx/html/mobile

COPY ssl/server.crt /etc/nginx/ssl/server.crt
COPY ssl/server.key /etc/nginx/ssl/server.key

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]