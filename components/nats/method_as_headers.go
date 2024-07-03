package nats

import (
	"fmt"
	"github.com/redpanda-data/benthos/v4/public/bloblang"
)

func init() {
	env := bloblang.GlobalEnvironment()
	if err := AttachAsHeaders(env); err != nil {
		panic(err)
	}
}

func asHeadersPluginSpec() *bloblang.PluginSpec {
	return bloblang.NewPluginSpec().
		Description(`Convert the incoming map to headers`).
		Category("NATS")
}

func asHeadersCtor(args *bloblang.ParsedParams) (bloblang.Method, error) {
	return func(data any) (any, error) {
		if data == nil {
			return nil, nil
		}

		var headersArray []Header
		switch res := data.(type) {
		case map[string]string:
			for k, v := range res {
				headersArray = append(headersArray, Header{Key: k, Values: []string{v}})
			}
		case map[string][]string:
			for k, v := range res {
				headersArray = append(headersArray, Header{Key: k, Values: v})
			}
		case map[string]interface{}:
			for k, v := range res {
				s, err := encodeToString(v)
				if err != nil {
					return nil, fmt.Errorf("failed to encode entry %s within the provided data; %w", k, err)
				}
				headersArray = append(headersArray, Header{Key: k, Values: []string{s}})
			}
		default:
			return nil, fmt.Errorf("expected data to be a map[string]string, map[string][]string or map[string]any")
		}

		return headersArray, nil
	}, nil
}

func AttachAsHeaders(env *bloblang.Environment) error {
	return env.RegisterMethodV2("as_headers", asHeadersPluginSpec(), asHeadersCtor)
}
