package utils

import "testing"

func TestArrayContainsInt(t *testing.T) {
	type args struct {
		s []int
		e int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{[]int{1, 2, 3, 4, 5}, 6},
			want: false,
		},
		{
			name: "",
			args: args{[]int{1, 2, 3, 4, 5}, 5},
			want: true,
		},
		{
			name: "",
			args: args{[]int{6, 2, 3, 4, 5}, 6},
			want: true,
		},
		{
			name: "",
			args: args{[]int{}, 6},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayContainsInt(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("ArrayContainsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
