FROM centos
EXPOSE 8888

RUN yum update -y

WORKDIR /app

COPY ./server ./server
COPY ./cli ./cli
CMD [ "./server" ]