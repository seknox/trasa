package testutils

import (
	"github.com/gorilla/websocket"
	"net/http"
)

const (
	MockOrgID  = "153f7582-5ae2-46ba-8c1c-79ef73fe296e"
	MockUserID = "13c45cfb-72ca-4177-b968-03604cab6a27"

	MockTrasaID      = "root"
	MockTrasaID2     = "bhbha@gmail.com"
	MocktrasaPass    = "changeme"
	MocktrasaPass2   = "Changeme@123"
	MocktotpSEC      = "AV2COXZHVG4OAFSF"
	MockupstreamUser = "root"
	MockupstreamPass = "root"

	MockHostCert = `ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBGAawajB0kLn5QKDC4KBxS8SBxJMzw6Hu+54Yw6O7vylfLXat3hmg/xo5VzRk/zpVLZ5ZG9FzX2TemmdvQC7aIY=`

	MockPrivateKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABlwAAAAdzc2gtcn
NhAAAAAwEAAQAAAYEArLShu9zqfvahOONl9nziHCew23RF89SqjB/E1i6H4E/5aFsGfoGV
g5HEvTCdKNQTZdp3xW5sNDvsp9kEoykkiBlknbR+bNfkq8hTKG+CU8VK8jCmZLqjLa7qhK
H3l0ChKtbT/k9xSdgoVn7sqrhg/nZBXD0lnWniZ+Day/HfgPzsejW8DSj5YTKSYDubINMw
5QoFB+f6MLe7pT6lK6yD5O7sK26g3Y/QztAHB+vrmbD91KL82O+Rf2qVZp94GP3R9DcHAB
KnyoMpudAvJ/67jhTlXo8eXywdp0gajxVom1LHhdm9NP3MZJK6tPRe+epfDVzLv9SM8v77
fh41WYn73paaj7r7xBHF072N9VyWqu2dEq9I2qlGu1ZjVZ1JZD17vopdK35Aq4wzNTo1Tv
LNxh8JbLeKm+Oi7xYNDF1bYhljZ1dlPyPvzjh839vuqTMh+HzZZ9Gh/5KQowgSKOdNvlTw
AWHfviub3LJYsEBXtQ7fGjJn2QfH7G7GD0wIFWZTAAAFmJ12BgiddgYIAAAAB3NzaC1yc2
EAAAGBAKy0obvc6n72oTjjZfZ84hwnsNt0RfPUqowfxNYuh+BP+WhbBn6BlYORxL0wnSjU
E2Xad8VubDQ77KfZBKMpJIgZZJ20fmzX5KvIUyhvglPFSvIwpmS6oy2u6oSh95dAoSrW0/
5PcUnYKFZ+7Kq4YP52QVw9JZ1p4mfg2svx34D87Ho1vA0o+WEykmA7myDTMOUKBQfn+jC3
u6U+pSusg+Tu7CtuoN2P0M7QBwfr65mw/dSi/NjvkX9qlWafeBj90fQ3BwASp8qDKbnQLy
f+u44U5V6PHl8sHadIGo8VaJtSx4XZvTT9zGSSurT0XvnqXw1cy7/UjPL++34eNVmJ+96W
mo+6+8QRxdO9jfVclqrtnRKvSNqpRrtWY1WdSWQ9e76KXSt+QKuMMzU6NU7yzcYfCWy3ip
vjou8WDQxdW2IZY2dXZT8j7844fN/b7qkzIfh82WfRof+SkKMIEijnTb5U8AFh374rm9yy
WLBAV7UO3xoyZ9kHx+xuxg9MCBVmUwAAAAMBAAEAAAGAZC5joxYC8KMf4mAGRXUrtClR7f
sEmOxEAgRrqdJT/0pk4qPqoHeKw0dLWHNattRObEbOMzhai/I21SaOChdTmZ8hPln0/C4/
92W81zfX4cAQOWz/GG8rONS+NTG+7X4P/0mer2Zl0PASdhoqLt3FZdYzE85kg2toadmFEc
i8XZZZloqVCw05m6g6QJhS5DedpT6qrGtkNd5eevxb03m/CI2PcaI6rks+VLlXWsD/aafB
lhb1lOFjBQZOo7jdcaHJtG59+cbR9ys5p12U5fPRAPFCgJ9G2zBN+iWxbJr+V+pSBSMuLR
S/2dsJ5+y3rc6GHToU4a1tMmSwsMPFKPboMRPaiXNlMlsPojdOh2JMRW2n7PaD+VppEAKn
OHz/5I3NBeEvEvgom2WQS8ooRgIMYPaDuu4DuVErn+trYKHuyO9YBNC2hxWiD1L1ZFn+uB
jBmHIHGkv096iFW1OAj2KE50xevDH+u1YLV/UgsdfCYjoOSzg5Z9KoOOy2XHck8gABAAAA
wEQ5QpOR7n270kjT+YwRT3aMIZVmvWNyG/uTIwArF+Ycv0wYgF2gg0VbRcq616e8efWJCv
Sh0kxePsUBB0g9QoeiZtAcIjQ3HhJwqgoBOehECwfWMW1DEnMfA/IheCjcf2rrlGZrjezx
QGa12Uh5hwsCniedGI71kWal0kbb6Z7YLKsWKlmad4q2k2GWALe4j/deWq4W2bW9/pSyQx
gmSnPUCnaRZuHB45KQqkgL0WxLQBcr+FOyuRy0YrYkPP05/AAAAMEA0rFmswszka1AHhF3
AUvPDuRabnw0sGudaUuyOVwEro9Yc363MkS1yjOg96AKFbh9UF8y8VbTfhe/OF/jZQQa5C
sk1Wovd/nnz0KtzSMiri24l2h/rEqjBZIuAAW0Fzlc7iKWXqSIPv1zw6+BCQvsdiD+NG3z
Xlpo8z15LtmBl4o+d7xIFEJTXqbN+Qx+BhKi2CRf0cOYoghlqlz5qERVGOeXvWsVKaSM9m
6me4cwBECJgI1vi+tr4P+isOzJoS1TAAAAwQDR2AjBypt0VEpz8E6Jw8r01Li3Y3R+ishc
PMK3eOK4V2hJl+Bn+5ZKB18bWG21m41nOZieox6Fl0S/W0NuhAoVKrw16XsE+vp686bYmJ
EGs/Ni79R+wfaaUg8CmcoVwkcxfrFUxlszksNrws+/JvjHyMlg5S6t4VPxWO64XIcvor9O
JX0dcGsOPgrCEt37wAAWhDjbUE3dxamnC+Cr7dMVy7AgguSQgkXBMoGsleTuB1PXOusn9z
ABz8oYP1eawwEAAAAiYmhyZzNzZUBCaGFyZ2Ficy1NYWNCb29rLVByby5sb2NhbAE=
-----END OPENSSH PRIVATE KEY-----`
)

var MockPrivateKey2 = `-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----`

var Mockupgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//TODO
		return true
	},
	Subprotocols: []string{"trasa", "guacamole", "livesessions", "xterm"},
}
