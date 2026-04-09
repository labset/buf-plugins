# buf-plugins

Custom [Buf](https://buf.build) lint plugins and a protobuf schema workspace to exercise them.

## Repository structure

```
.
├── plugin/         # Custom buf lint plugins (Go)
│   ├── cmd/
│   │   └── api-lint-plugin/        # Plugin for API lint rules
│   ├── internal/
│   │   └── api/                    # Lint rule implementations
│   └── dist/                       # Compiled plugin binaries
└── schema/         # Protobuf schema workspace
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

## Requirements

This project uses [mise](https://mise.jdx.dev) to manage tools. Install it then run:

```bash
mise install
```

This installs Go and Buf, and adds the compiled plugin binaries in `plugin/dist/` to your `PATH`.

## Plugins

Custom buf lint plugins are written in Go using the [bufplugin SDK](https://github.com/bufbuild/bufplugin-go).

### Install

Install the plugin directly with `go install`:

```bash
go install github.com/viqueen/buf-plugins/plugin/cmd/api-lint-plugin@latest
```

This places the binary in `$(go env GOBIN)` (or `$(go env GOPATH)/bin`). Make sure that directory is on your `PATH` so `buf` can discover the plugin.

### Build from source

```bash
cd plugin
mise run build
```

Binaries are output to `plugin/dist/` and automatically available on `PATH` via mise.

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

## Schema

The `schema` workspace defines protobuf APIs for two domains:

- **user.v1** — `UserService` with Create, Get, Update, Delete, List operations
- **team.v1** — `TeamService` with Create, Get, Update, Delete, List operations, plus `AddTeamMember`, `AddTeamMembers`, and `RemoveTeamMember`

Cross-domain references use ref types (e.g. `user.v1.UserRef`) rather than plain string IDs.

### Lint

```bash
cd schema
mise run lint
```
