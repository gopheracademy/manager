#!/bin/bash

sudo apt update
sudo apt upgrade -y

sudo apt install libpng-dev build-essential -y
curl -sL https://deb.nodesource.com/setup_14.x | sudo -E bash -
sudo apt install nodejs -y

curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list

sudo apt update
sudo apt install yarn -y

# verify setup
node -v && npm -v && yarn -v

# add strapi user 
sudo adduser --shell /bin/bash --disabled-login --gecos "" --quiet strapi

sudo mkdir /srv/content

sudo chown strapi:strapi /srv/content 

# database
# https://www.digitalocean.com/community/tutorials/how-to-install-and-use-postgresql-on-ubuntu-18-04#step-4-%E2%80%94-creating-a-new-database

# app
cd /srv/content 
sudo su strapi && git clone https://github.com/gopheracademy/manager
cd manager
# temporary
sudo su strapi && git checkout cms


