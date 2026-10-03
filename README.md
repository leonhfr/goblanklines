# goblanklines

A Go linter that inserts blank lines according to a set of rules, to make code easier to scan.

It runs as a [golangci-lint](https://golangci-lint.run/) module plugin, built on [`golang.org/x/tools/go/analysis`](https://pkg.go.dev/golang.org/x/tools/go/analysis). Every rule provides an automatic fix.

## Rules

Each rule has a settings key. All rules run by default; the `rules` setting restricts the analyzer to a subset.

| key          | rule                                                                                                                                                                                                                   |
| ------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `block`      | Blank line after a block statement (`if`, `for`, `switch`, `select`, `defer`), unless the block ends its enclosing block, is followed by `else`/`case`/`default`, or is a defer immediately after an error-check `if`. |
| `testhelper` | Blank line after a run of `t.Parallel()`/`t.Helper()` calls, unless the run ends its block.                                                                                                                            |
| `testify`    | Blank line after a run of testify `assert`/`require` calls (package-level or `*Assertions` method calls), unless the run ends its block.                                                                              |
| `testrun`    | Blank lines before and after a run of `t.Run(...)` calls, unless the run starts or ends its block.                                                                                                                     |
| `funcdecl`   | Blank line between two consecutive function/method declarations, unless both are one-liners.                                                                                                                           |

The `testhelper`, `testify`, and `testrun` rules type-check their targets: a call only matches if its receiver implements `testing.TB`, or if the called function resolves to the `testify` `assert`/`require` package.

## Installation

Build a custom golangci-lint binary that includes this plugin. Add it to `.custom-gcl.yml`:

```yaml
version: v2.14.0
name: goblanklines-golangci-lint
destination: bin
plugins:
  - module: github.com/leonhfr/goblanklines
    version: latest
```

Then build it:

```bash
golangci-lint custom
```

See the [module plugin system docs](https://golangci-lint.run/docs/plugins/module-plugins/) for the general mechanism.

## Configuration

Enable the linter in `.golangci.yml` and configure it under `linters.settings.custom`:

```yaml
linters:
  enable:
    - goblanklines
  settings:
    custom:
      goblanklines:
        type: module
        description: Inserts blank lines according to a set of rules.
        original-url: github.com/leonhfr/goblanklines
        settings:
          rules: [] # empty = all rules; or e.g. ["block", "forloop"]
```

Run it with fixes applied:

```bash
./bin/goblanklines-golangci-lint run --fix ./...
```

## Go API

The analyzer is also usable directly, outside golangci-lint, with any `go/analysis` driver such as `singlechecker` or `analysistest`:

```go
import "github.com/leonhfr/goblanklines"

analyzer := goblanklines.New()
```

## Development

This repo uses [mise](https://mise.jdx.dev/) to pin tool versions and [hk](https://hk.jdx.dev/) to run checks.

```bash
mise run check # verify formatting, lint, and tests
mise run fix   # apply formatting and lint fixes
mise run test  # run the Go test suite
```

## License

MIT, see [LICENSE](LICENSE).
