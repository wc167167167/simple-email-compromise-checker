FROM nginx:stable-alpine-perl
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/
CMD ["nginx", "-g", "daemon off;"]