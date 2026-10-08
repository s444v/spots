package signaling

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestDecodeJoin(t *testing.T) {
	in := []byte(`{"type":"join","room":"lobby","name":"stepan"}`)
	var msg wsMessage
	if err := json.Unmarshal(in, &msg); err != nil {
		t.Fatalf("unmarshal %v", err)
	}
	if msg.Type != "join" {
		t.Errorf("Type = %q, want %q", msg.Type, "join")
	}
	if msg.Room != "lobby" {
		t.Errorf("Room = %q, want %q", msg.Room, "lobby")
	}
	if msg.Name != "stepan" {
		t.Errorf("Name = %q, want %q", msg.Name, "stepan")
	}
}
func TestDataRoundTripByteForByte(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "sdp offer",
			data: `{"sdp":"v=0\r\no=- 46117 2 IN IP4 127.0.0.1\r\n","type":"offer"}`,
		},
		{
			name: "ice candidate",
			data: `{"candidate":"candidate:1 1 UDP 2122252543 192.168.0.2 54321 typ host","sdpMid":"0","sdpMLineIndex":0}`,
		},
		{
			name: "order of keys is preserved",
			data: `{"z":1,"a":2,"m":3}`,
		},
		{
			name: "numbers are not normalized",
			data: `{"a":1.0,"b":1e3,"c":12345678901234567890,"d":-0}`,
		},
		{
			name: "unicode escapes are preserved",
			data: `{"s":"\u041f\u0440\u0438\u0432\u0435\u0442","emoji":"\ud83d\ude00"}`,
		},
		{
			name: "nested structures",
			data: `{"a":[1,{"b":null},[true,false]],"c":{"d":{"e":"f"}}}`,
		},
		{
			name: "array",
			data: `[1,2,3]`,
		},
		{
			name: "string",
			data: `"just a string"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := []byte(`{"type":"signal","data":` + tt.data + `}`)

			var msg wsMessage
			if err := json.Unmarshal(in, &msg); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if !bytes.Equal(msg.Data, []byte(tt.data)) {
				t.Errorf("after decode: data = %s, want %s", msg.Data, tt.data)
			}

			out, err := json.Marshal(msg)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			var raw map[string]json.RawMessage
			if err := json.Unmarshal(out, &raw); err != nil {
				t.Fatalf("unmarshal output: %v", err)
			}
			if !bytes.Equal(raw["data"], []byte(tt.data)) {
				t.Errorf("after re-encode: data = %s, want %s", raw["data"], tt.data)
			}
		})
	}
}

func TestEncodePeerLeftHasNoExtraFields(t *testing.T) {
	out, err := json.Marshal(wsMessage{Type: "peer-left"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(out, &fields); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	for _, key := range []string{"room", "name", "data"} {
		if _, ok := fields[key]; ok {
			t.Errorf("unexpected field %q in %s", key, out)
		}
	}
	allowed := map[string]bool{"type": true}
	for key := range fields {
		if !allowed[key] {
			t.Errorf("extra field %q in %s", key, out)
		}
	}

	if want := `{"type":"peer-left"}`; string(out) != want {
		t.Errorf("got %s, want %s", out, want)
	}
}
