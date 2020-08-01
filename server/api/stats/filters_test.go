package stats

import (
	"github.com/huandu/go-sqlbuilder"
	"testing"
	"time"
)

func Test_addTimeFilter(t *testing.T) {
	type args struct {
		timeFilter string
		now        time.Time
		sb         sqlbuilder.SelectBuilder
	}

	//now:=time.Unix(1257894000,0)//2009-11-09 18:15:00
	now := time.Date(2009, 11, 10, 18, 15, 0, 0, time.UTC)

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("col1").From("table1")
	tests := []struct {
		name     string
		args     args
		wantSQL  string
		wantArgs []interface{}
	}{
		{
			name:     "Today",
			args:     args{"Today", now, *sb},
			wantSQL:  "SELECT col1 FROM table1 WHERE login_time > ?",
			wantArgs: []interface{}{int64(1257811200000000000)},
		},
		{
			name:     "Yesterday",
			args:     args{"Yesterday", now, *sb},
			wantSQL:  "SELECT col1 FROM table1 WHERE login_time < ? AND login_time > ?",
			wantArgs: []interface{}{int64(1257811200000000000), int64(1257724800000000000)},
		},
		{
			name:     "7 Days",
			args:     args{"7 Days", now, *sb},
			wantSQL:  "SELECT col1 FROM table1 WHERE login_time > ?",
			wantArgs: []interface{}{int64(1257206400000000000)},
		},
		{
			name:     "30 Days",
			args:     args{"30 Days", now, *sb},
			wantSQL:  "SELECT col1 FROM table1 WHERE login_time > ?",
			wantArgs: []interface{}{int64(1255219200000000000)},
		},
		{
			name:     "90 Days",
			args:     args{"90 Days", now, *sb},
			wantSQL:  "SELECT col1 FROM table1 WHERE login_time > ?",
			wantArgs: []interface{}{int64(1250035200000000000)},
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addTimeFilter(tt.args.timeFilter, tt.args.now, tt.args.sb)
			gotSqlStr, gotArgs := got.Build()
			if gotSqlStr != tt.wantSQL {
				t.Errorf("addTimeFilter() = %v, want %v", gotSqlStr, tt.wantSQL)
			}
			if len(gotArgs) != len(tt.wantArgs) {
				t.Fatalf("incorrect number of args  got= %v, want %v", len(gotArgs), len(tt.wantArgs))
			}
			for i, a := range gotArgs {
				if gotArgs[i] != tt.wantArgs[i] {
					t.Errorf("incorrect  args  got= %v, want %v", (a), (tt.wantArgs[i]))
				}
			}

		})
	}
}
