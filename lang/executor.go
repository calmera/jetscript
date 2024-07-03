package lang

import "github.com/redpanda-data/benthos/v4/public/bloblang"

func wrapExecuter(be *bloblang.Executor) Executor {
	return &blobExecutor{
		be: be,
	}
}

type blobExecutor struct {
	be *bloblang.Executor
}

func (e *blobExecutor) Overlay(val any, onto *any) error {
	return e.be.Overlay(val, onto)
}

func (e *blobExecutor) Query(val any) (any, error) {
	return e.be.Query(val)
}
