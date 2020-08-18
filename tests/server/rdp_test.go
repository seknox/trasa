package server_test

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/seknox/guacamole"
	"github.com/seknox/trasa/server/accessproxy/rdpproxy"
	"github.com/seknox/trasa/server/models"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeWS(t *testing.T) {

	type args struct {
		params models.ConnectionParams
	}
	tests := []struct {
		name        string
		args        args
		wantStatus  bool
		wantErrMsg  string
		wantErrCode string
	}{
		{
			name: "should fail when hostname is incorrect",
			args: args{models.ConnectionParams{
				TotpCode:  "",
				Privilege: "root",
				Password:  "root",
				OptHeight: 1500,
				OptWidth:  1500,
				Hostname:  "127.0.3.1",
			}},
			wantErrMsg:  "",
			wantErrCode: "3339",
			wantStatus:  false,
		},

		{
			name: "should pass",
			args: args{models.ConnectionParams{
				TotpCode:  getTotpCode(totpSEC),
				Privilege: "root",
				Password:  "root",
				OptHeight: 1500,
				OptWidth:  1500,
				Hostname:  "127.0.0.1:33899",
			}},
			wantErrMsg:  "",
			wantErrCode: "3339",
			wantStatus:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErrCode, _, gotStatus := connectGuac(t, &tt.args.params)
			if tt.wantStatus != gotStatus {
				t.Errorf("connectGuac() status = %t, wantErr %t", gotStatus, tt.wantStatus)
				return
			}
			if tt.wantErrCode != gotErrCode {
				t.Errorf("connectGuac() errCode = %s, wantErrCode %s", gotErrCode, tt.wantErrCode)
				return
			}

			//TODO verify errMsg expectations

			//if tt.wantErrMsg!=gotErrMsg {
			//	t.Errorf("connectGuac() status = %v, wantErr %v", gotStatus, tt.wantStatus)
			//	return
			//}

		})
	}
}

func connectGuac(t *testing.T, params *models.ConnectionParams) (err_code, err_msg string, status bool) {
	rdpProxy := rdpproxy.NewProxy()

	s := httptest.NewServer(AddTestUserContextWS(rdpProxy.ServeWS))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	paramBytes, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	ws.WriteMessage(websocket.TextMessage, paramBytes)

	inst := waitForErrorOrTFA(t, ws)

	if inst == nil {
		t.Fatal(`did not expected tfa or error response`)
	}

	switch inst.Opcode {
	case guacamole.TfaOpcode:
		totpResp := guacamole.NewInstruction(
			guacamole.TfaOpcode,
			params.TotpCode,
		)
		ws.WriteMessage(websocket.TextMessage, totpResp.Byte())
		break

	case "error":
		if len(inst.Args) != 2 {
			t.Fatalf(`invalid guacamole error instruction: %s`, inst.String())
		}
		return inst.Args[1], inst.Args[0], false

	default:
		t.Fatalf(`error in test code`)
	}

	inst = waitForErrorOrTFA(t, ws)

	if inst == nil {
		return "", "", true
	}

	if inst.Opcode == "error" {
		if len(inst.Args) != 2 {
			t.Fatalf(`invalid guacamole error instruction: %s`, inst.String())
		}
		return inst.Args[0], inst.Args[1], false
	}

	t.Fatalf(`unexpected response: %s`, inst.String())
	return "", "", false

}

func waitForErrorOrTFA(t *testing.T, ws *websocket.Conn) *guacamole.Instruction {
	//Wait for tfa response
	for i := 0; i < 100; i++ {
		_, b, err := ws.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}

		inst, err := guacamole.Parse(b)
		if err != nil {
			t.Fatalf(`guacamole instruction parse error: %v`, err)
		}

		if inst.Opcode == guacamole.TfaOpcode || inst.Opcode == "error" {
			return inst
		}

	}
	return nil
}
