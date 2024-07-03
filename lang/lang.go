package lang

type (
	Environment interface {
		Parse(script string) (Executor, error)
	}
	Executor interface {
		Query(val any) (any, error)
		Overlay(val any, onto *any) error
	}
)
