FROM biocontainers/biocontainers:latest

LABEL base.image="biocontainers:latest"
LABEL version="1"
LABEL software="Philosopher"
LABEL software.version="20180302"
LABEL description="A tool for Proteomics data analysis and post-processing filtering"
LABEL website="https://prvst.github.io/philosopher/"
LABEL documentation="https://prvst.github.io/philosopher/documentation.html"
LABEL license=""
LABEL BIOTOOLS=""
LABEL tags="Proteomics"

MAINTAINER Felipe da Veiga Leprevost <felipe@leprevost.com.br>

USER biodocker

RUN wget https://github.com/prvst/philosopher/releases/download/20180302/philosopher_linux_amd64 -P /home/biodocker/bin/Philosopher/ && \
  chmod -R 755 /home/biodocker/bin/Philosopher/*

ENV PATH /home/biodocker/bin/Philosopher/:$PATH

WORKDIR /data/

CMD ["philosopher_linux_amd64"]
