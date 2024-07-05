package lang

import "github.com/redpanda-data/benthos/v4/public/bloblang"

type (
	Environment interface {
		Parse(script string) (*bloblang.Executor, error)
		DumpComponents(outdir string) error
	}
)
