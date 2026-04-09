# buf-plugins

Custom [Buf](https://buf.build) lint plugins, written in Go with the [bufplugin SDK](https://github.com/bufbuild/bufplugin-go).

## Usage

### Install

Install the plugin with `go install`:

```bash
go install github.com/viqueen/buf-plugins/plugin/cmd/api-lint-plugin@latest
```

This drops the binary in `$(go env GOBIN)` (or `$(go env GOPATH)/bin`). Make sure that directory is on your `PATH` so `buf` can discover the plugin.

### Configure `buf.yaml`

Reference the plugin and opt into the rules you want from your schema's `buf.yaml`:

```yaml
version: v2
modules:
  - path: protos

plugins:
  - plugin: api-lint-plugin

lint:
  use:
    - FILE_NAME_CONVENTION
    - UPDATE_REQUEST_FIELD_MASK
    - REPEATED_FIELD_VALIDATION
```

Then run `buf lint` as usual.

### Available plugins

| Plugin | Description |
|---|---|
| `api-lint-plugin` | Enforces API lint rules |

### Lint rules

#### `api-lint-plugin`

| Rule ID | Default | Description |
|---|---|---|
| `FILE_NAME_CONVENTION` | yes | Proto files must be named `enums.proto`, `models.proto`, `refs.proto`, or `service_<name>.proto` |
| `UPDATE_REQUEST_FIELD_MASK` | yes | `UpdateXxxRequest` messages must have a `google.protobuf.FieldMask update_mask` field to support partial updates |
| `REPEATED_FIELD_VALIDATION` | yes | Repeated fields in request messages (including nested messages) must have a `max_items` constraint to prevent unbounded input attacks |

## Repository structure

```
.
├── plugin/         # Custom buf lint plugins (Go)
│   ├── cmd/
│   │   └── api-lint-plugin/        # Plugin for API lint rules
│   ├── internal/
│   │   └── api/                    # Lint rule implementations
│   └── dist/                       # Compiled plugin binaries
└── schema/         # Protobuf schema workspace used to exercise the plugins
    └── protos/
        ├── user/v1/                # User API
        │   ├── enums.proto
        │   ├── models.proto
        │   ├── refs.proto
        │   └── service_user.proto
        └── team/v1/                # Team API
            ├── enums.proto
            ├── models.proto
            ├── refs.proto
            └── service_team.proto
```

The `schema/` workspace defines protobuf APIs for two domains:

- **user.v1** — `UserService` with Create, Get, Update, Delete, List operations
- **team.v1** — `TeamService` with Create, Get, Update, Delete, List operations, plus `AddTeamMember`, `AddTeamMembers`, and `RemoveTeamMember`

Cross-domain references use ref types (e.g. `user.v1.UserRef`) rather than plain string IDs.

## Development

This project uses [mise](https://mise.jdx.dev) to manage tools. Install it, then from the repo root run:

```bash
mise install
```

This installs Go and Buf, and puts the compiled plugin binaries from `plugin/dist/` on your `PATH`.

### Build the plugin

```bash
cd plugin
mise run build
```

Binaries are output to `plugin/dist/`.

### Lint the sample schema

```bash
cd schema
mise run lint
```

This runs `buf lint` against the `schema/` workspace using the locally built `api-lint-plugin`, which is the quickest way to iterate on rule changes.
