FROM ubuntu
# Note: change this url to fetch the
# binary for your arch (or specify your buid/run platformx)
WORKDIR /app
RUN apt update; apt install wget -y
RUN wget https://storage.googleapis.com/better-unified/linux/386/cheek
RUN chmod +x cheek
ENTRYPOINT [ "./cheek" ]