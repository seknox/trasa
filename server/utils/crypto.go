package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"strings"

	"github.com/tstranex/u2f"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ssh"
	"gopkg.in/square/go-jose.v2"

	"encoding/base64"

	"github.com/hashicorp/vault/shamir"
	"golang.org/x/crypto/nacl/secretbox"
)

func LoadJSONWebKey(json []byte, pub bool) (*jose.JSONWebKey, error) {
	var jwk jose.JSONWebKey
	err := jwk.UnmarshalJSON(json)
	if err != nil {
		return nil, err
	}
	if !jwk.Valid() {
		return nil, errors.New("invalid JWK key")
	}
	if jwk.IsPublic() != pub {
		return nil, errors.New("priv/pub JWK key mismatch")
	}
	return &jwk, nil
}

// LoadPublicKey loads a public key from PEM/DER/JWK-encoded data.
func LoadPublicKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	// Try to load SubjectPublicKeyInfo
	pub, err0 := x509.ParsePKIXPublicKey(input)
	if err0 == nil {
		return pub, nil
	}

	cert, err1 := x509.ParseCertificate(input)
	if err1 == nil {
		return cert.PublicKey, nil
	}

	jwk, err2 := LoadJSONWebKey(data, true)
	if err2 == nil {
		return jwk, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s', '%s' and '%s'", err0, err1, err2)
}

// LoadPrivateKey loads a private key from PEM/DER/JWK-encoded data.
func LoadPrivateKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	var priv interface{}
	priv, err0 := x509.ParsePKCS1PrivateKey(input)
	if err0 == nil {
		return priv, nil
	}

	priv, err1 := x509.ParsePKCS8PrivateKey(input)
	if err1 == nil {
		return priv, nil
	}

	priv, err2 := x509.ParseECPrivateKey(input)
	if err2 == nil {
		return priv, nil
	}

	jwk, err3 := LoadJSONWebKey(input, false)
	if err3 == nil {
		return jwk, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s', '%s', '%s' and '%s'", err0, err1, err2, err3)
}

// GeneratePrivateKey creates a RSA Private Key of specified byte size
func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	//log.Println("Private Key generated")
	return privateKey, nil
}

// EncodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// GeneratePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func GeneratePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}

func GetEcdsaPublicKeyBytes(pub *ecdsa.PublicKey) []byte {
	val := elliptic.Marshal(pub.Curve, pub.X, pub.Y)

	return val
}

func GetEcdsaPublicKeyFromBytes(r *u2f.Registration, pub []byte) error {
	x, y := elliptic.Unmarshal(elliptic.P256(), pub)
	if x == nil {
		return errors.New("u2f: invalid public key")
	}
	r.PubKey.Curve = elliptic.P256()
	r.PubKey.X = x
	r.PubKey.Y = y

	return nil
}

/////////////////////////////////////////////////////////////////
////////////////////		AES Encryption 		/////////////////
/////////////////////////////////////////////////////////////////
const (
	KeySize   = 32
	NonceSize = 12
)

// AESGenKey creates a new random secret key.
func AESGenKey() (*[KeySize]byte, error) {
	key := new([KeySize]byte)
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, err
	}

	return key, nil
}

// AESGenNonce creates a new random nonce.
func AESGenNonce() (*[NonceSize]byte, error) {
	nonce := new([NonceSize]byte)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

var (
	ErrEncrypt = errors.New("secret: encryption failed")
	ErrDecrypt = errors.New("secret: decryption failed")
)

// AESEncrypt computes AES GCM encryption
func AESEncrypt(key, message []byte) ([]byte, error) {
	if string(key) == "" {
		return nil, fmt.Errorf("%s", "ecryption key not found")
	}

	var buf []byte

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("enc.aes.NewCipher: %s", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("enc.cipher.NewGCM: %s", err)
	}

	nonce, err := AESGenNonce()
	if err != nil {
		return nil, fmt.Errorf("enc.AESGenNonce: %s", err)
	}

	buf = append(buf, nonce[:]...)
	buf = gcm.Seal(buf, nonce[:], message, nil)
	return buf, nil
}

// AESDecrypt computes AES GCM decryption
func AESDecrypt(key, message []byte) ([]byte, error) {
	if len(message) <= NonceSize {
		return nil, ErrDecrypt
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("dec.aes.NewCipher: %s", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("dec.cipher.NewGCM: %s", err)
	}

	nonce := make([]byte, NonceSize)
	copy(nonce, message)

	// Decrypt the message, using the sender ID as the additional
	// data requiring authentication.
	out, err := gcm.Open(nil, nonce, message[NonceSize:], nil)
	if err != nil {
		return nil, fmt.Errorf("dec.gcm.Open: %s", err)
	}
	return out, nil
}

func AESDecryptHexString(key []byte, message string) ([]byte, error) {
	msgbytes, err := hex.DecodeString(message)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex value: %s", err)
	}
	if len(msgbytes) <= NonceSize {
		return nil, ErrDecrypt
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("dec.aes.NewCipher: %s", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("dec.cipher.NewGCM: %s", err)
	}

	nonce := make([]byte, NonceSize)
	copy(nonce, msgbytes)

	// Decrypt the message, using the sender ID as the additional
	// data requiring authentication.
	out, err := gcm.Open(nil, nonce, msgbytes[NonceSize:], nil)
	if err != nil {
		return nil, fmt.Errorf("dec.gcm.Open: %s", err)
	}
	return out, nil
}

func EncryptorAndSharder(secretData string) (string, []string) {
	//fmt.Println("initial root token: ", secretData)
	SecretKeyBytes, err := hex.DecodeString(GetRandomString(32))
	if err != nil {
		logrus.Error(err)
	}

	//fmt.Println("encryption key: ", base64.StdEncoding.EncodeToString(SecretKeyBytes))
	// get encrypted key
	encryptedSecretKey := NaclEnCrypt(secretData, SecretKeyBytes)
	//fmt.Println("encrypted root token: ", encryptedSecretKey)

	// shard secret key. these shards will be required whenever users need to get root token (decrypted value)
	shardedKeys := ShamirSharder(SecretKeyBytes, 5, 3)

	return encryptedSecretKey, shardedKeys

}

func DeducerAndDecryptor(shards [][]byte, secretData string) (string, error) {

	// deduce encryption key
	deducedSecretKey, err := ShamirDeducer(shards)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	//fmt.Println("deduced encryption Key: ", deducedSecretKey)
	// we now retrieve(by decryption) encrypted data from deducedSecretKey.
	data, err := NaclDeCrypt(secretData, deducedSecretKey)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	//fmt.Println("deduced root Token: ", data)
	return data, nil
}

func NaclEnCrypt(secretData string, secretkeyByte []byte) string {

	var secretKey [32]byte
	copy(secretKey[:], secretkeyByte)

	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	encryptedData := secretbox.Seal(nonce[:], []byte(secretData), &nonce, &secretKey)

	return base64.StdEncoding.EncodeToString(encryptedData)

}

func NaclDeCrypt(encryptedData string, decryptionKey []byte) (string, error) {

	decodedEncryptedData, _ := base64.StdEncoding.DecodeString(encryptedData)

	var decryptNonce [24]byte
	copy(decryptNonce[:], decodedEncryptedData[:24])

	var decretkey [32]byte
	//decodedkey, _ := base64.StdEncoding.DecodeString(decryptionKey)
	copy(decretkey[:], decryptionKey)

	decrypted, ok := secretbox.Open(nil, decodedEncryptedData[24:], &decryptNonce, &decretkey)
	if !ok {
		return "", fmt.Errorf("error in secret decryption")
	}

	return string(decrypted), nil

}

func ShamirSharder(key []byte, shards, threshold int) []string {
	keys, err := shamir.Split(key, shards, threshold)
	if err != nil {
		logrus.Error(err)
	}

	var shardedKeys []string
	for _, v := range keys {
		str := base64.StdEncoding.EncodeToString(v)
		shardedKeys = append(shardedKeys, str)
	}

	return shardedKeys
}

func ShamirDeducer(keys [][]byte) ([]byte, error) {
	var val []byte
	val, err := shamir.Combine(keys)
	if err != nil {
		logrus.Error(err)
		return val, err
	}

	return val, nil
}

func HoldVaultRootKryShards() {

}

func EncodeBase64(buf []byte) string {
	s := base64.URLEncoding.EncodeToString(buf)
	return strings.TrimRight(s, "=")
}

func DecodeBase64(s string) ([]byte, error) {
	for i := 0; i < len(s)%4; i++ {
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

/////////////////////////////////////////////
//////  ECDH
///////////////////////////////////////////////

// ECDHGenKeyPair generated public and private key pair to be used in curve25519.ScalarBaseMult()
// This function should be called in both client and server independently.
// Reference from https://cr.yp.to/ecdh.html. Inspiration from https://github.com/aead/ecdh
func ECDHGenKeyPair() (privateKey *[32]byte, publicKey *[32]byte, err error) {
	priv := new([KeySize]byte)
	pub := new([KeySize]byte)

	_, err = io.ReadFull(rand.Reader, priv[:])
	if err != nil {
		return priv, pub, err
	}

	priv[0] &= 248
	priv[31] &= 127
	priv[31] |= 64

	curve25519.ScalarBaseMult(pub, priv)

	privateKey = priv
	publicKey = pub

	return
}

// ECDHComputeSecret takes private key, remote peers public key and computes secret key.
func ECDHComputeSecret(yourPrivateKey *[32]byte, remotePublicKey *[32]byte) (secret []byte) {

	var sec [32]byte

	curve25519.ScalarMult(&sec, yourPrivateKey, remotePublicKey)

	secret = sec[:]
	return
}
