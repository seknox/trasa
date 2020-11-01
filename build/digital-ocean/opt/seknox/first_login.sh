#!/usr/bin/env sh

sudo systemctl stop trasa

echo "Enter TRASA server domain"
read trasadomain

sudo sed -i -e 's|<TRASA_LISTEN_ADDR>|'$trasadomain'|g' /etc/trasa/config/config.toml



echo "Do you want to retrieve certificate from Let's Encrypt? (Y/N) "
read ans



if [ $ans = 'Y' ] || [ $ans = 'y' ]; then
  sudo sed -i -e 's|<AUTO_CERT>|true|g' /etc/trasa/config/config.toml
fi

if [ $ans = 'N' ] || [ $ans = 'n' ]; then
  sudo sed -i -e 's|<AUTO_CERT>|false|g' /etc/trasa/config/config.toml
fi





cp -f /etc/skel/.bashrc /root/.bashrc
sudo systemctl restart trasa
