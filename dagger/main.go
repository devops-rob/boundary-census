package main

import (
	"fmt"
	"os"
	"path"

	"dagger.io/dagger"
)

// Dagger build pipeline
func main() {
	// remove the build folder if it exists
	os.RemoveAll("./build")

	// define build matrix
	oses := []string{"linux", "darwin"}
	arches := []string{"amd64", "arm64"}

	builder, err := NewBuilder(oses, arches)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	src := builder.DaggerClient.Host().Directory(".")
	if err != nil {
		fmt.Printf("Error getting reference to host directory: %s", err)
		os.Exit(1)
	}

	test(builder, src)
	build(builder, src)
	createDockerContainers(builder, src)

	if builder.HasError() {
		os.Exit(1)
	}
}

func test(builder *Builder, src *dagger.Directory) {
	if builder.ctx.Err() != nil {
		return
	}

	done := builder.LogStartSection("Running unit tests")

	_, err := builder.DaggerClient.Container().
		From("golang:latest").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "test", "-v", "./..."}).
		ExitCode(builder.ctx)

	if err != nil {
		builder.LogError("unable to test application", err)
	}

	defer done()
}

// build the application for multiple architectures
func build(builder *Builder, src *dagger.Directory) {
	if builder.ctx.Err() != nil {
		return
	}

	done := builder.LogStartSection("Building application for all architectures")
	defer done()

	// create a build container for go
	golang := builder.DaggerClient.Container().From("golang:latest")
	golang = golang.WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithEnvVariable("CGO_ENABLED", "0")

	builder.WithArchitectures(func(goos, goarch string) error {
		done := builder.LogSubSection(fmt.Sprintf("building %s %s", goos, goarch))

		p := path.Join("build/", goos, goarch)
		build := golang.WithEnvVariable("GOOS", goos).
			WithEnvVariable("GOARCH", goarch).
			WithExec([]string{"go", "build", "-o", path.Join(p, "census")})

		_, err := build.Directory(p).Export(builder.ctx, p)
		if err != nil {
			builder.LogError("failed to build", err)
		}

		done(fmt.Sprintf("%s %s done", goos, goarch))

		return nil
	})
}

// create a Docker container for the built architectures
func createDockerContainers(builder *Builder, src *dagger.Directory) {
	if builder.ctx.Err() != nil {
		return
	}

	done := builder.LogStartSection("Building Docker containers for all architectures")
	defer done()

	alpine := builder.DaggerClient.Container().From("golang:latest")

	builder.WithArchitectures(func(goos, goarch string) error {
		done := builder.LogSubSection(fmt.Sprintf("building %s %s", goos, goarch))
		done(fmt.Sprintf("%s %s done", goos, goarch))

		dir := builder.DaggerClient.Host().Directory(path.Join("./build", goos, goarch))
		alpine.WithMountedDirectory("/tmp", dir).
			Exec(dagger.ContainerExecOpts{
				Args: []string{"cp", "/tmp/census", "/bin/census"},
			}).
			WithEntrypoint([]string{"/bin/census"}).ExitCode(builder.ctx)

		return nil
	})
}
