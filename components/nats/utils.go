package nats

import "encoding/json"

func encodeData(data any) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	switch v := data.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return json.Marshal(data)
	}
}

func encodeToString(data any) (string, error) {
	if data == nil {
		return "", nil
	}

	b, err := encodeData(data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
