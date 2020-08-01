package consts

type SSH_AUTH_TYPE string

const (
	SSH_AUTH_TYPE_PUB      = "PUB"      // trasa cert
	SSH_AUTH_TYPE_DACERT   = "DACERT"   //trasa  device agent certificate  (tfa already done, service already chosen)
	SSH_AUTH_TYPE_PASSWORD = "PASSWORD" //using email and password

)
