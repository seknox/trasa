package middlewares

import (
	"fmt"
	"testing"
)

func Test_permissionChecker(t *testing.T) {
	type args struct {
		assignedRole    string
		requestEndpoint string
		method          string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{"orgAdmin", "services", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"normal", "services", "POST"},
			want: false,
		},

		{
			name: "",
			args: args{"normal", "users", "POST"},
			want: false,
		},

		{
			name: "",
			args: args{"normal", "system", "GET"},
			want: false,
		},

		{
			name: "",
			args: args{"orgAdmin", "system", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "org", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "crypto", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "groups", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "logs", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "devices", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "stats", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "gateway", "GET"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "idp", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"normal", "gateway", "POST"},
			want: false,
		},

		{
			name: "",
			args: args{"normalUser", "stats", "POST"},
			want: false,
		},

		{
			name: "",
			args: args{"normalUser", "devices", "POST"},
			want: false,
		},

		{
			name: "",
			args: args{"normalUser", "org", "POST"},
			want: false,
		},

		{
			name: "",
			args: args{"normalUser", "system", "POST"},
			want: false,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s accessing %s method on %s path", tt.args.assignedRole, tt.args.method, tt.args.requestEndpoint), func(t *testing.T) {
			if got := permissionChecker(tt.args.assignedRole, tt.args.requestEndpoint, tt.args.method); got != tt.want {
				t.Errorf("permissionChecker() = %v, want %v", got, tt.want)
			}
		})
	}
}
