FROM ubuntu:18.04

LABEL base.image="biocontainers:latest"
LABEL version="1"
LABEL software="Philosopher"
LABEL description="A complete toolkit for shotgun proteomics data analysis"
LABEL website="https://philosopher.nesvilab.org"
LABEL documentation="https://philosopher/wiki"
LABEL license="GPL-3.0"
LABEL BIOTOOLS=""
LABEL tags="Proteomics"

# MAINTAINER Felipe da Veiga Leprevost <felipevl@umich.edu>

USER biodocker

ADD philosopher /home/biodocker/bin/

ENV PATH /home/biodocker/bin/:$PATH

WORKDIR /data/
