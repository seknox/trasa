# Trasapam

Pluggable Authentication Module for native two factor authentication agents for *nix platforms.

## Status
Currently only tested on Ubuntu 19.04, 20.04, Centos 6,7,8.

## Building
Use make file.  
In termminal (inside the project folder), `$ make` will (1) build `trasapam.so` file, (2) copy file in `/lib/security` and (3) restart `sshd` daemon

## Config
Config file trasapam.toml should be copied in `/etc/trasa/config/trasapam.toml`. Make sure to edit config file according to your service detail and trasa-server lcoation.

## Project structure

- CGO codes are included in `trasapam.go` and `trasapamUtils.go` file.
- `trasapamUtils.go` includes utility functions to initialize pam module and conversing with users.
- `trasapam.go` exports `pam_sm_authenticate` function which is go port of relevant C extern function where our logic for processing 2fa is performed.
- `utils.go` file has utility for reading config file, handleing 2fa requests and logging


## Debugging
Set `debug = true` in `trasapam.toml` config file
Create log file at `/var/log/trasapam.log`.

## Note:
Using go-curl package for http request because for some reason, go's built in httpclient is failing to make http request when invoked from inside pam module.