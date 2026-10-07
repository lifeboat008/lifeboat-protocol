# lifeboat-protocol

The Go domain model for sponsor-funded open-source maintenance and rescue work.

## Owns

- Project, sponsor plan, budget, evidence, claim, approval, and rescue state types.
- Allowed transitions and invariants such as no overspending or double payment.
- Versioned domain events.

## Does not own

Stellar RPC, HTTP, persistence, or GitHub webhook parsing.

## First implementation slice

Model a plan and claim; require an authorized approval before the claim is payable; test duplicate and over-budget transitions. The next consumer is `lifeboat-ledger`, with `lifeboat-api` also importing the domain model.

Product PRD and architecture live in the parent `lifeboat/docs` folder in the local workspace.
