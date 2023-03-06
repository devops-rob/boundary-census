package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"dagger.io/dagger"
)

var hasTTY = flag.Bool("tty", false, "does the output terminal have tty")

// Dagger build pipeline
func main() {
	flag.Parse()

	// remove the build folder if it exists
	os.RemoveAll("./build")

	// define build matrix
	oses := []string{"linux", "darwin"}
	arches := []string{"amd64", "arm64"}

	builder, err := NewBuilder(oses, arches, *hasTTY)
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

// run the application unit tests
func test(builder *Builder, src *dagger.Directory) {
	if builder.ctx.Err() != nil {
		return
	}

	testContainer := builder.DaggerClient.Container().
		From("golang:latest").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src")

	testContainer = builder.WithModCache(testContainer).
		WithExec([]string{"go", "test", "-v", "./..."})

	done := builder.LogStartSectionWithContainer("Running unit tests", testContainer)
	defer done()

	_, err := testContainer.
		ExitCode(builder.ctx)

	if err != nil {
		builder.LogError("unable to test application", err)
	}
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

	// add the mod cache
	golang = builder.WithModCache(golang)

	// build for all architectures
	builder.WithArchitectures(func(goos, goarch string) error {
		p := path.Join("build/", goos, goarch)
		build := golang.WithEnvVariable("GOOS", goos).
			WithEnvVariable("GOARCH", goarch).
			WithExec([]string{"go", "build", "-o", path.Join(p, "census")})

		done := builder.LogSubSectionWithContainer(fmt.Sprintf("building %s %s", goos, goarch), build)

		_, err := build.Directory(p).Export(builder.ctx, p)
		if err != nil {
			builder.LogError("failed to build", err)
		}

		done(fmt.Sprintf("%s %s complete: ./build/%s/%s", goos, goarch, goos, goarch))

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
