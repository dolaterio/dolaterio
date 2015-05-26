FROM scratch
ADD dolater.bin /dolater
ADD migrate.bin /migrate
ADD config_base.yml /config.yml
EXPOSE 7000
CMD ["/dolater"]
