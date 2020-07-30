package logs

import "testing"

func Test_getMinioPath(t *testing.T) {
	type args struct {
		sessionID   string
		sessionType string
		orgID       string
		year        string
		month       string
		day         string
	}

	type testc struct {
		name  string
		args  args
		want  string
		want1 string
	}

	tests := []testc{
		{
			name: "blank fields",
			args: args{
				sessionID:   "1234",
				sessionType: "",
				orgID:       "abc",
				year:        "1111",
				month:       "2",
				day:         "3",
			},
			want:  "abc/1111/2/3/1234.session",
			want1: "unknown",
		},
		{
			name: "ssh bucket",
			args: args{
				sessionID:   "1234567890",
				sessionType: "ssh",
				orgID:       "abcd",
				year:        "1919",
				month:       "02",
				day:         "03",
			},
			want:  "abcd/1919/02/03/1234567890.session",
			want1: "trasa-ssh-logs",
		},
		{
			name: "rdp bucket",
			args: args{
				sessionID:   "1234567890",
				sessionType: "rdp",
				orgID:       "abcd",
				year:        "1919",
				month:       "02",
				day:         "03",
			},
			want:  "abcd/1919/02/03/1234567890.guac",
			want1: "trasa-guac-logs",
		},
		{
			name: "http bucket",
			args: args{
				sessionID:   "1234567890",
				sessionType: "http",
				orgID:       "abcd",
				year:        "1919",
				month:       "02",
				day:         "03",
			},
			want:  "abcd/1919/02/03/1234567890.mp4",
			want1: "trasa-https-logs",
		},
		{
			name: "http raw bucket",
			args: args{
				sessionID:   "1234567890",
				sessionType: "http-raw",
				orgID:       "abcd",
				year:        "1919",
				month:       "02",
				day:         "03",
			},
			want:  "abcd/1919/02/03/1234567890.http-raw",
			want1: "trasa-https-logs",
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getMinioPath(tt.args.sessionID, tt.args.sessionType, tt.args.orgID, tt.args.year, tt.args.month, tt.args.day)
			if got != tt.want {
				t.Errorf("getMinioPath() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getMinioPath() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
