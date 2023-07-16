# syntax=docker/dockerfile:1
FROM ubuntu:latest

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get -qq update -y && apt-get -qq upgrade -y

RUN mkdir -p /opt/dockerorc

COPY release/dockerorc-node /opt/dockerorc/dockerorc-node
RUN chmod -R +x /opt/dockerorc

CMD /opt/dockerorc/dockerorc-node