package utils

import "testing"

func TestSqlReplacer(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{"SELECT * FROM users where id=?"},
			want: "SELECT * FROM users where id=$1",
		},
		{
			name: "",
			args: args{"? ? ? ? ? ? ?"},
			want: "$1 $2 $3 $4 $5 $6 $7",
		},
		{
			name: "",
			args: args{""},
			want: "",
		},
		{
			name: "",
			args: args{"??"},
			want: "$1$2",
		},

		{
			name: "",
			args: args{"?,?"},
			want: "$1,$2",
		},
		{
			name: "works with line break",
			args: args{`?
?`},
			want: `$1
$2`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SqlReplacer(tt.args.src); got != tt.want {
				t.Errorf("SqlReplacer() = %v, want %v", got, tt.want)
			}
		})
	}
}
