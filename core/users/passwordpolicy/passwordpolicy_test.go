package passwordpolicy

import (
	"testing"
	"time"
)

func Test_checkisExpired(t *testing.T) {
	type args struct {
		thenTime int64
		timeZone string
		days     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "when timestamp is 0",
			args: args{0, "Asia/Kathmandu", "34"},
			want: true,
		},
		{
			name: "when timestamp is now",
			args: args{time.Now().Unix(), "Asia/Kathmandu", "1"},
			want: false,
		},
		{
			name: "when timestamp -100 hrs",
			args: args{time.Now().Add(-100 * time.Hour).Unix(), "Asia/Kathmandu", "1"},
			want: true,
		},
		{
			name: "invalid timestamp ",
			args: args{time.Now().Add(-100 * time.Hour).Unix(), "Nepal/athmandu", "1"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkisExpired(tt.args.thenTime, tt.args.timeZone, tt.args.days); got != tt.want {
				t.Errorf("checkisExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
