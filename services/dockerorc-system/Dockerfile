# syntax=docker/dockerfile:1
FROM ubuntu:latest

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get -qq update -y && apt-get -qq upgrade -y

RUN mkdir -p /opt/dockerorc

COPY release/dockerorc-system /opt/dockerorc/dockerorc-system
RUN chmod -R +x /opt/dockerorc

EXPOSE 6443

CMD /opt/dockerorc/dockerorc-system