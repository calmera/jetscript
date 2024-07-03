package nats

import (
	"github.com/redpanda-data/benthos/v4/public/bloblang"
)

type Header struct {
	Key    string
	Values []string
}

func init() {
	env := bloblang.GlobalEnvironment()
	if err := AttachHeader(env); err != nil {
		panic(err)
	}
}

func headerPluginSpec() *bloblang.PluginSpec {
	return bloblang.NewPluginSpec().
		Description("Construct a header. Values can only be strings").
		Category("NATS").
		Param(bloblang.NewStringParam("key").Description("The header key")).
		Param(bloblang.NewStringParam("value").Description("The header value"))
}

func headerCtor(args *bloblang.ParsedParams) (bloblang.Function, error) {
	key, err := args.GetString("key")
	if err != nil {
		return nil, err
	}

	value, err := args.GetString("value")
	if err != nil {
		return nil, err
	}

	return func() (any, error) {
		return Header{Key: key, Values: []string{value}}, nil
	}, nil
}

func AttachHeader(env *bloblang.Environment) error {
	return env.RegisterFunctionV2("header", asMsgPluginSpec(), headerCtor)
}
