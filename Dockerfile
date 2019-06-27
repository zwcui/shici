FROM golang:latest

MAINTAINER zwcui<zwcui2017@163.com>

ENV kpdir /go/src/shici

RUN mkdir -p ${kpdir}

ADD . ${kpdir}/

WORKDIR ${kpdir}

RUN go build -v

EXPOSE 8083

ENTRYPOINT ["./shici"]