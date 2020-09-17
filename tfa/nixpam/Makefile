MODULE := trasapam

debian:
	go build  -buildmode=c-shared -o ${MODULE}.so 
	sudo cp trasapam.so /lib/security/
	sudo systemctl restart sshd


