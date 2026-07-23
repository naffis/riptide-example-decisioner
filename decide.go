package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// DefaultBump is the fixed CPM margin this sample adds above the resolved base floor.
const DefaultBump = "0.2500"

// DecideRequest mirrors the protojson shape of riptide.decision.v1.DecideRequest
// (field names match the JSON the Riptide host posts to EXTERNAL_ENDPOINT models).
type DecideRequest struct {
	Request *DecisionRequest `json:"request"`
}

// DecisionRequest is the inner decision payload.
type DecisionRequest struct {
	Point      string        `json:"point"`
	TenantID   string        `json:"tenantId"`
	RequestID  string        `json:"requestId"`
	Features   *Features     `json:"features"`
	Candidates []Candidate   `json:"candidates"`
	Policy     *TenantPolicy `json:"policy"`
	Budget     *Budget       `json:"budget"`
}

type Features struct {
	Categorical map[string]string  `json:"categorical"`
	Numeric     map[string]float64 `json:"numeric"`
	Flags       map[string]bool    `json:"flags"`
}

type Candidate struct {
	ID    string            `json:"id"`
	Kind  string            `json:"kind"`
	Attrs map[string]string `json:"attrs"`
	Price string            `json:"price"`
}

type TenantPolicy struct {
	Floor      string            `json:"floor"`
	MaxBid     string            `json:"maxBid"`
	AllowedIDs []string          `json:"allowedIds"`
	Caps       map[string]string `json:"caps"`
}

type Budget struct {
	LatencyBudgetMs int32 `json:"latencyBudgetMs"`
}

// DecideResponse mirrors protojson DecideResponse.
type DecideResponse struct {
	Response *DecisionResponse `json:"response"`
}

type DecisionResponse struct {
	Price      *Price    `json:"price"`
	Confidence float64   `json:"confidence"`
	Rationale  string    `json:"rationale"`
	Model      *ModelRef `json:"model"`
}

type Price struct {
	CPM    string `json:"cpm"`
	Accept bool   `json:"accept"`
}

type ModelRef struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

// DecideFloor applies a deterministic floor bump for DECISION_POINT_FLOOR.
func DecideFloor(ctx context.Context, in *DecisionRequest) (*DecisionResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if in == nil {
		return nil, fmt.Errorf("nil DecisionRequest")
	}
	point := strings.TrimSpace(in.Point)
	if point != "" && point != "DECISION_POINT_FLOOR" && point != "DECISION_POINT_UNSPECIFIED" && point != "2" {
		return nil, fmt.Errorf("only DECISION_POINT_FLOOR is supported (got %q)", point)
	}

	base := resolveBaseFloor(in)
	bump, err := decimal.NewFromString(DefaultBump)
	if err != nil {
		return nil, err
	}
	out := base.Add(bump)
	tier := "sample_floor_bump"

	if in.Policy != nil && in.Policy.MaxBid != "" {
		if cap, err := decimal.NewFromString(in.Policy.MaxBid); err == nil && !cap.IsZero() && out.GreaterThan(cap) {
			out = cap
			tier = "sample_floor_bump_clipped_max_bid"
		}
	}

	cpm := out.StringFixed(4)
	return &DecisionResponse{
		Price:      &Price{CPM: cpm, Accept: true},
		Confidence: 0.75,
		Rationale:  "sample decisioner: " + tier + " base=" + base.StringFixed(4) + " bump=" + DefaultBump + " cpm=" + cpm,
		Model:      &ModelRef{ID: "sample_floor_bump", Version: "v1"},
	}, nil
}

func resolveBaseFloor(in *DecisionRequest) decimal.Decimal {
	best := decimal.Zero
	consider := func(s string) {
		if s == "" {
			return
		}
		d, err := decimal.NewFromString(s)
		if err != nil || d.IsZero() {
			return
		}
		if d.GreaterThan(best) {
			best = d
		}
	}
	if in.Policy != nil {
		consider(in.Policy.Floor)
	}
	if in.Features != nil {
		for _, k := range []string{"deal_floor", "route_floor", "tenant_default_floor"} {
			consider(in.Features.Categorical[k])
		}
	}
	for _, c := range in.Candidates {
		if c.Attrs != nil {
			consider(c.Attrs["placement_floor"])
			consider(c.Attrs["pod_floor"])
		}
		consider(c.Price)
	}
	return best
}
