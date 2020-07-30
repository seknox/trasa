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
			args: args{"orgAdmin", "/api/v1/services/path", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"normal", "/api/v1/services/path", "POST"},
			want: false,
		},

		{
			name: "",
			args: args{"normal", "/api/v1/users", "POST"},
			want: false,
		},

		{
			name: "",
			args: args{"normal", "/api/v1/system/path", "GET"},
			want: false,
		},

		{
			name: "",
			args: args{"orgAdmin", "/api/v1/system/path", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "/api/v1/org/path", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "/api/v1/crypto/path", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "/api/v1/groups/path", "PUT"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "/api/v1/logs/path", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "/api/v1/devices/path", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "/api/v1/stats/path", "POST"},
			want: true,
		},

		{
			name: "",
			args: args{"orgAdmin", "/api/v1/gateway/path", "GET"},
			want: false,
		},

		{
			name: "",
			args: args{"orgAdmin", "/api/v1/idp/path", "POST"},
			want: false,
		},

		{
			name: "",
			args: args{"normal", "/api/v1/gateway/path", "POST"},
			want: false,
		},

		{
			name: "",
			args: args{"normalUser", "/api/v1/stats/path", "POST"},
			want: false,
		},

		{
			name: "",
			args: args{"normalUser", "/api/v1/devices/path", "POST"},
			want: false,
		},

		{
			name: "",
			args: args{"normalUser", "/api/v1/org/path", "POST"},
			want: false,
		},

		{
			name: "",
			args: args{"normalUser", "/api/v1/system/path", "POST"},
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
