package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/redpanda-data/benthos/v4/public/bloblang"
)

func init() {
	env := bloblang.GlobalEnvironment()
	if err := AttachAsMsg(env); err != nil {
		panic(err)
	}
}

func asMsgPluginSpec() *bloblang.PluginSpec {
	return bloblang.NewPluginSpec().
		Description(`Construct a message to be sent to NATS.

When data is a string, it will be converted into a byte array. However if the data is a byte array it will be sent 
as is. In all other cases, it will be marshalled into a JSON byte array.
`).
		Category("NATS").
		Param(bloblang.NewStringParam("subject").Description("The subject to send the message to")).
		Param(bloblang.NewAnyParam("headers").Description("The headers of the message").Optional())
}

func asMsgCtor(args *bloblang.ParsedParams) (bloblang.Method, error) {
	subject, err := args.GetString("subject")
	if err != nil {
		return nil, err
	}

	headers, err := args.Get("headers")
	if err != nil {
		return nil, err
	}

	var headersArray []Header
	if headers != nil {
		ok := false
		headersArray, ok = headers.([]Header)
		if !ok {
			return nil, fmt.Errorf("headers must be an array of Header objects")
		}
	}

	return func(data any) (any, error) {
		encodedData, err := encodeData(data)
		if err != nil {
			return nil, fmt.Errorf("failed to encode data; %w", err)
		}

		result := nats.NewMsg(subject)
		result.Data = encodedData

		for _, h := range headersArray {
			for idx, v := range h.Values {
				if idx == 0 {
					result.Header.Set(h.Key, v)
				} else {
					result.Header.Add(h.Key, v)
				}
			}
		}

		return result, nil
	}, nil
}

func AttachAsMsg(env *bloblang.Environment) error {
	return env.RegisterMethodV2("as_msg", asMsgPluginSpec(), asMsgCtor)
}
