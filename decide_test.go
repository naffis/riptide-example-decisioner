package main

import (
	"context"
	"testing"
	"time"
)

func TestDecideFloor_bumpsBase(t *testing.T) {
	out, err := DecideFloor(context.Background(), &DecisionRequest{
		Point:  "DECISION_POINT_FLOOR",
		Policy: &TenantPolicy{Floor: "2.0000"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Price.CPM != "2.2500" {
		t.Fatalf("cpm=%s want 2.2500", out.Price.CPM)
	}
	if !out.Price.Accept {
		t.Fatal("expected accept")
	}
}

func TestDecideFloor_clipsMaxBid(t *testing.T) {
	out, err := DecideFloor(context.Background(), &DecisionRequest{
		Point:  "DECISION_POINT_FLOOR",
		Policy: &TenantPolicy{Floor: "2.0000", MaxBid: "2.1000"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Price.CPM != "2.1000" {
		t.Fatalf("cpm=%s want 2.1000", out.Price.CPM)
	}
}

func TestDecideFloor_honorsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := DecideFloor(ctx, &DecisionRequest{Point: "DECISION_POINT_FLOOR", Policy: &TenantPolicy{Floor: "1.0000"}})
	if err == nil {
		t.Fatal("expected canceled context error")
	}
}

func TestDecideFloor_rejectsWrongPoint(t *testing.T) {
	_, err := DecideFloor(context.Background(), &DecisionRequest{Point: "DECISION_POINT_BID"})
	if err == nil {
		t.Fatal("expected wrong-point error")
	}
}

func TestDecideFloor_deadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(time.Millisecond)
	_, err := DecideFloor(ctx, &DecisionRequest{Point: "DECISION_POINT_FLOOR", Policy: &TenantPolicy{Floor: "1.0000"}})
	if err == nil {
		t.Fatal("expected deadline error")
	}
}
