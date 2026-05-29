package analytics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func d(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestPeriod_Previous(t *testing.T) {
	tests := []struct {
		name           string
		from, to       string
		now            string
		wantFrom, want string
	}{
		{
			name: "single day -> yesterday",
			from: "2026-05-16", to: "2026-05-16", now: "2026-05-16",
			wantFrom: "2026-05-15", want: "2026-05-15",
		},
		{
			name: "equal-length 7 days",
			from: "2026-05-10", to: "2026-05-16", now: "2026-05-16",
			wantFrom: "2026-05-03", want: "2026-05-09",
		},
		{
			name: "MTD by today",
			from: "2026-05-01", to: "2026-05-24", now: "2026-05-24",
			wantFrom: "2026-04-01", want: "2026-04-24",
		},
		{
			name: "MTD by end of month",
			from: "2026-05-01", to: "2026-05-31", now: "2026-06-15",
			wantFrom: "2026-04-01", want: "2026-04-30",
		},
		{
			name: "MTD March -> Feb clamp (non-leap)",
			from: "2026-03-01", to: "2026-03-31", now: "2026-04-10",
			wantFrom: "2026-02-01", want: "2026-02-28",
		},
		{
			name: "MTD Jan -> prior Dec",
			from: "2026-01-01", to: "2026-01-15", now: "2026-01-15",
			wantFrom: "2025-12-01", want: "2025-12-15",
		},
		{
			name: "MTD day 31 today -> Feb clamp",
			from: "2026-03-01", to: "2026-03-31", now: "2026-03-31",
			wantFrom: "2026-02-01", want: "2026-02-28",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewPeriod(d(tt.from), d(tt.to))
			require.NoError(t, err)
			prev := p.Previous(d(tt.now))
			require.Equal(t, tt.wantFrom, prev.From.Format("2006-01-02"), "from")
			require.Equal(t, tt.want, prev.To.Format("2006-01-02"), "to")
		})
	}
}

func TestPeriod_EqualLengthMatchesDays(t *testing.T) {
	p, err := NewPeriod(d("2026-05-10"), d("2026-05-16"))
	require.NoError(t, err)
	prev := p.Previous(d("2026-05-20")) // now not today -> equal-length
	require.Equal(t, p.Days(), prev.Days())
}

func TestNewPeriod_InvalidWhenFromAfterTo(t *testing.T) {
	_, err := NewPeriod(d("2026-05-16"), d("2026-05-10"))
	require.ErrorIs(t, err, ErrInvalidPeriod)
}

func TestPeriod_DaysAndDates(t *testing.T) {
	p, err := NewPeriod(d("2026-05-01"), d("2026-05-03"))
	require.NoError(t, err)
	require.Equal(t, 3, p.Days())
	require.Len(t, p.Dates(), 3)
	require.Equal(t, "2026-05-01", p.Dates()[0].Format("2006-01-02"))
	require.Equal(t, "2026-05-03", p.Dates()[2].Format("2006-01-02"))
}
