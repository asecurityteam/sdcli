FROM local/test/sdcli
USER root
RUN pip3 install pytest-testinfra
RUN mkdir /test
ADD *.py /test/
WORKDIR /test/
ENTRYPOINT ["py.test", "-v"]