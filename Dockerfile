FROM ubuntu:jammy
ARG CHEEK_ARCH="linux/amd64"
EXPOSE 8081
WORKDIR /app
RUN apt update; apt install wget -y
RUN wget -P /usr/local/bin https://storage.googleapis.com/cheek-scheduler/${CHEEK_ARCH}/cheek
RUN chmod +x /usr/local/bin/cheek
ENTRYPOINT [ "cheek", "run" ]
CMD ["schedule_spec.yaml"]