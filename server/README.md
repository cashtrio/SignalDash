# Server

## Package manger

install `asdf` package manager.

```go
go install github.com/asdf-vm/asdf/cmd/asdf@v0.18.0
```

install golang plugin for asdf.

```sh
asdf plugin add golang
```

install target go version and anchor it.

```sh
asdf install golang 1.24.0
asdf set golang 1.24.0
```

check the current version.

```sh
 asdf current golang
Name            Version         Source                                              Installed
golang          1.24.0          ~/SignalDash/server/.tool-versions true
```

## Mocking

download binary from mockery release.

- https://vektra.github.io/mockery/latest/installation/#installation

code interface and set `go generate` directive.

```go
//go:generate mockery --name=MyDummyInterface --dir=./ --output=./mocks
type MyDummyInterface interface {
	Log() (string, error)
}
```

run the command to generate mock codes.

```sh
go generate ./...
```
