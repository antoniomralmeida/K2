FROM mongo:5.0.5

WORKDIR /var/www/k2

COPY ./bin/*.bin ./bin/
COPY ./bin/*.env ./bin/
COPY ./log/*.md ./log/

COPY ./k2web/pub/ ./k2web/pub/
RUN ls -la ./k2web/pub/

EXPOSE 8080

CMD ["mongod" ]

