package utils

import (
	"reflect"
	"testing"

	"github.com/seknox/trasa/models"
)

func TestArrayContainsString(t *testing.T) {
	type args struct {
		s []string
		e string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{[]string{"apple", "ball", "cat", "dog"}, "dog"},
			want: true,
		},
		{
			name: "",
			args: args{[]string{"apple", "ball", "cat", "dog"}, "rat"},
			want: false,
		},
		{
			name: "",
			args: args{[]string{"", "ball", "cat", "dog"}, "rat"},
			want: false,
		},
		{
			name: "",
			args: args{[]string{"", "ball", "cat", "dog"}, ""},
			want: true,
		},

		{
			name: "",
			args: args{[]string{"test", "ball", "cat", "dog"}, ""},
			want: false,
		},
		{
			name: "",
			args: args{[]string{"test", "ball", "", ""}, "uu"},
			want: false,
		},
		{
			name: "",
			args: args{[]string{}, "uu"},
			want: false,
		},
		{
			name: "",
			args: args{[]string{"uu"}, "uu"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayContainsString(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("ArrayContainsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainFromEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{"name@example.com"},
			want: "example.com",
		},
		{
			name: "",
			args: args{"name@localhost"},
			want: "localhost",
		},
		{
			name: "",
			args: args{"1@2"},
			want: "2",
		},
		{
			name: "",
			args: args{""},
			want: "",
		},
		{
			name: "",
			args: args{"@"},
			want: "",
		},
		{
			name: "",
			args: args{"name@"},
			want: "",
		},
		{
			name: "",
			args: args{"@local"},
			want: "local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DomainFromEmail(tt.args.email)
			if got != tt.want {
				t.Errorf("DomainFromEmail() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestNormalizeString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{"sTb asjd"},
			want: "stb asjd",
		},
		{
			name: "",
			args: args{" ab "},
			want: "ab",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeString(tt.args.s); got != tt.want {
				t.Errorf("NormalizeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToStringArr(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "",
			args:    args{[]interface{}{"hello", "%^&", "123"}},
			want:    []string{"hello", "%^&", "123"},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{[]interface{}{"hello", "%^&", models.Policy{}}},
			want:    []string{"hello", "%^&"},
			wantErr: true,
		},

		{
			name:    "",
			args:    args{[]interface{}{"hello", "%^&", nil}},
			want:    []string{"hello", "%^&"},
			wantErr: true,
		},
		{
			name:    "",
			args:    args{[]interface{}{nil, "%^&", nil}},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToStringArr(tt.args.vals)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToStringArr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToStringArr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
