FROM alpine:3.7
LABEL author="gang.tao@outlook.com"

RUN mkdir /home/candy
WORKDIR /home/candy

COPY ./server/server /home/candy/
COPY ./client/client /home/candy/

EXPOSE 50051

ENTRYPOINT ["/home/candy/server"]