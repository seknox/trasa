package utils

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"reflect"
	"testing"

	"github.com/tstranex/u2f"
	"gopkg.in/square/go-jose.v2"
)

func TestAESDecryptHexString(t *testing.T) {
	type args struct {
		key     []byte
		message string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AESDecryptHexString(tt.args.key, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("AESDecryptHexString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AESDecryptHexString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeBase64(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeBase64(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeBase64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAESDecrypt(t *testing.T) {
	type args struct {
		key     []byte
		message []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AESDecrypt(tt.args.key, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("AESDecrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AESDecrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeducerAndDecryptor(t *testing.T) {
	type args struct {
		shards     [][]byte
		secretData string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeducerAndDecryptor(tt.args.shards, tt.args.secretData)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeducerAndDecryptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeducerAndDecryptor() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeBase64(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeBase64(tt.args.buf); got != tt.want {
				t.Errorf("EncodeBase64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodePrivateKeyToPEM(t *testing.T) {
	type args struct {
		privateKey *rsa.PrivateKey
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodePrivateKeyToPEM(tt.args.privateKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodePrivateKeyToPEM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAESEncrypt(t *testing.T) {
	type args struct {
		key     []byte
		message []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AESEncrypt(tt.args.key, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("AESEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AESEncrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncryptorAndSharder(t *testing.T) {
	type args struct {
		secretData string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := EncryptorAndSharder(tt.args.secretData)
			if got != tt.want {
				t.Errorf("EncryptorAndSharder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("EncryptorAndSharder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAESGenKey(t *testing.T) {
	tests := []struct {
		name    string
		want    *[KeySize]byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AESGenKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("AESGenKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AESGenKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAESGenNonce(t *testing.T) {
	tests := []struct {
		name    string
		want    *[NonceSize]byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AESGenNonce()
			if (err != nil) != tt.wantErr {
				t.Errorf("AESGenNonce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AESGenNonce() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratePrivateKey(t *testing.T) {
	type args struct {
		bitSize int
	}
	tests := []struct {
		name    string
		args    args
		want    *rsa.PrivateKey
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePrivateKey(tt.args.bitSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GeneratePrivateKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratePublicKey(t *testing.T) {
	type args struct {
		privatekey *rsa.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePublicKey(tt.args.privatekey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GeneratePublicKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEcdsaPublicKeyBytes(t *testing.T) {
	type args struct {
		pub *ecdsa.PublicKey
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEcdsaPublicKeyBytes(tt.args.pub); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEcdsaPublicKeyBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEcdsaPublicKeyFromBytes(t *testing.T) {
	type args struct {
		r   *u2f.Registration
		pub []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetEcdsaPublicKeyFromBytes(tt.args.r, tt.args.pub); (err != nil) != tt.wantErr {
				t.Errorf("GetEcdsaPublicKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHoldVaultRootKryShards(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestLoadJSONWebKey(t *testing.T) {
	type args struct {
		json []byte
		pub  bool
	}
	tests := []struct {
		name    string
		args    args
		want    *jose.JSONWebKey
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadJSONWebKey(tt.args.json, tt.args.pub)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadJSONWebKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadJSONWebKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadPrivateKey(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadPrivateKey(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadPrivateKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadPublicKey(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadPublicKey(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadPublicKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNaclDeCrypt(t *testing.T) {
	type args struct {
		encryptedData string
		decryptionKey []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NaclDeCrypt(tt.args.encryptedData, tt.args.decryptionKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NaclDeCrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NaclDeCrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNaclEnCrypt(t *testing.T) {
	type args struct {
		secretData    string
		secretkeyByte []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NaclEnCrypt(tt.args.secretData, tt.args.secretkeyByte); got != tt.want {
				t.Errorf("NaclEnCrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShamirDeducer(t *testing.T) {
	type args struct {
		keys [][]byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ShamirDeducer(tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShamirDeducer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShamirDeducer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShamirSharder(t *testing.T) {
	type args struct {
		key       []byte
		shards    int
		threshold int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShamirSharder(tt.args.key, tt.args.shards, tt.args.threshold); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShamirSharder() = %v, want %v", got, tt.want)
			}
		})
	}
}
