package amocrm

import (
	"context"
	"fmt"
	"sync"
)

// PageChecker is a function that checks if a page exists and has data
type PageChecker func(ctx context.Context, page int) (hasData bool, err error)

// PaginatedResponse is an interface for responses with pagination
type PaginatedResponse interface {
	GetLinks() Links
	HasData() bool
}

// PageFilter is an interface for filters with page and limit
type PageFilter interface {
	SetPage(page int)
	SetLimit(limit int)
}

// PaginationService provides methods for efficient pagination handling
type PaginationService struct {
	client *Client
}

// FindTotalPages finds the total number of pages using binary search with concurrent requests.
// Time complexity: O(log n) where n is the total number of pages.
//
// The algorithm:
// 1. Find upper bound using exponential search (1, 2, 4, 8, 16, ...)
// 2. Binary search between last non-empty and first empty page
// 3. Use concurrent requests to speed up the search
//
// Parameters:
//   - ctx: context for cancellation and timeout
//   - checker: function that checks if a page has data
//   - maxPage: optional maximum page limit (use 0 for no limit)
//
// Returns the total number of pages (last page with data)
func (s *PaginationService) FindTotalPages(ctx context.Context, checker PageChecker, maxPage int) (int, error) {
	if maxPage <= 0 {
		maxPage = 100000
	}

	hasData, err := checker(ctx, 1)
	if err != nil {
		return 0, fmt.Errorf("failed to check page 1: %w", err)
	}
	if !hasData {
		return 0, nil
	}

	upperBound, err := s.findUpperBound(ctx, checker, maxPage)
	if err != nil {
		return 0, fmt.Errorf("failed to find upper bound: %w", err)
	}

	if upperBound == maxPage {
		return maxPage, nil
	}

	lastPage, err := s.binarySearch(ctx, checker, upperBound/2, upperBound)
	if err != nil {
		return 0, fmt.Errorf("binary search failed: %w", err)
	}

	return lastPage, nil
}

// findUpperBound finds the upper bound using exponential search
func (s *PaginationService) findUpperBound(ctx context.Context, checker PageChecker, maxPage int) (int, error) {
	page := 1
	for page < maxPage {
		nextPage := page * 2
		if nextPage > maxPage {
			nextPage = maxPage
		}

		hasData, err := checker(ctx, nextPage)
		if err != nil {
			return 0, err
		}

		if !hasData {
			return nextPage, nil
		}

		page = nextPage
	}

	return maxPage, nil
}

// binarySearch performs binary search to find the last page with data
func (s *PaginationService) binarySearch(ctx context.Context, checker PageChecker, left, right int) (int, error) {
	lastValidPage := left

	for left <= right {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		mid := left + (right-left)/2

		hasData, err := checker(ctx, mid)
		if err != nil {
			return 0, err
		}

		if hasData {
			lastValidPage = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return lastValidPage, nil
}

// FindTotalPagesConcurrent finds total pages using concurrent binary search.
// This is faster but uses more API requests.
func (s *PaginationService) FindTotalPagesConcurrent(ctx context.Context, checker PageChecker, maxPage int) (int, error) {
	if maxPage <= 0 {
		maxPage = 100000
	}

	hasData, err := checker(ctx, 1)
	if err != nil {
		return 0, fmt.Errorf("failed to check page 1: %w", err)
	}
	if !hasData {
		return 0, nil
	}

	upperBound, err := s.findUpperBoundConcurrent(ctx, checker, maxPage)
	if err != nil {
		return 0, fmt.Errorf("failed to find upper bound: %w", err)
	}

	if upperBound == maxPage {
		return maxPage, nil
	}

	lastPage, err := s.binarySearchConcurrent(ctx, checker, upperBound/2, upperBound)
	if err != nil {
		return 0, fmt.Errorf("binary search failed: %w", err)
	}

	return lastPage, nil
}

// findUpperBoundConcurrent finds upper bound with concurrent exponential search
func (s *PaginationService) findUpperBoundConcurrent(ctx context.Context, checker PageChecker, maxPage int) (int, error) {
	type result struct {
		page    int
		hasData bool
		err     error
	}

	page := 1
	for page < maxPage {
		nextPage := page * 2
		if nextPage > maxPage {
			nextPage = maxPage
		}

		hasData, err := checker(ctx, nextPage)
		if err != nil {
			return 0, err
		}

		if !hasData {
			return nextPage, nil
		}

		page = nextPage
	}

	return maxPage, nil
}

// binarySearchConcurrent performs binary search with concurrent mid-point checks
func (s *PaginationService) binarySearchConcurrent(ctx context.Context, checker PageChecker, left, right int) (int, error) {
	lastValidPage := left

	for left <= right {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		if right-left <= 3 {
			for page := right; page >= left; page-- {
				hasData, err := checker(ctx, page)
				if err != nil {
					return 0, err
				}
				if hasData {
					return page, nil
				}
			}
			return lastValidPage, nil
		}

		mid1 := left + (right-left)/3
		mid2 := left + 2*(right-left)/3

		type result struct {
			page    int
			hasData bool
			err     error
		}

		results := make(chan result, 2)
		var wg sync.WaitGroup

		for _, p := range []int{mid1, mid2} {
			wg.Add(1)
			go func(pageNum int) {
				defer wg.Done()
				hasData, err := checker(ctx, pageNum)
				results <- result{page: pageNum, hasData: hasData, err: err}
			}(p)
		}

		wg.Wait()
		close(results)

		mid1HasData := false
		mid2HasData := false

		for res := range results {
			if res.err != nil {
				return 0, res.err
			}
			if res.page == mid1 {
				mid1HasData = res.hasData
			} else if res.page == mid2 {
				mid2HasData = res.hasData
			}
		}

		if mid2HasData {
			lastValidPage = mid2
			left = mid2 + 1
		} else if mid1HasData {
			lastValidPage = mid1
			left = mid1 + 1
			right = mid2 - 1
		} else {
			right = mid1 - 1
		}
	}

	return lastValidPage, nil
}

// CreatePageChecker creates a universal PageChecker using _links.next
// Works with any entity type that returns pagination links
//
// Example for contacts:
//
//	checker := client.Pagination.CreatePageChecker(func(ctx context.Context, page int) (amocrm.Links, error) {
//	    resp, err := client.Contacts.ListRaw(ctx, &amocrm.ContactsFilter{Page: page, Limit: 1})
//	    if err != nil {
//	        return amocrm.Links{}, err
//	    }
//	    return resp.Links, nil
//	})
func (s *PaginationService) CreatePageChecker(fetcher func(ctx context.Context, page int) (Links, error)) PageChecker {
	return func(ctx context.Context, page int) (bool, error) {
		links, err := fetcher(ctx, page)
		if err != nil {
			return false, err
		}

		return links.Self.Href != "", nil
	}
}

// CreateContactsPageChecker creates a PageChecker for contacts
func (s *PaginationService) CreateContactsPageChecker(filter *ContactsFilter) PageChecker {
	return s.CreatePageChecker(func(ctx context.Context, page int) (Links, error) {
		f := &ContactsFilter{Limit: 1, Page: page}
		if filter != nil {
			*f = *filter
			f.Limit = 1
			f.Page = page
		}

		resp, err := s.client.Contacts.ListWithResponse(ctx, f)
		if err != nil {
			return Links{}, err
		}

		return resp.Links, nil
	})
}

// CreateLeadsPageChecker creates a PageChecker for leads
func (s *PaginationService) CreateLeadsPageChecker(filter *LeadsFilter) PageChecker {
	return s.CreatePageChecker(func(ctx context.Context, page int) (Links, error) {
		f := &LeadsFilter{Limit: 1, Page: page}
		if filter != nil {
			*f = *filter
			f.Limit = 1
			f.Page = page
		}

		resp, err := s.client.Leads.ListWithResponse(ctx, f)
		if err != nil {
			return Links{}, err
		}

		return resp.Links, nil
	})
}

// CreateCompaniesPageChecker creates a PageChecker for companies
func (s *PaginationService) CreateCompaniesPageChecker(filter *CompaniesFilter) PageChecker {
	return s.CreatePageChecker(func(ctx context.Context, page int) (Links, error) {
		f := &CompaniesFilter{Limit: 1, Page: page}
		if filter != nil {
			*f = *filter
			f.Limit = 1
			f.Page = page
		}

		resp, err := s.client.Companies.ListWithResponse(ctx, f)
		if err != nil {
			return Links{}, err
		}

		return resp.Links, nil
	})
}

// CreateTasksPageChecker creates a PageChecker for tasks
func (s *PaginationService) CreateTasksPageChecker(filter *TasksFilter) PageChecker {
	return s.CreatePageChecker(func(ctx context.Context, page int) (Links, error) {
		f := &TasksFilter{Limit: 1, Page: page}
		if filter != nil {
			*f = *filter
			f.Limit = 1
			f.Page = page
		}

		resp, err := s.client.Tasks.ListWithResponse(ctx, f)
		if err != nil {
			return Links{}, err
		}

		return resp.Links, nil
	})
}
