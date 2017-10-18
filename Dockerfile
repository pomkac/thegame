FROM ubuntu:14.04

WORKDIR /root

RUN apt-get update 
RUN apt-get install -y wget git

RUN wget https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz
	
RUN sudo tar -C /usr/local -xzf go1.9.linux-amd64.tar.gz && \
    mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg

ENV PATH=${PATH}:/usr/local/go/bin GOROOT=/usr/local/go GOPATH=/root/go

RUN go get github.com/pomkac/mnemonic

RUN go get github.com/pomkac/thegame 

RUN go build github.com/pomkac/thegame
 
EXPOSE 80

CMD /root/thegame