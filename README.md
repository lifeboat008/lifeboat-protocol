<p align="center"><img src="assets/logo.svg" alt="Lifeboat logo" width="112"></p>

# lifeboat-protocol

The Go domain model for sponsor-funded open-source maintenance and rescue work.

## Owns

- Project, sponsor plan, budget, evidence, claim, approval, and rescue state types.
- Validation and integer-budget invariants.

## Does not own

Stellar RPC, HTTP, persistence, or GitHub webhook parsing.

`Project`, `Plan`, `Budget`, `Evidence`, `Claim`, `Decision`, and `RescueTask` are the shared domain types. Validation keeps reviewer decisions tied to the named reviewer, rescue tasks tied to the named steward, claims tied to eligible work, and budgets in integer stroops. `Budget.Reserve` and `Budget.Settle` implement the arithmetic; `lifeboat-api` performs their changes atomically in SQLite. `lifeboat-ledger` consumes approved claims but never decides who gets paid.

Product PRD and architecture live in the parent `lifeboat/docs` folder in the local workspace.
