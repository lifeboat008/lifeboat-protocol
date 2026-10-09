<p align="center"><img src="assets/logo.svg" alt="Lifeboat logo" width="112"></p>

# lifeboat-protocol

[![Go CI](https://github.com/lifeboat008/lifeboat-protocol/actions/workflows/ci.yml/badge.svg)](https://github.com/lifeboat008/lifeboat-protocol/actions/workflows/ci.yml)

The Go domain model for sponsor-funded open-source maintenance and rescue work.

## Owns

- Project, sponsor plan, budget, evidence, claim, approval, and rescue state types.
- Validation and integer-budget invariants.

## Does not own

Stellar RPC, HTTP, persistence, or GitHub webhook parsing.

`Project`, `Plan`, `Budget`, `Evidence`, `Claim`, `Decision`, and `RescueTask` are the shared domain types. Validation keeps reviewer decisions tied to the named reviewer, rescue tasks tied to the named steward, claims tied to eligible work, and budgets in integer stroops. `Budget.Reserve` and `Budget.Settle` implement the arithmetic; `lifeboat-api` performs their changes atomically in SQLite. `lifeboat-ledger` consumes approved claims but never decides who gets paid.

## Quick start

With Go 1.24 or newer, run `go test ./...` and `go vet ./...`. This module has no network service or credentials.

The [product requirements](https://github.com/lifeboat008/lifeboat-api/blob/main/product/docs/PRD.md), [architecture](https://github.com/lifeboat008/lifeboat-api/blob/main/product/docs/ARCHITECTURE.md), and [Wave plan](https://github.com/lifeboat008/lifeboat-api/blob/main/product/docs/WAVE.md) are versioned in `lifeboat-api`.

## Documentation

Full documentation, including the API reference, role guides, security notes, and operations runbooks, is at [cjay-1.gitbook.io/lifeboat-docs](https://cjay-1.gitbook.io/lifeboat-docs/).

## Contributing and security

See [CONTRIBUTING.md](CONTRIBUTING.md) for pull requests, [SECURITY.md](SECURITY.md) for private vulnerability reports, and [LICENSE](LICENSE) for MIT terms.

Maintainers: [lifeboat008](https://github.com/lifeboat008). Discuss public work in [issues](https://github.com/lifeboat008/lifeboat-protocol/issues); report vulnerabilities privately as described in SECURITY.md. This pilot has not had a formal security audit.
