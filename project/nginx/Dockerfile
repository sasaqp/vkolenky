FROM nginx:alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY frontend-telegram /usr/share/nginx/html/telegram
COPY frontend-mobile /usr/share/nginx/html/mobile

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]