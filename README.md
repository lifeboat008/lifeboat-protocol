<p align="center"><img src="assets/logo.svg" alt="Lifeboat logo" width="112"></p>

<h1 align="center">lifeboat-protocol</h1>

<p align="center">The Go domain model for sponsor-funded open-source maintenance and rescue work.</p>

<p align="center">
  <a href="https://github.com/lifeboat008/lifeboat-protocol/actions/workflows/ci.yml"><img src="https://github.com/lifeboat008/lifeboat-protocol/actions/workflows/ci.yml/badge.svg" alt="Go CI"></a>
  <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT license">
  <img src="https://img.shields.io/badge/go-1.24%2B-00ADD8?logo=go&logoColor=white" alt="Go 1.24+">
  <img src="https://img.shields.io/badge/Stellar-testnet%20pilot-7D00FF" alt="Stellar testnet pilot">
  <img src="https://img.shields.io/badge/money-integer%20stroops-success" alt="Integer stroops">
  <img src="https://img.shields.io/badge/status-pilot-orange" alt="Pilot status">
  <a href="https://cjay-1.gitbook.io/lifeboat-docs/"><img src="https://img.shields.io/badge/docs-GitBook-3884FF" alt="Documentation"></a>
</p>

## Contents

- [What is Lifeboat](#what-is-lifeboat)
- [What this repository does](#what-this-repository-does)
- [Domain types](#domain-types)
- [Rules it enforces](#rules-it-enforces)
- [Quick start](#quick-start)
- [Use it in Go](#use-it-in-go)
- [How the four repositories fit together](#how-the-four-repositories-fit-together)
- [Project status](#project-status)
- [Open work](#open-work)
- [Documentation](#documentation)
- [Contributing and security](#contributing-and-security)
- [Maintainers](#maintainers)
- [Contributors](#contributors)
- [License](#license)

## What is Lifeboat

Lifeboat helps companies keep the open-source projects they depend on healthy. A sponsor commits a budget to a defined maintenance plan. A maintainer submits evidence of finished work as a claim. A named human reviewer approves or rejects it. Only an approved claim is paid, through a Stellar payment whose transaction hash anyone can verify.

GitHub activity is evidence, never approval. Lifeboat does not take control of a repository and never pays because a project looks inactive.

## What this repository does

`lifeboat-protocol` owns the shared vocabulary of the product: the types, their validation rules, and the budget arithmetic.

| Owns | Does not own |
| --- | --- |
| Project, plan, budget, evidence, claim, decision, and rescue task types | Stellar RPC or payments |
| Validation and integer-budget invariants | HTTP and persistence |
| `Budget.Reserve` and `Budget.Settle` arithmetic | GitHub webhook parsing |

It has no network code, no database, and no credentials.

## Domain types

| Type | Purpose |
| --- | --- |
| `Project` | An enrolled repository with its steward, GitHub App installation, and state (`active`, `needs_help`, `rescue_open`) |
| `Plan` | A sponsor's budget, named reviewer, eligible work, and end date |
| `Budget` | Total, reserved, and paid stroops, with `Remaining()` |
| `Evidence` | A GitHub event link that supports review but proves nothing alone |
| `Claim` | A maintainer's request for payment for one piece of evidence |
| `Decision` | The reviewer's approval or rejection, with a required reason |
| `RescueTask` | A steward's recorded request for new help on a project |

Claim states: `submitted`, `approved`, `rejected`, `payment_pending`, `paid`.

## Rules it enforces

- Money is `int64` **stroops** (10,000,000 stroops = 1 XLM). No floats.
- A plan's sponsor and reviewer must be different.
- The testnet plan asset is `XLM`, with a positive budget and at least one eligible work type.
- A claim must match its plan, fall inside the plan period, and use an eligible work type.
- Only the reviewer named on the plan can decide a claim, once, with a non-blank reason.
- A rescue task must come from the project's own steward.
- `Reserve` refuses zero, negative, over-budget, and overflowing amounts. `Settle` refuses amounts above what is reserved.

`lifeboat-api` applies the same budget changes atomically in SQLite, so two approvals cannot both spend the last of a budget.

## Quick start

With Go 1.24 or newer:

```bash
git clone https://github.com/lifeboat008/lifeboat-protocol.git
cd lifeboat-protocol
go test ./...
go vet ./...
```

No network access or credentials are needed.

## Use it in Go

```go
import protocol "github.com/lifeboat008/lifeboat-protocol"

budget := protocol.Budget{TotalStroops: 5_000_000_000} // 500 XLM
budget, err := budget.Reserve(1_200_000_000)           // approve a 120 XLM claim
if err != nil {
    // over budget or invalid amount
}
remaining := budget.Remaining() // 3_800_000_000
```

## How the four repositories fit together

```text
lifeboat-protocol ──► lifeboat-ledger ──► lifeboat-api ◄── lifeboat-github
```

| Repository | Responsibility |
| --- | --- |
| [lifeboat-protocol](https://github.com/lifeboat008/lifeboat-protocol) | Domain types, validation, budget arithmetic (this repository) |
| [lifeboat-ledger](https://github.com/lifeboat008/lifeboat-ledger) | Stellar testnet payment construction, submission, reconciliation |
| [lifeboat-api](https://github.com/lifeboat008/lifeboat-api) | HTTP API, persistence, authorization, audit history |
| [lifeboat-github](https://github.com/lifeboat008/lifeboat-github) | GitHub webhook verification and evidence submission |

## Project status

Lifeboat is a **Stellar testnet pilot**. It does not move real money, has had no formal security audit, and is a single-operator design. Custody, payout asset, jurisdiction, and legal review must be settled before mainnet. See [Project status](https://cjay-1.gitbook.io/lifeboat-docs/project-status) and the [Mainnet gate](https://cjay-1.gitbook.io/lifeboat-docs/security-and-governance/mainnet-gate).

## Open work

- [Model disputed claims and reviewer conflict rules](https://github.com/lifeboat008/lifeboat-protocol/issues/1) (complexity: high)

The product requirements, architecture, and Wave plan are versioned in `lifeboat-api` ([PRD](https://github.com/lifeboat008/lifeboat-api/blob/main/product/docs/PRD.md), [architecture](https://github.com/lifeboat008/lifeboat-api/blob/main/product/docs/ARCHITECTURE.md), [Wave plan](https://github.com/lifeboat008/lifeboat-api/blob/main/product/docs/WAVE.md)).

## Documentation

Full documentation, including the [protocol reference](https://cjay-1.gitbook.io/lifeboat-docs/protocol-reference), API reference, role guides, security notes, and operations runbooks, is at **[cjay-1.gitbook.io/lifeboat-docs](https://cjay-1.gitbook.io/lifeboat-docs/)**.

## Contributing and security

- Read [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request. Run `gofmt`, `go vet ./...`, and `go test ./...`, and add tests for changed behavior.
- Use synthetic data only. Never commit keys, tokens, or real sponsor or maintainer details.
- Preserve human approval before any payout.
- Report vulnerabilities privately as described in [SECURITY.md](SECURITY.md). Do not post exploit details in a public issue.

## Maintainers

| Maintainer | Contact |
| --- | --- |
| [lifeboat008](https://github.com/lifeboat008) | [Open an issue](https://github.com/lifeboat008/lifeboat-protocol/issues) for public work; use SECURITY.md for private reports |

## Contributors

<a href="https://github.com/lifeboat008/lifeboat-protocol/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=lifeboat008/lifeboat-protocol" alt="Contributors">
</a>

## License

[MIT](LICENSE). This pilot has not had a formal security audit.
