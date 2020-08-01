package utils

import (
	"net"
	"testing"
)

func TestGetIPFromAddr(t *testing.T) {
	type args struct {
		addr net.Addr
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "doesnt break with nil value",
			args: args{nil},
			want: "0.0.0.0",
		},
		{
			name: "",
			args: args{&net.IPAddr{IP: net.ParseIP("1.1.1.1")}},
			want: "1.1.1.1",
		},
		{
			name: "doesnt break with 0.0.0.0",
			args: args{&net.IPAddr{IP: net.ParseIP("0.0.0.0")}},
			want: "0.0.0.0",
		},
		{
			name: "doesnt break with invalid ip",
			args: args{&net.IPAddr{IP: net.ParseIP("1.2.3.")}},
			want: "0.0.0.0",
		},
		{
			name: "works with private ip",
			args: args{&net.IPAddr{IP: net.ParseIP("192.168.0.1")}},
			want: "192.168.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIPFromAddr(tt.args.addr); got != tt.want {
				t.Errorf("GetIPFromAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPrivateIP(t *testing.T) {
	type args struct {
		ip net.IP
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "doesnt break with nil value",
			args: args{nil},
			want: true,
		},
		{
			name: "",
			args: args{net.ParseIP("1.1.1.1")},
			want: false,
		},
		{
			name: "",
			args: args{net.ParseIP("127.0.0.1")},
			want: true,
		},
		{
			name: "",
			args: args{net.ParseIP("192.168.99.99")},
			want: true,
		},
		{
			name: "",
			args: args{net.ParseIP("10.10.212.2")},
			want: true,
		},
		{
			name: "",
			args: args{net.ParseIP("10.10.1.1")},
			want: true,
		},
		{
			name: "",
			args: args{net.ParseIP("192.169.99.99")},
			want: false,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPrivateIP(tt.args.ip); got != tt.want {
				t.Errorf("IsPrivateIP(%v) = %v, want %v", tt.args.ip, got, tt.want)
			}
		})
	}
}
