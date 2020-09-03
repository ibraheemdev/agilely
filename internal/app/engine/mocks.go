package engine

import (
	"context"
	"encoding/json"
	"net/http"
)

type mockClientStateReadWriter struct {
	state mockClientState
}

type mockClientState map[string]string

func newMockClientStateRW(keyValue ...string) mockClientStateReadWriter {
	state := mockClientState{}
	for i := 0; i < len(keyValue); i += 2 {
		key, value := keyValue[i], keyValue[i+1]
		state[key] = value
	}

	return mockClientStateReadWriter{state}
}

func (m mockClientStateReadWriter) ReadState(r *http.Request) (ClientState, error) {
	return m.state, nil
}

func (m mockClientStateReadWriter) WriteState(w http.ResponseWriter, cs ClientState, evs []ClientStateEvent) error {
	var state mockClientState

	if cs != nil {
		state = cs.(mockClientState)
	} else {
		state = mockClientState{}
	}

	for _, ev := range evs {
		switch ev.Kind {
		case ClientStateEventPut:
			state[ev.Key] = ev.Value
		case ClientStateEventDel:
			delete(state, ev.Key)
		}
	}

	b, err := json.Marshal(state)
	if err != nil {
		return err
	}

	w.Header().Set("test_session", string(b))
	return nil
}

func (m mockClientState) Get(key string) (string, bool) {
	val, ok := m[key]
	return val, ok
}

type mockEmailRenderer struct{}

func (m mockEmailRenderer) Load(names ...string) error { return nil }

func (m mockEmailRenderer) Render(ctx context.Context, name string, data HTMLData) ([]byte, string, error) {
	switch name {
	case "text":
		return []byte("a development text e-mail template"), "text/plain", nil
	case "html":
		return []byte("a development html e-mail template"), "text/html", nil
	default:
		panic("shouldn't get here")
	}
}

type mockLogger struct{}

func (m mockLogger) Info(s string)  {}
func (m mockLogger) Error(s string) {}
