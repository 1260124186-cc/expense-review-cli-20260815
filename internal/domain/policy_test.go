package domain

import "testing"

func TestDefaultPolicyProvidesCategoryCaps(t *testing.T) {
	if got, want := DefaultPolicy().capFor("meals"), int64(7_500); got != want {
		t.Fatalf("meal cap = %d, want %d", got, want)
	}
}
