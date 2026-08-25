package serial

import (
	"encoding/json"
	"fmt"
	"time"
)

type Envelope struct {
	Version int             `json:"version"`
	Type    string          `json:"type"`
	At      time.Time       `json:"at"`
	Data    json.RawMessage `json:"data"`
}

func Encode(typ string, v any) ([]byte, error) {
	d, e := json.Marshal(v)
	if e != nil {
		return nil, e
	}
	return json.Marshal(Envelope{Version: 1, Type: typ, At: time.Now().UTC(), Data: d})
}
func Decode(b []byte, out any) error {
	var e Envelope
	if err := json.Unmarshal(b, &e); err != nil {
		return err
	}
	if e.Version != 1 {
		return fmt.Errorf("unsupported envelope version %d", e.Version)
	}
	return json.Unmarshal(e.Data, out)
}
