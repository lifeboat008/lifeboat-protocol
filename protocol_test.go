package protocol

import (
	"testing"
	"time"
)

func TestBudgetReserveAndSettle(t *testing.T) {
	budget := Budget{TotalStroops: 100}
	reserved, err := budget.Reserve(60)
	if err != nil || reserved.Remaining() != 40 {
		t.Fatalf("reserve: %v %#v", err, reserved)
	}
	if _, err := reserved.Reserve(41); err == nil {
		t.Fatal("allowed overspending")
	}
	settled, err := reserved.Settle(60)
	if err != nil || settled.PaidStroops != 60 || settled.ReservedStroops != 0 {
		t.Fatalf("settle: %v %#v", err, settled)
	}
}

func TestDecisionNeedsCorrectReviewer(t *testing.T) {
	claim := Claim{ID: "claim", State: ClaimSubmitted}
	plan := Plan{ReviewerID: "reviewer"}
	decision := Decision{ClaimID: "claim", ReviewerID: "intruder", Approved: true, Reason: "reviewed", DecidedAt: time.Now()}
	if decision.Validate(claim, plan) == nil {
		t.Fatal("unauthorized reviewer was accepted")
	}
}

func TestRescueRequiresSteward(t *testing.T) {
	project := Project{ID: "project", StewardID: "steward"}
	task := RescueTask{ID: "rescue", ProjectID: "project", StewardID: "stranger", Reason: "needs help", Scope: "release", ReviewerID: "reviewer", OpenedAt: time.Now()}
	if task.Validate(project) == nil {
		t.Fatal("unauthorized rescue task was accepted")
	}
}
