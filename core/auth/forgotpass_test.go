package auth

import (
	"testing"

	"github.com/seknox/trasa/consts"
	"github.com/seknox/trasa/core/notif"
	"github.com/seknox/trasa/core/redis"
	"github.com/seknox/trasa/models"
)

func Test_forgotPassTfaResp(t *testing.T) {
	redis.InitStoreMock()
	notifstore := notif.InitStoreMock()

	type args struct {
		userDetails models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{models.User{
				ID:         "123",
				OrgID:      "abc",
				UserName:   "test",
				FirstName:  "F",
				MiddleName: "M",
				LastName:   "Last",
				Email:      "user@example",
			}},
			wantErr: false,
		},

		{
			name: "",
			args: args{models.User{
				ID:         "123asdsa",
				OrgID:      "abc555vassad",
				UserName:   "test",
				FirstName:  "F",
				MiddleName: "M",
				LastName:   "Last",
				Email:      "user@example",
			}},
			wantErr: false,
		},

		// TODO: Add test cases.
	}

	for _, t := range tests {
		notifstore.On("SendEmail", t.args.userDetails.OrgID, consts.EMAIL_USER_CRUD).Return(nil)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := forgotPassTfaResp(tt.args.userDetails); (err != nil) != tt.wantErr {
				t.Errorf("forgotPassTfaResp() error = %v, wantErr %v", err, tt.wantErr)
			}

		})

	}
	notifstore.AssertExpectations(t)

}
