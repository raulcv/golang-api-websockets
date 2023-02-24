from postgres:10.3

# Copy the up.sql content file into 1.sql, this file will be into docker container 
COPY up.sql /docker-entrypoint-initdb.d/1.sql

# this commads run the postgres server
CMD ["postgres"]