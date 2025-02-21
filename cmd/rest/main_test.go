package main

import "testing"

func TestIntMin(t *testing.T) {

	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{
			name: "a less than b",
			a:    0,
			b:    1,
			want: 0,
		},
		{
			name: "b less than a",
			a:    1,
			b:    0,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IntMin(tt.a, tt.b)

			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}
