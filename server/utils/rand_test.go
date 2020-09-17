package utils

import (
	"reflect"
	"testing"
)

func TestGetRandomBytes(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "",
			args: args{0},
			want: []byte(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRandomBytes(tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRandomBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRandomID(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{0},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRandomString(tt.args.length); got != tt.want {
				t.Errorf("GetRandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUUID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUUID(); got != tt.want {
				t.Errorf("GetUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
