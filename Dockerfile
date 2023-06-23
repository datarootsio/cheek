FROM ubuntu
ARG CHEEK_ARCH="linux/amd64"
EXPOSE 8081
WORKDIR /app
RUN apt update; apt install wget -y
RUN wget https://storage.googleapis.com/better-unified/${CHEEK_ARCH}/cheek
RUN chmod +x cheek
ENTRYPOINT [ "./cheek" ]
CMD ["run", "job_spec.yaml"]