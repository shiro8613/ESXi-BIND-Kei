FROM debian

RUN mkdir /data

COPY ./build/dns /data/dns

RUN chmod +x /data/dns

COPY ./entrypoint.sh /entrypoint.sh

CMD ["/bin/sh", "/entrypoint.sh"]