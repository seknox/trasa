package utils

import "testing"

func TestCalculateTotp(t *testing.T) {
	type args struct {
		dbcode string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := CalculateTotp(tt.args.dbcode)
			if got != tt.want {
				t.Errorf("CalculateTotp() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CalculateTotp() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("CalculateTotp() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
