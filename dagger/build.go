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
	DaggerClient      *dagger.Client
	ctx               context.Context
	ctxCancel         func()
	document          *glint.Document
	documentCloseFunc func()
	oses              []string
	arches            []string
	uiComponents      []glint.Component
	hasError          bool
	stringRenderer    *glint.StringRenderer
	hasTTY            bool
	modCache          *dagger.CacheVolume
}

func NewBuilder(oses, arches []string, withTTY bool) (*Builder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(ioutil.Discard))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("Error connecting to Dagger Engine: %w", err)
	}

	return &Builder{client, ctx, cancel, nil, nil, oses, arches, nil, false, nil, withTTY, nil}, nil
}

// WithModCache adds a go module cache to the container
func (b *Builder) WithModCache(c *dagger.Container) *dagger.Container {
	if b.modCache == nil {
		b.modCache = b.DaggerClient.CacheVolume("modcache")
	}

	return c.WithMountedCache("/go/pkg/mod", b.modCache)
}

func (b *Builder) newDocument() {
	d := glint.New()
	b.document = d

	if !b.hasTTY {
		sr := &glint.StringRenderer{}
		d.SetRenderer(sr)
		b.stringRenderer = sr
	}

	// close function should be called for all documents
	b.documentCloseFunc = func() {
		if b.document == nil {
			return
		}

		b.document.Close()
		if !b.hasTTY {
			fmt.Println(b.stringRenderer.Builder.String())
			b.document = nil
		}
	}
}

func (b *Builder) HasError() bool {
	return b.hasError
}

// Logs the start of a new build section
func (b *Builder) LogStartSectionWithContainer(message string, container *dagger.Container) func() {
	// create a new Glint document for the output
	b.newDocument()

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
		// if we have a container try and get std out and err
		if container != nil && b.document != nil {

			out, _ := container.Stderr(b.ctx)
			if out != "" {
				b.uiComponents = append(b.uiComponents,
					glint.Layout(glint.Text(out)).MarginLeft(4))
			}

			out, _ = container.Stdout(b.ctx)
			if out != "" {
				b.uiComponents = append(b.uiComponents,
					glint.Layout(glint.Text(out)).MarginLeft(4))
			}

			// update the output to show the container logs
			b.document.Set(b.uiComponents...)
		}

		// close the document and flush any output
		b.documentCloseFunc()
	}
}

func (b *Builder) LogStartSection(message string) func() {
	return b.LogStartSectionWithContainer(message, nil)
}

func (b *Builder) LogSubSection(message string) func(message string) {
	return b.LogSubSectionWithContainer(message, nil)
}

func (b *Builder) LogSubSectionWithContainer(message string, container *dagger.Container) func(message string) {
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

		outComponents := []glint.Component{
			glint.Layout(glint.Text(message)).MarginLeft(4),
		}

		// if we have a container try and get std out and err
		if container != nil {
			out, _ := container.Stderr(b.ctx)
			if out != "" {
				outComponents = append(outComponents,
					glint.Layout(glint.Text(out)).MarginLeft(8))
			}

			out, _ = container.Stdout(b.ctx)
			if out != "" {
				outComponents = append(outComponents,
					glint.Layout(glint.Text(out)).MarginLeft(8))
			}
		}

		// add the new component
		b.uiComponents = append(
			b.uiComponents, outComponents...,
		)

		b.document.Set(b.uiComponents...)
	}
}

func (b *Builder) LogError(message string, err error) {
	b.hasError = true

	// cancel the context
	b.ctxCancel()

	// clear the existing renderer
	b.documentCloseFunc()

	// create a new Glint document for the output
	b.newDocument()

	// output the text
	b.document.Set(glint.Layout(
		glint.Style(
			glint.Layout(
				glint.Layout(glint.Text("Error:")),
				glint.Layout(glint.Text(message)).MarginLeft(1),
			).Row(), glint.Color("red")),
		glint.Layout(
			glint.Text(err.Error()),
		).Row().MarginLeft(4),
	))

	go b.document.Render(context.Background())
	b.documentCloseFunc()
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
