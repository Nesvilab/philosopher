FROM biocontainers/biocontainers:latest

LABEL base.image="biocontainers:latest"
LABEL version="1"
LABEL software="Philosopher"
LABEL software.version="20190313"
LABEL description="A complete toolkit for shotgun proteomics data analysis"
LABEL website="https://philosopher.nesvilab.org"
LABEL documentation="https://github.com/Nesvilab/philosopher/wiki"
LABEL license="GPL-3.0"
LABEL BIOTOOLS=""
LABEL tags="Proteomics"

MAINTAINER Felipe da Veiga Leprevost <felipevl@umich.edu>

USER biodocker

RUN wget https://github.com/prvst/philosopher/releases/download/20190313/philosopher_linux_amd64 -P /home/biodocker/bin/Philosopher/ && \
  chmod -R 755 /home/biodocker/bin/Philosopher/*

ENV PATH /home/biodocker/bin/Philosopher/:$PATH

WORKDIR /data/

CMD ["philosopher_linux_amd64"]
