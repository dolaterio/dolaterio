FROM scratch
ADD dolater /
ADD config_base.yml /config.yml
EXPOSE 7000
CMD ["/dolater"]
