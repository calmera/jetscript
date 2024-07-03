package lang

type (
	Environment interface {
		Parse(script string) (Executor, error)
		DumpComponents(outdir string) error
	}
	Executor interface {
		Query(val any) (any, error)
		Overlay(val any, onto *any) error
	}
)
