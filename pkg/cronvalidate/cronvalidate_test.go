package cronvalidate

import "testing"

func TestValidateExprStrictCron(t *testing.T) {
	tests := []struct {
		expr string
		ok   bool
	}{
		{"0 0 * * *", true},
		{"@daily", false},
		{"@every 1h", false},
		{"CRON_TZ=UTC 0 0 * * *", true},
	}
	for _, tc := range tests {
		err := ValidateExpr(tc.expr)
		got := err == nil
		if got != tc.ok {
			t.Fatalf("expr=%q want ok=%v got err=%v", tc.expr, tc.ok, err)
		}
	}
}
