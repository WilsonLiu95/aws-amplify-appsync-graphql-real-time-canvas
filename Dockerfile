FROM node:10-slim

RUN userdel -r node && useradd -m -u 1000 -s /bin/bash faas

COPY index.js /opt/application/index.js
COPY run.sh /opt/application/run.sh

WORKDIR /opt/application
USER root
CMD /opt/application/run.sh
