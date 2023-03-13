package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"

	"dagger.io/dagger"
)

var hasTTY = flag.Bool("tty", false, "does the output terminal have tty")
var publishImage = flag.Bool("publish-image", false, "should the images be published")

// var dockerRegistry = flag.String("docker-registry", "devops-rob/census", "registry for docker images")
var dockerRegistry = flag.String("docker-registry", "nicholasjackson/census", "registry for docker images")

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

	src = src.WithoutDirectory("shipyard")

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

		done(fmt.Sprintf("build complete output: ./build/%s/%s", goos, goarch))

		return nil
	})
}

// create a Docker container for the built architectures
func createDockerContainers(builder *Builder, src *dagger.Directory) {
	if builder.ctx.Err() != nil {
		return
	}

	containers := []*dagger.Container{}

	done := builder.LogStartSection("Building Docker containers for all architectures")
	defer done()

	builder.WithArchitectures(func(goos, goarch string) error {
		if goos != "linux" {
			return nil
		}

		done := builder.LogSubSection(fmt.Sprintf("building %s %s", goos, goarch))

		platform := dagger.Platform(fmt.Sprintf("%s/%s", goos, goarch))

		dir := builder.DaggerClient.Host().Directory(path.Join("./build", goos, goarch))

		build := builder.DaggerClient.
			Container(dagger.ContainerOpts{Platform: platform}).
			From("alpine:latest").
			WithMountedDirectory("/tmp", dir).
			WithExec([]string{"cp", "/tmp/census", "/bin/census"}, dagger.ContainerWithExecOpts{}).
			WithEntrypoint([]string{"/bin/census"})

		containers = append(containers, build)

		// export a local docker container for the current architecture
		if goarch == runtime.GOARCH {
			err := exportLocalImage(builder, build, platform)
			if err != nil {
				builder.LogError("error exporting image", err)
			}
		}

		sha, _ := builder.GitSHA()
		if *publishImage {
			_, err := builder.DaggerClient.Container().
				Publish(
					builder.ctx,
					fmt.Sprintf("%s:%s", *dockerRegistry, sha),
					dagger.ContainerPublishOpts{PlatformVariants: containers},
				)

			if err != nil {
				builder.LogError("error publishing image", err)
			}
		}

		done(fmt.Sprintf("published container %s:%s for architecture %s/%s", *dockerRegistry, sha, goos, goarch))

		return nil
	})
}

func exportLocalImage(builder *Builder, container *dagger.Container, platform dagger.Platform) error {
	done := builder.LogSubSection(fmt.Sprintf("exporting %s", platform))

	outputDirectory := path.Join(os.TempDir(), "build_output")
	os.MkdirAll(outputDirectory, os.ModePerm)

	imagePath := path.Join(outputDirectory, "localexport.tar")
	//defer os.RemoveAll(outputDirectory)

	// export the image to the local path
	_, err := builder.DaggerClient.
		Container(dagger.ContainerOpts{Platform: platform}).
		Export(
			builder.ctx,
			imagePath,
			dagger.ContainerExportOpts{PlatformVariants: []*dagger.Container{container}},
		)

	if err != nil {
		return err
	}

	// run docker import to import to the local system
	cmd := exec.Command(
		"docker",
		"load",
		"-i",
		imagePath,
	)

	out := &bytes.Buffer{}
	cmd.Stdout = out
	cmd.Stderr = out

	err = cmd.Run()
	if err != nil {
		return err
	}

	// get the sha from the output
	r, err := regexp.Compile(`sha256:(.*)`)
	if err != nil {
		return err
	}

	res := r.FindStringSubmatch(out.String())
	if len(res) != 2 {
		return fmt.Errorf("expected sha from docker load output, got: %s", out.String())
	}

	sha := res[1]
	cmd = exec.Command(
		"docker",
		"tag",
		sha,
		fmt.Sprintf("%s:local", *dockerRegistry),
	)

	out = &bytes.Buffer{}
	cmd.Stdout = out
	cmd.Stderr = out

	err = cmd.Run()
	if err != nil {
		return err
	}

	done(fmt.Sprintf("export container complete for architecture %s", platform))
	return nil
}

// creates a github release for the built artifacts
func createGitHubRelease(builder *Builder) {
}
