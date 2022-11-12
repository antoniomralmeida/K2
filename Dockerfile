FROM mongo:5.0.5

WORKDIR /var/www/k2

COPY ./bin/*.bin ./bin/
COPY ./bin/*.env ./bin/
COPY ./log/*.json ./log/

COPY ./k2web/web/pub/css/* ./k2web/web/pub/css/
COPY ./k2web/web/pub/ebnf/* ./k2web/web/pub/ebnf/
COPY ./k2web/web/pub/img/* ./k2web/web/pub/img/
COPY ./k2web/web/pub/js/* ./k2web/web/pub/js/
COPY ./k2web/web/pub/scss/* ./k2web/web/pub/scss/
COPY ./k2web/web/pub/vendor/* ./k2web/web/pub/vendor/
COPY ./k2web/web/pub/view/* ./k2web/web/pub/view/


EXPOSE 8080

CMD ["mongod" ]



