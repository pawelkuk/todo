package cli

import (
	"fmt"
	"testing"
	"time"
)

func TestParseDue(t *testing.T) {
	tt := []struct {
		in      string
		want    time.Duration
		wantErr bool
	}{
		{
			in:      "1h",
			want:    time.Hour,
			wantErr: false,
		},
		{
			in:      "11h",
			want:    11 * time.Hour,
			wantErr: false,
		},
		{
			in:      "1h1d",
			want:    25 * time.Hour,
			wantErr: false,
		},
		{
			in:      "1h1d1w",
			want:    7*24*time.Hour + 25*time.Hour,
			wantErr: false,
		},
		{
			in:      "1w1m",
			want:    (7 + 30) * 24 * time.Hour,
			wantErr: false,
		},
		{
			in:      "1m1y",
			want:    (30 + 365) * 24 * time.Hour,
			wantErr: false,
		},
		{
			in:      "0h0d0w0m0y",
			want:    0,
			wantErr: false,
		},
		{
			in:      "99h0d0w100m0y",
			want:    99*time.Hour + 100*30*24*time.Hour,
			wantErr: false,
		},
	}
	for idx, test := range tt {
		t.Run(fmt.Sprintf("test_%d", idx), func(t *testing.T) {
			got, err := parseDue(test.in)
			if err != nil && !test.wantErr {
				t.Errorf("got err %v; expected nil", err)
			}
			if err == nil && test.wantErr {
				t.Errorf("wanted error; err==nil; got: %v", got)
			}
			if got != test.want {
				t.Errorf("got: %v, want %v", got, test.want)
			}
		})
	}
}
