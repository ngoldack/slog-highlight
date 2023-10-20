package main

import (
	"context"
	"os"

	"dagger.io/dagger"
)

type ciType string

const CI = ciType("ci")

func main() {
	ctx := context.WithValue(context.Background(), CI, os.Getenv("CI") == "true")

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	cache := client.CacheVolume("go-mod-cache")

	if err := test(ctx, client, cache); err != nil {
		panic(err)
	}

	if err := lint(ctx, client); err != nil {
		panic(err)
	}
}

func test(ctx context.Context, client *dagger.Client, cache *dagger.CacheVolume) error {
	client = client.Pipeline("test")

	srcDir := client.Host().Directory(".")

	c := client.Container().From("golang:latest")
	c = c.WithDirectory("/src", srcDir).WithWorkdir("/src")
	c = c.WithMountedCache("/go/pkg/mod", cache)

	if ci, ok := ctx.Value(CI).(bool); ok && ci {
		c = c.WithExec([]string{"go", "test", "-race", "-v", "-coverprofile=coverage.txt", "-covermode=atomic", "./..."})
	} else {
		c = c.WithExec([]string{"go", "test", "-race", "-v", "./..."})
	}

	// copy coverage.txt back
	output := c.Directory("/src").File("coverage.txt")
	_, err := output.Export(ctx, "coverage.txt")
	if err != nil {
		return err
	}

	if _, err := c.Stderr(ctx); err != nil {
		return err
	}

	return nil
}

func lint(ctx context.Context, client *dagger.Client) error {
	client = client.Pipeline("lint")

	srcDir := client.Directory().Directory("src")
	c := client.Container().From("golangci/golangci-lint:latest")
	c = c.WithDirectory("/src", srcDir).WithWorkdir("/src")

	if ci, ok := ctx.Value(CI).(bool); ok && ci {
		c.WithExec([]string{"golangci-lint", "run", "--out-format=github-actions"})
	} else {
		c.WithExec([]string{"golangci-lint", "run"})
	}

	return nil
}
