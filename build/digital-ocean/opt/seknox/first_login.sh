#!/usr/bin/env sh

sudo systemctl stop trasa

echo "Enter TRASA server domain"
read trasadomain

sudo sed -i -e 's|<TRASA_LISTEN_ADDR>|'$trasadomain'|g' /Users/bhrg3se/seknox/code/trasa/trasa-oss/build/etc/trasa/config/config.toml



echo "Do you want to retrieve certificate from Let's Encrypt? (Y/N) "
read ans



if [ $ans = 'Y' ] || [ $ans = 'y' ]; then
  sudo sed -i -e 's|<AUTO_CERT>|true|g' /Users/bhrg3se/seknox/code/trasa/trasa-oss/build/etc/trasa/config/config.toml
fi



cp -f /etc/skel/.bashrc /root/.bashrc
sudo systemctl start trasa
