package stats

import (
	"reflect"
	"testing"
)

func Test_sortIps(t *testing.T) {
	type args struct {
		arr []ipcount
	}
	type testc struct {
		name string
		args args
		want []firstOctet
	}

	tests := []testc{
		{
			name: "Zero length array",
			args: args{[]ipcount{}},
			want: []firstOctet{},
		},

		{
			name: "Nil array",
			args: args{nil},
			want: []firstOctet{},
		},

		{
			name: "Single value",
			args: args{[]ipcount{{
				Count: 4,
				IP:    "127.0.0.1",
			}}},
			want: []firstOctet{{
				Key:   "127",
				Name:  "127.0.0.0/8",
				Value: 0,
				Children: []secondOctet{{
					Key:   "0",
					Name:  "127.0.0.0/16",
					Value: 0,
					Children: []thirdOctet{{
						Key:   "0",
						Name:  "127.0.0.0/24",
						Value: 0,
						Children: []fourthOctet{{
							Key:   "1",
							Name:  "127.0.0.1",
							Value: 4,
						}},
					}},
				}},
			}},
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortIps(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortIps() = %v, want %v", got, tt.want)
			}
		})
	}
}
