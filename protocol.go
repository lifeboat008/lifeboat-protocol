package protocol

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type ProjectState string

const (
	ProjectActive     ProjectState = "active"
	ProjectNeedsHelp  ProjectState = "needs_help"
	ProjectRescueOpen ProjectState = "rescue_open"
)

type ClaimState string

const (
	ClaimSubmitted      ClaimState = "submitted"
	ClaimApproved       ClaimState = "approved"
	ClaimRejected       ClaimState = "rejected"
	ClaimPaymentPending ClaimState = "payment_pending"
	ClaimPaid           ClaimState = "paid"
)

type Project struct {
	ID             string       `json:"id"`
	Repository     string       `json:"repository"`
	InstallationID int64        `json:"installation_id"`
	StewardID      string       `json:"steward_id"`
	State          ProjectState `json:"state"`
}

func (p Project) Validate() error {
	if p.ID == "" || p.Repository == "" || p.InstallationID <= 0 || p.StewardID == "" {
		return errors.New("project id, repository, installation, and steward are required")
	}
	if !strings.Contains(p.Repository, "/") {
		return errors.New("repository must use owner/name format")
	}
	if p.State != ProjectActive && p.State != ProjectNeedsHelp && p.State != ProjectRescueOpen {
		return errors.New("invalid project state")
	}
	return nil
}

type Plan struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"project_id"`
	SponsorID     string    `json:"sponsor_id"`
	ReviewerID    string    `json:"reviewer_id"`
	Asset         string    `json:"asset"`
	BudgetStroops int64     `json:"budget_stroops"`
	EligibleWork  []string  `json:"eligible_work"`
	PeriodEndsAt  time.Time `json:"period_ends_at"`
}

func (p Plan) Validate(now time.Time) error {
	if p.ID == "" || p.ProjectID == "" || p.SponsorID == "" || p.ReviewerID == "" {
		return errors.New("plan id, project, sponsor, and reviewer are required")
	}
	if p.SponsorID == p.ReviewerID {
		return errors.New("sponsor and reviewer must be separate")
	}
	if p.Asset != "XLM" || p.BudgetStroops <= 0 || len(p.EligibleWork) == 0 {
		return errors.New("testnet plan requires a positive XLM budget and eligible work")
	}
	if !p.PeriodEndsAt.After(now) {
		return errors.New("plan period must end in the future")
	}
	return nil
}

type Budget struct {
	TotalStroops    int64 `json:"total_stroops"`
	ReservedStroops int64 `json:"reserved_stroops"`
	PaidStroops     int64 `json:"paid_stroops"`
}

func (b Budget) Remaining() int64 {
	return b.TotalStroops - b.ReservedStroops - b.PaidStroops
}

func (b Budget) Validate() error {
	if b.TotalStroops < 0 || b.ReservedStroops < 0 || b.PaidStroops < 0 || b.ReservedStroops > b.TotalStroops || b.PaidStroops > b.TotalStroops-b.ReservedStroops {
		return errors.New("invalid budget balance")
	}
	return nil
}

func (b Budget) Reserve(amount int64) (Budget, error) {
	if err := b.Validate(); err != nil {
		return Budget{}, err
	}
	if amount <= 0 || amount > b.Remaining() || b.ReservedStroops > math.MaxInt64-amount {
		return Budget{}, errors.New("claim exceeds available budget")
	}
	b.ReservedStroops += amount
	return b, nil
}

func (b Budget) Settle(amount int64) (Budget, error) {
	if err := b.Validate(); err != nil {
		return Budget{}, err
	}
	if amount <= 0 || amount > b.ReservedStroops || b.PaidStroops > math.MaxInt64-amount {
		return Budget{}, errors.New("invalid settlement amount")
	}
	b.ReservedStroops -= amount
	b.PaidStroops += amount
	return b, nil
}

type Evidence struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"project_id"`
	InstallationID int64     `json:"installation_id"`
	DeliveryID     string    `json:"delivery_id"`
	Kind           string    `json:"kind"`
	URL            string    `json:"url"`
	ObservedAt     time.Time `json:"observed_at"`
}

func (e Evidence) Validate() error {
	if e.ID == "" || e.ProjectID == "" || e.InstallationID <= 0 || e.DeliveryID == "" || e.Kind == "" || e.URL == "" || e.ObservedAt.IsZero() {
		return errors.New("incomplete evidence")
	}
	if !strings.HasPrefix(e.URL, "https://github.com/") {
		return errors.New("evidence must link to GitHub")
	}
	return nil
}

type Claim struct {
	ID                 string     `json:"id"`
	ProjectID          string     `json:"project_id"`
	PlanID             string     `json:"plan_id"`
	MaintainerID       string     `json:"maintainer_id"`
	EvidenceID         string     `json:"evidence_id"`
	WorkType           string     `json:"work_type"`
	Summary            string     `json:"summary"`
	AmountStroops      int64      `json:"amount_stroops"`
	DestinationAccount string     `json:"destination_account"`
	State              ClaimState `json:"state"`
}

func (c Claim) Validate(plan Plan, now time.Time) error {
	if c.ID == "" || c.ProjectID != plan.ProjectID || c.PlanID != plan.ID || c.MaintainerID == "" || c.EvidenceID == "" || c.Summary == "" || c.DestinationAccount == "" {
		return errors.New("incomplete claim or plan mismatch")
	}
	if c.AmountStroops <= 0 || now.After(plan.PeriodEndsAt) {
		return errors.New("claim amount or period is invalid")
	}
	for _, work := range plan.EligibleWork {
		if c.WorkType == work {
			return nil
		}
	}
	return fmt.Errorf("work type %q is outside the plan", c.WorkType)
}

type Decision struct {
	ClaimID    string    `json:"claim_id"`
	ReviewerID string    `json:"reviewer_id"`
	Approved   bool      `json:"approved"`
	Reason     string    `json:"reason"`
	DecidedAt  time.Time `json:"decided_at"`
}

func (d Decision) Validate(claim Claim, plan Plan) error {
	if d.ClaimID != claim.ID || d.ReviewerID != plan.ReviewerID || claim.State != ClaimSubmitted || strings.TrimSpace(d.Reason) == "" || d.DecidedAt.IsZero() {
		return errors.New("unauthorized or incomplete claim decision")
	}
	return nil
}

type RescueTask struct {
	ID         string    `json:"id"`
	ProjectID  string    `json:"project_id"`
	StewardID  string    `json:"steward_id"`
	Reason     string    `json:"reason"`
	Scope      string    `json:"scope"`
	ReviewerID string    `json:"reviewer_id"`
	OpenedAt   time.Time `json:"opened_at"`
}

func (r RescueTask) Validate(project Project) error {
	if r.ID == "" || r.ProjectID != project.ID || r.StewardID != project.StewardID || r.Reason == "" || r.Scope == "" || r.ReviewerID == "" || r.OpenedAt.IsZero() {
		return errors.New("unauthorized or incomplete rescue task")
	}
	return nil
}
