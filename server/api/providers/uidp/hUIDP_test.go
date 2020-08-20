package uidp

import (
	"reflect"
	"testing"

	"github.com/seknox/trasa/server/models"
)

func TestPreConfiguredIdps(t *testing.T) {
	type args struct {
		idp models.IdentityProvider
		uc  models.UserContext
	}
	tests := []struct {
		name string
		args args
		want models.IdentityProvider
	}{
		{
			name: "",
			args: args{},
			want: models.IdentityProvider{},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PreConfiguredIdps(tt.args.idp, tt.args.uc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PreConfiguredIdps() = %v, want %v", got, tt.want)
			}
		})
	}
}
