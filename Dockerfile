FROM scratch

COPY build/urlinfo /

VOLUME ["/data"]
EXPOSE 8080

CMD ["/urlinfo", "-datafile", "/data/dataset.txt"]