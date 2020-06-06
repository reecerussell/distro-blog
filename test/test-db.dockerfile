FROM mysql

ENV MYSQL_ROOT_PASSWORD=password

COPY test-db-setup.sql /docker-entrypoint-initdb.d/