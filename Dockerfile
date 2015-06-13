FROM scratch
ADD api.bin /api
ADD migrate.bin /migrate
ADD worker.bin /worker
ADD config_base.yml /config.yml
EXPOSE 7000
CMD ["/worker"]
