#!/bin/sh
sudo yum localinstall https://dev.mysql.com/get/mysql80-community-release-el7-3.noarch.rpm -y
sudo yum-config-manager --disable mysql80-community
sudo yum-config-manager --enable mysql57-community
sudo yum install -y mysql-community-client