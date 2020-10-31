FROM alpine:3.7
LABEL author="gang.tao@outlook.com"

RUN mkdir /home/candy
WORKDIR /home/candy

COPY ./candy /home/candy/

ENTRYPOINT ["/home/candy/candy"]