FROM local/test/sdcli
USER root
RUN apt-get update && apt-get -y install python3-testinfra
RUN mkdir /test
ADD *.py /test/
WORKDIR /test/
ENTRYPOINT ["py.test", "-v"]