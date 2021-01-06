package consts

type SSH_AUTH_TYPE string

const (
	SSH_AUTH_TYPE_PUB      = "PUB"      // trasa private key
	SSH_AUTH_TYPE_CERT     = "CERT"     // trasa cert
	SSH_AUTH_TYPE_PASSWORD = "PASSWORD" //using email and password

)
