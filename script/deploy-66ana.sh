#!/usr/bin/env bash

apt-get update

apt-get install -y tzdata
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
echo "Asia/Shanghai" > /etc/timezone
dpkg-reconfigure -f noninteractive tzdata

export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:${GOPATH//://bin:}/bin

shopt -s expand_aliases
alias proxy='env all_proxy=socks5://10.0.0.220:1080 http_proxy=http://10.0.0.220:8123 https_proxy=http://10.0.0.220:8123'

echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.bashrc
echo 'export GOPATH=$HOME/go' >> $HOME/.bashrc
echo 'export PATH=$PATH:${GOPATH//://bin:}/bin' >> $HOME/.bashrc

echo "alias proxy='env all_proxy=socks5://10.0.0.220:1080 http_proxy=http://10.0.0.220:8123 https_proxy=http://10.0.0.220:8123'" >> $HOME/.bashrc

proxy wget -O go1.8.linux-amd64.tar.gz https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
tar xvf go1.8.linux-amd64.tar.gz -C /usr/local/

go version

apt-get install git -y

retry=5
while [ $retry -gt 0 ]; do
    echo go get glide retry: $retry
    proxy go get -u -v github.com/Masterminds/glide
    if [ $? -eq 0 ]; then
        break;
    fi
    retry=$[$retry-1]
done

if [ $retry -eq 0 ]; then
    echo glide install failed
    exit 1
fi

glide help

mkdir -p ~/go/src/github.com/zwh8800
cd ~/go/src/github.com/zwh8800
proxy git clone https://github.com/zwh8800/66ana
cd 66ana

retry=5
while [ $retry -gt 0 ]; do
    echo glide install retry: $retry
    proxy glide install
    if [ $? -eq 0 ]; then
        break;
    fi
    retry=$[$retry-1]
done

if [ $retry -eq 0 ]; then
    echo glide install failed
    exit 1
fi

cd ~

apt-get -y install libtool pkg-config build-essential autoconf automake
apt-get -y install libzmq-dev

git clone git://github.com/jedisct1/libsodium.git
cd libsodium
./autogen.sh
./configure && make check
make install
ldconfig

cd ..

retry=5
while [ $retry -gt 0 ]; do
    echo go get zeromq retry: $retry
    proxy wget -O zeromq-4.2.1.tar.gz https://github.com/zeromq/libzmq/releases/download/v4.2.1/zeromq-4.2.1.tar.gz

    if [ $? -eq 0 ]; then
        break;
    fi
    retry=$[$retry-1]
done

if [ $retry -eq 0 ]; then
    echo glide install failed
    exit 1
fi

tar -xvf zeromq-4.2.1.tar.gz
cd zeromq-4.2.1
./autogen.sh
./configure && make check
make install
ldconfig

cd ~/go/src/github.com/zwh8800/66ana
go build

if [ $? -eq 0 ]; then
    echo build ok
else
    echo build failed
fi

cd ~
ln -s ~/go/src/github.com/zwh8800/66ana 66ana
