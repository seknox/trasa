package tfa

import "testing"

func Test_rsaVerify(t *testing.T) {
	type args struct {
		signedChallenge   string
		originalChallenge string
		publicKeyPEM      string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"should fail when signed challenge is blank",
			args{
				signedChallenge:   "",
				originalChallenge: "",
				publicKeyPEM:      "",
			},
			false,
			true,
		},
		{
			"should fail when signed challenge is incorrect",
			args{
				signedChallenge:   "ASKNKJSCJKNSCSNCIUNSCS",
				originalChallenge: "stadysvayt7",
				publicKeyPEM:      "'-----BEGIN RSA PUBLIC KEY-----\nMIICCgKCAgEAqhqiLyqj9N8JytyH/iW82ynJBSRGlj7F016s9xzTK+UpfS7xA71j\nLDq1Zx4M5U+Iv1QWaDVcmh+8wMiviCSbEuE8tVPZANT/PdObo/Zl+Zbu/24jzqiT\nlldH6dK6SYj3f6FsOHz+ort2yJU4AXXQA7nTg4p4pvtr7x8dswLiQjt0Oqzrr2g1\n5Xm5GWCfUKJKVepNi+8NN/Y0sCSiIfeYD6/ug2MnyzEZfxiohM2UYVG31rhQV1uy\nr/ZPqH+IVM89oVqT8N78fOv6R5XEXnlz1begzPVkLT9rrMdnj+NDO+ooQk5f6t/M\nuZvbay+vUzwWUPoU861z5aDLljRCEDVE/AR/h2uXCJNSj63oT+yEFDLYmz+RZ1zr\nMt92okl/9KUgC2nNfsyGqH8XoYeplQdxeGtXfhfrbrEVFFk0pZGXwWCKlvs/aGfu\nvN/LHaMdFJnB1hhw8igLM1v2YcNU9Cra5Ldo26cq//yydpbUOYgfN7jQpbBnneNx\nqYxqxLg5gjtBsM1bQJhWNyu6ZSpzKdW5jCiGBClWl3NdFbWzke60OLAZ4mSS31+U\nujGzuLSQkvP7RnLXXNGTBa0rHxxFdPGwtEQEh1RWPxKXeXErauFvmU+boqva5qvq\nMlmfQqmpRnAtwZ63poVRAXOEP/I9QUf/iPkgMsQ9xulLWLlhBpwIz4ECAwEAAQ==\n-----END RSA PUBLIC KEY-----\n",
			},
			false,
			true,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rsaVerify(tt.args.signedChallenge, tt.args.originalChallenge, tt.args.publicKeyPEM)
			if (err != nil) != tt.wantErr {
				t.Errorf("rsaVerify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("rsaVerify() got = %v, want %v", got, tt.want)
			}
		})
	}
}
