#!/usr/bin/env bash

sudo apt-get -y install libtool pkg-config build-essential autoconf automake
sudo apt-get -y install libzmq-dev

git clone git://github.com/jedisct1/libsodium.git
cd libsodium
./autogen.sh
./configure && make check
sudo make install
sudo ldconfig

cd ..

wget https://github.com/zeromq/libzmq/releases/download/v4.2.1/zeromq-4.2.1.tar.gz
tar -xvf zeromq-4.2.1.tar.gz
cd zeromq-4.2.1
./autogen.sh
./configure && make check
sudo make install
sudo ldconfig
