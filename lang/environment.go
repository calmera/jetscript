package lang

import (
	"encoding/json"
	"fmt"
	"github.com/redpanda-data/benthos/v4/public/bloblang"
	"github.com/sirupsen/logrus"
	"os"
)

func NewEnvironment() Environment {
	return &blobEnvironment{
		be: bloblang.GlobalEnvironment(),
	}
}

type blobEnvironment struct {
	be *bloblang.Environment
}

func (e *blobEnvironment) Parse(script string) (*bloblang.Executor, error) {
	be, err := e.be.Parse(script)
	if err != nil {
		return nil, err
	}
	return be, nil
}

func (e *blobEnvironment) DumpComponents(outdir string) error {
	index := map[string]summary{}

	// create the output dir if it does not exist
	if err := os.MkdirAll(outdir, 0755); err != nil {
		return err
	}

	e.be.WalkFunctions(func(name string, spec *bloblang.FunctionView) {
		key := name
		td := spec.TemplateData()
		index[key] = summary{
			Name:       td.Name,
			Kind:       KindFunction,
			Categories: []string{td.Category},
			Status:     td.Status,
		}

		// write the function to a file
		f, err := os.Create(outdir + "/" + key + ".json")
		if err != nil {
			logrus.Error(err)
			return
		}

		b, err := spec.FormatJSON()
		if err != nil {
			logrus.Error(err)
			return
		}

		if _, err := f.Write(b); err != nil {
			logrus.Error(err)
			return
		}
	})

	e.be.WalkMethods(func(name string, spec *bloblang.MethodView) {
		key := name
		td := spec.TemplateData()
		var cats []string
		for _, cat := range td.Categories {
			cats = append(cats, cat.Category)
		}

		index[key] = summary{
			Name:       td.Name,
			Kind:       KindMethod,
			Categories: cats,
			Status:     td.Status,
		}

		// write the function to a file
		f, err := os.Create(outdir + "/" + key + ".json")
		if err != nil {
			logrus.Error(err)
			return
		}

		b, err := spec.FormatJSON()
		if err != nil {
			logrus.Error(err)
			return
		}

		if _, err := f.Write(b); err != nil {
			logrus.Error(err)
			return
		}
	})

	// write the index to a file
	f, err := os.Create(outdir + "/__index.json")
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}

	b, err := json.Marshal(index)
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}

	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	return nil
}

type Kind string

var (
	KindFunction Kind = "function"
	KindMethod   Kind = "method"
)

type summary struct {
	Name       string   `json:"name"`
	Kind       Kind     `json:"kind"`
	Categories []string `json:"categories"`
	Status     string   `json:"status"`
}
