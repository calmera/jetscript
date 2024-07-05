package lang

import (
	"encoding/json"
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/redpanda-data/benthos/v4/public/bloblang"
	"github.com/redpanda-data/benthos/v4/public/service"
)

func Process(exec *bloblang.Executor, msg *nats.Msg) ([]*nats.Msg, error) {
	res, err := ToBenthosMessage(msg).BloblangQueryValue(exec)
	if err != nil {
		if !errors.Is(err, bloblang.ErrRootDeleted) {
			return nil, err
		}

		return nil, nil
	}

	if res == nil {
		return nil, nil
	}

	switch res := res.(type) {
	case *nats.Msg:
		return []*nats.Msg{res}, nil
	default:
		return nil, errors.New("expected result to be a NATS message")
	}
}

func ToBenthosMessage(msg *nats.Msg) *service.Message {
	result := service.NewMessage(msg.Data)

	result.MetaSet("nats_subject", msg.Subject)

	for k, v := range msg.Header {
		result.MetaSet(k, v[0])
	}

	return result
}

func FromBenthosMessage(msg *service.Message, defaultSubject string) (*nats.Msg, error) {
	sub, fnd := msg.MetaGet("nats_subject")
	if !fnd {
		sub = defaultSubject
	}

	b, err := msg.AsBytes()
	if err != nil {
		return nil, err
	}

	result := nats.NewMsg(sub)
	result.Data = b

	err = msg.MetaWalk(func(k, v string) error {
		result.Header.Set(k, v)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

type MessageWrapper struct {
	msg *nats.Msg
}

func WrapMessage(msg *nats.Msg) MessageWrapper {
	return MessageWrapper{msg: msg}
}

// MetaSetMut sets the value of a metadata key to any value.
func (w *MessageWrapper) MetaSetMut(key string, value any) {
	switch v := value.(type) {
	case string:
		w.msg.Header.Set(key, v)
	default:
		// not really happy with eating up this error, but it's the only way
		b, _ := json.Marshal(value)
		w.msg.Header.Set(key, string(b))
	}
}

// MetaGetStr returns a metadata value if a key exists as a string, otherwise an
// empty string.
func (w *MessageWrapper) MetaGetStr(key string) string {
	return w.msg.Header.Get(key)
}

// MetaGetMut returns a metadata value if a key exists.
func (w *MessageWrapper) MetaGetMut(key string) (any, bool) {
	res := w.msg.Header.Get(key)
	if res == "" {
		return nil, false
	}

	return res, true
}

// MetaDelete removes the value of a metadata key.
func (w *MessageWrapper) MetaDelete(key string) {
	w.msg.Header.Del(key)
}

// MetaIterMut iterates each metadata key/value pair.
func (w *MessageWrapper) MetaIterMut(f func(k string, v any) error) error {
	for k, v := range w.msg.Header {
		err := f(k, v[0])
		if err != nil {
			return err
		}
	}

	return nil
}

// MetaIterStr iterates each metadata key/value pair with the value serialised
// as a string.
func (w *MessageWrapper) MetaIterStr(f func(k, v string) error) error {
	for k, v := range w.msg.Header {
		err := f(k, v[0])
		if err != nil {
			return err
		}
	}

	return nil
}
