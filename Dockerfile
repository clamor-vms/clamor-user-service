FROM drone/ca-certs

ADD src/user /

CMD ["/user"]
