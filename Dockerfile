FROM scratch
COPY tmaze /usr/bin/tmaze
ENV HOME=/home/user
ENTRYPOINT ["/usr/bin/tmaze"]
