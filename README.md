# run

[![CI](https://github.com/dio/run/actions/workflows/ci.yml/badge.svg)](https://github.com/dio/run/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/dio/run.svg)](https://pkg.go.dev/github.com/dio/run)
[![Go Report Card](https://goreportcard.com/badge/github.com/dio/run)](https://goreportcard.com/report/github.com/dio/run)

A drop-in replacement for [oklog/run](https://github.com/oklog/run) with automatic panic recovery.

## Features

- **Drop-in replacement**: Compatible with `oklog/run.Group` API
- **Automatic panic recovery**: All panics in actors are caught and converted to errors
- **Custom panic handlers**: Optional custom handling of panic situations
- **Zero dependencies**: Only depends on `oklog/run`
- **100% test coverage**: Thoroughly tested with comprehensive test suite

## Installation

```bash
go get github.com/dio/run
```

## Quick Start

```go
package main

import (
    "errors"
    "fmt"
    "log"

    "github.com/dio/run"
)

func main() {
    g := run.New()

    // Add an actor that might panic
    g.Add(func() error {
        // This panic will be caught and converted to an error
        panic("something went wrong")
    }, func(err error) {
        log.Printf("Actor interrupted: %v", err)
    })

    // Add a normal actor
    g.Add(func() error {
        return errors.New("normal error")
    }, func(err error) {
        log.Printf("Actor interrupted: %v", err)
    })

    if err := g.Run(); err != nil {
        log.Printf("Group failed: %v", err)
    }
}
```

## Usage

### Basic Usage

The `run.Group` works exactly like `oklog/run.Group` but with automatic panic recovery:

```go
g := run.New()

g.Add(execute, interrupt)

err := g.Run()
```

### Custom Panic Handler

You can provide a custom panic handler to process panics in your own way:

```go
g := run.NewWithHandler(func(panicValue any, stackTrace string) error {
    // Log the panic with custom formatting
    log.Printf("PANIC: %v\nStack trace:\n%s", panicValue, stackTrace)

    // Return a custom error
    return fmt.Errorf("actor panicked: %v", panicValue)
})
```

### Error Handling

When an actor panics, it's automatically converted to a `PanicError`:

```go
err := g.Run()
if err != nil {
    if panicErr, ok := err.(*run.PanicError); ok {
        fmt.Printf("Actor panicked with value: %v\n", panicErr.Value)
        fmt.Printf("Stack trace:\n%s\n", panicErr.StackTrace)
    }
}
```

## API Reference

### Types

#### Group

```go
type Group struct {
    // PanicHandler is called when an actor panics.
    // If nil, panics are converted to PanicError and returned.
    PanicHandler func(panicValue any, stackTrace string) error
}
```

#### PanicError

```go
type PanicError struct {
    Value      any    // The panic value
    StackTrace string // Stack trace at panic time
}
```

### Functions

#### New

```go
func New() *Group
```

Creates a new Group with panic recovery enabled.

#### NewWithHandler

```go
func NewWithHandler(handler func(panicValue any, stackTrace string) error) *Group
```

Creates a new Group with a custom panic handler.

### Methods

#### Add

```go
func (g *Group) Add(execute func() error, interrupt func(error))
```

Add an actor to the group. This is a drop-in replacement for `oklog/run.Group.Add()`.

#### Run

```go
func (g *Group) Run() error
```

Execute all actors. This is identical to `oklog/run.Group.Run()`.

## Development

### Prerequisites

- Go 1.24+
- golangci-lint (installed via `go tool`)

### Building

```bash
go build ./...
```

### Testing

```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Linting

```bash
go tool golangci-lint run
```

### Formatting

```bash
go tool golangci-lint fmt
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Ensure all tests pass (`go test ./...`)
6. Run linting (`go tool golangci-lint run`)
7. Commit your changes (`git commit -am 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [oklog/run](https://github.com/oklog/run) - The original actor model implementation
- Inspired by the need for safer concurrent programming in Go
