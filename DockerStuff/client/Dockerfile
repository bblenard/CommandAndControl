FROM centos
EXPOSE 8888


RUN yum update -y

WORKDIR /app

COPY ./ClientID ./...
COPY ./client ./client
CMD [ "./client" ]