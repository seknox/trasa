package stats

import (
	"testing"
)

func Test_sortIps1(t *testing.T) {

	//TODO
	arr := []ipcount{
		{10, "194.168.0.1"},
		{200, "192.168.0.1"},
		{20, "192.168.0.1"},
		{40, "192.188.0.1"},
		{20, "192.168.0.1"},
		{20, "192.168.0.1"},
		{1, "192.168.0.1"},
		{20, "192.168.0.1"},
		{20, "192.168.2.1"},
		{20, "192.168.2.11"},
		{24, "192.168.2.5"},
		{22, "192.168.2.5"},
		{6, "192.166.0.1"},
		{12, "192.166.0.6"},
		{11, "192.166.0.9"},
		{3, "172.168.0.1"},
		{30, "10.168.0.1"},
	}

	t.Run("", func(t *testing.T) {
		got := sortIps(arr)

		if len(got) != 4 {
			t.Errorf("incorrect number of first octets: want %d, got %d", 4, len(got))
		}

		if got[0].Key < got[1].Key {
			t.Errorf("not sorted: want %d, got %d", 4, len(got))
		}
		//if(got[2].Children[0].Children[0].Key <got[2].Children[0].Children[1].Key){
		//	t.Errorf("not sorted: want %d, got %d",4,len(got))
		//}

		//TODO

	})

}
