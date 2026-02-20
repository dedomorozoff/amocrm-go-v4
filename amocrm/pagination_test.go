package amocrm

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
)

func TestFindTotalPages(t *testing.T) {
	tests := []struct {
		name          string
		totalPages    int
		maxPage       int
		expectedCalls int
	}{
		{
			name:          "Empty result",
			totalPages:    0,
			maxPage:       1000,
			expectedCalls: 1,
		},
		{
			name:          "Single page",
			totalPages:    1,
			maxPage:       1000,
			expectedCalls: 2,
		},
		{
			name:          "10 pages",
			totalPages:    10,
			maxPage:       1000,
			expectedCalls: 7,
		},
		{
			name:          "100 pages",
			totalPages:    100,
			maxPage:       1000,
			expectedCalls: 13,
		},
		{
			name:          "1000 pages",
			totalPages:    1000,
			maxPage:       10000,
			expectedCalls: 20,
		},
		{
			name:          "500 pages",
			totalPages:    500,
			maxPage:       10000,
			expectedCalls: 18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var callCount int32

			checker := func(ctx context.Context, page int) (bool, error) {
				atomic.AddInt32(&callCount, 1)
				return page <= tt.totalPages, nil
			}

			service := &PaginationService{client: nil}

			result, err := service.FindTotalPages(context.Background(), checker, tt.maxPage)
			if err != nil {
				t.Fatalf("FindTotalPages() error = %v", err)
			}

			if result != tt.totalPages {
				t.Errorf("FindTotalPages() = %v, want %v", result, tt.totalPages)
			}

			calls := int(atomic.LoadInt32(&callCount))
			t.Logf("Total API calls: %d (expected ~%d, actual pages: %d)", calls, tt.expectedCalls, tt.totalPages)

			if calls > tt.expectedCalls+5 {
				t.Errorf("Too many API calls: %d, expected around %d", calls, tt.expectedCalls)
			}
		})
	}
}

func TestFindTotalPagesConcurrent(t *testing.T) {
	tests := []struct {
		name       string
		totalPages int
		maxPage    int
	}{
		{
			name:       "10 pages concurrent",
			totalPages: 10,
			maxPage:    1000,
		},
		{
			name:       "100 pages concurrent",
			totalPages: 100,
			maxPage:    1000,
		},
		{
			name:       "500 pages concurrent",
			totalPages: 500,
			maxPage:    10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var callCount int32

			checker := func(ctx context.Context, page int) (bool, error) {
				atomic.AddInt32(&callCount, 1)
				return page <= tt.totalPages, nil
			}

			service := &PaginationService{client: nil}

			result, err := service.FindTotalPagesConcurrent(context.Background(), checker, tt.maxPage)
			if err != nil {
				t.Fatalf("FindTotalPagesConcurrent() error = %v", err)
			}

			if result != tt.totalPages {
				t.Errorf("FindTotalPagesConcurrent() = %v, want %v", result, tt.totalPages)
			}

			calls := int(atomic.LoadInt32(&callCount))
			t.Logf("Total API calls (concurrent): %d (actual pages: %d)", calls, tt.totalPages)
		})
	}
}

func TestFindTotalPagesWithError(t *testing.T) {
	checker := func(ctx context.Context, page int) (bool, error) {
		return false, fmt.Errorf("simulated error at page %d", page)
	}

	service := &PaginationService{client: nil}

	_, err := service.FindTotalPages(context.Background(), checker, 1000)
	if err == nil {
		t.Error("Expected error but got nil")
	}
	if err != nil {
		t.Logf("Got expected error: %v", err)
	}
}

func TestFindTotalPagesWithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	checker := func(ctx context.Context, page int) (bool, error) {
		return page <= 100, nil
	}

	service := &PaginationService{client: nil}

	_, err := service.FindTotalPages(ctx, checker, 1000)
	if err == nil {
		t.Error("Expected context cancellation error but got nil")
	} else {
		t.Logf("Got expected cancellation error: %v", err)
	}
}

func BenchmarkFindTotalPages(b *testing.B) {
	benchmarks := []struct {
		name       string
		totalPages int
	}{
		{"10 pages", 10},
		{"100 pages", 100},
		{"1000 pages", 1000},
		{"5000 pages", 5000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			checker := func(ctx context.Context, page int) (bool, error) {
				return page <= bm.totalPages, nil
			}

			service := &PaginationService{client: nil}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := service.FindTotalPages(context.Background(), checker, 100000)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkFindTotalPagesConcurrent(b *testing.B) {
	benchmarks := []struct {
		name       string
		totalPages int
	}{
		{"10 pages", 10},
		{"100 pages", 100},
		{"1000 pages", 1000},
		{"5000 pages", 5000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			checker := func(ctx context.Context, page int) (bool, error) {
				return page <= bm.totalPages, nil
			}

			service := &PaginationService{client: nil}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := service.FindTotalPagesConcurrent(context.Background(), checker, 100000)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
