package main

import (
	"strings"
	"testing"

	"github.com/vshulcz/deja-vu/internal/bench"
)

// The gate this benchmark exists to be: an update that only added bytes must
// not rewrite the records already on file. That is what the new-transcript path
// got wrong for two months (#3500), and wall time on a shared runner is too
// noisy to catch it — the bytes are not.
func TestBenchIngestOnlyARewriteRewritesTheStore(t *testing.T) {
	hermeticEnv(t)
	report, err := measureIngest(bench.Seed)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Classes) != 4 {
		t.Fatalf("measured %d update classes, want the four a store sees: %#v", len(report.Classes), report.Classes)
	}
	for _, c := range report.Classes {
		rewriteExpected := strings.HasPrefix(c.Name, "rewritten")
		if c.Rewritten != rewriteExpected {
			t.Errorf("%q rewrote the store = %v, want %v", c.Name, c.Rewritten, rewriteExpected)
		}
		if c.RecordsMB <= 0 {
			t.Errorf("%q left no records behind, so the pass measured nothing", c.Name)
		}
	}
	if added := report.Classes[2].AddedKB; added <= 0 {
		t.Errorf("a new transcript added %.2f KB of records, so it was not indexed", added)
	}
}
