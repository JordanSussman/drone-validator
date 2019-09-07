FROM drone/ca-certs

ADD drone-validator /bin/

ENTRYPOINT ["/bin/drone-validator"]
