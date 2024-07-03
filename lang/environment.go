package lang

import "github.com/redpanda-data/benthos/v4/public/bloblang"

func NewEnvironment() Environment {
	return &blobEnvironment{
		be: bloblang.GlobalEnvironment(),
	}
}

type blobEnvironment struct {
	be *bloblang.Environment
}

func (e *blobEnvironment) Parse(script string) (Executor, error) {
	be, err := e.be.Parse(script)
	if err != nil {
		return nil, err
	}
	return wrapExecuter(be), nil
}
