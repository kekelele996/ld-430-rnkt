package util

import "testing"

func TestPaginationHelpers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		page     int64
		pageSize int64
		wantPage int64
		wantSize int64
		wantOff  int64
	}{
		{name: "normal", page: 2, pageSize: 20, wantPage: 2, wantSize: 20, wantOff: 20},
		{name: "zero page", page: 0, pageSize: 20, wantPage: 1, wantSize: 20, wantOff: 0},
		{name: "negative page", page: -3, pageSize: 20, wantPage: 1, wantSize: 20, wantOff: 0},
		{name: "zero size", page: 3, pageSize: 0, wantPage: 3, wantSize: 20, wantOff: 40},
		{name: "oversize", page: 2, pageSize: 500, wantPage: 2, wantSize: 100, wantOff: 100},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NormalizePage(tt.page); got != tt.wantPage {
				t.Fatalf("NormalizePage(%d) = %d, want %d", tt.page, got, tt.wantPage)
			}
			if got := NormalizePageSize(tt.pageSize); got != tt.wantSize {
				t.Fatalf("NormalizePageSize(%d) = %d, want %d", tt.pageSize, got, tt.wantSize)
			}
			if got := Offset(tt.page, tt.pageSize); got != tt.wantOff {
				t.Fatalf("Offset(%d,%d) = %d, want %d", tt.page, tt.pageSize, got, tt.wantOff)
			}
		})
	}
}
