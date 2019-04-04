FROM ubuntu:18.04

# install golang
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install wget -y
RUN wget https://storage.googleapis.com/golang/go1.11.2.linux-amd64.tar.gz
RUN tar -xvf go1.11.2.linux-amd64.tar.gz
RUN mv go /usr/local

# set path for golang
ENV GOPATH=$HOME/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

RUN mkdir $HOME/go

RUN go get github.com/cvhariharan/SLB

WORKDIR $HOME/go/src/app
COPY . .

RUN go build -o main

EXPOSE 9000
RUN ls
CMD ["./main", "./conf.json"]