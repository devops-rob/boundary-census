package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"dagger.io/dagger"
	"github.com/mitchellh/go-glint"
	gc "github.com/mitchellh/go-glint/components"
)

type Builder struct {
	DaggerClient *dagger.Client
	ctx          context.Context
	ctxCancel    func()
	document     *glint.Document
	oses         []string
	arches       []string
	uiComponents []glint.Component
	hasError     bool
}

func NewBuilder(oses, arches []string) (*Builder, error) {
	d := glint.New()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(ioutil.Discard))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("Error connecting to Dagger Engine: %w", err)
	}

	return &Builder{client, ctx, cancel, d, oses, arches, nil, false}, nil
}

func (b *Builder) HasError() bool {
	return b.hasError
}

// Logs the start of a new build section
func (b *Builder) LogStartSection(message string) func() {
	b.uiComponents = []glint.Component{}
	b.uiComponents = append(b.uiComponents,
		glint.Style(
			glint.Layout(
				gc.Spinner(),
				glint.Layout(glint.Text(message)).MarginLeft(1),
				glint.Layout(gc.Stopwatch(time.Now())).MarginLeft(1),
			).Row(), glint.Color("green")),
	)

	b.document.Set(b.uiComponents...)
	go b.document.Render(context.Background())

	return func() {
		b.document.Close()
	}
}

func (b *Builder) LogSubSection(message string) func(message string) {
	lc := glint.Layout(
		glint.Text(message),
		glint.Layout(gc.Stopwatch(time.Now())).MarginLeft(1),
	).Row().MarginLeft(4)

	b.uiComponents = append(b.uiComponents, lc)
	b.document.Set(b.uiComponents...)

	return func(message string) {
		for i, c := range b.uiComponents {
			if c == lc {
				b.uiComponents = append(b.uiComponents[:i], b.uiComponents[i+1:]...)
				break
			}
		}

		b.uiComponents = append(b.uiComponents, glint.Layout(glint.Text(message)).MarginLeft(4))

		b.document.Set(b.uiComponents...)
	}
}

func (b *Builder) LogError(message string, err error) {
	b.hasError = true

	// cancel the context
	b.ctxCancel()

	// clear the existing renderer
	b.document.Close()

	// create a new renderer
	b.document = glint.New()

	// output the text
	b.document.Set(glint.Layout(
		glint.Style(
			glint.Layout(
				glint.Layout(glint.Text("Error")).MarginLeft(1),
				glint.Layout(glint.Text(message)).MarginLeft(1),
			).Row(), glint.Color("red")),
		glint.Layout(
			glint.Text(err.Error()),
		).Row().MarginLeft(4),
	))

	go b.document.Render(context.Background())
	b.document.Close()
}

func (b *Builder) WithArchitectures(work func(os, arch string) error) {
	wg := sync.WaitGroup{}

	for _, goos := range b.oses {
		for _, goarch := range b.arches {
			wg.Add(1)
			go func(goos, goarch string) {
				err := work(goos, goarch)
				if err != nil {
					b.LogError("error executing architecture", err)
				}

				wg.Done()
			}(goos, goarch)
		}
	}

	wg.Wait()
}
