FROM scratch
ADD dolater /
ADD config_base.yml /config.yml
CMD ["/dolater"]
