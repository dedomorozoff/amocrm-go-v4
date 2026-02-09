# Pagination Service

The Pagination Service provides efficient methods to find the total number of pages in AmoCRM API responses using binary search with O(log n) time complexity.

## Problem

AmoCRM API doesn't return the total number of pages in pagination responses. To find the last page, you would typically need to iterate through all pages sequentially, which takes O(n) time where n is the total number of pages.

## Solution

This library implements two algorithms:

1. **FindTotalPages** - Sequential binary search (O(log n))
2. **FindTotalPagesConcurrent** - Concurrent binary search (faster but uses more API calls)

Both algorithms use:
- Exponential search to find upper bound (1, 2, 4, 8, 16, ...)
- Binary search to find exact last page

## Usage

### Basic Example - Contacts

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ALipckin/amocrm-go-v4/amocrm"
)

func main() {
    client := amocrm.NewClient(
        amocrm.WithSubdomain("your-subdomain"),
        amocrm.WithPermanentToken("your-token"),
    )

    ctx := context.Background()

    // Create page checker for contacts
    filter := &amocrm.ContactsFilter{
        Limit: 250,
    }
    checker := client.Pagination.CreateContactsPageChecker(filter)

    // Find total pages
    totalPages, err := client.Pagination.FindTotalPages(ctx, checker, 10000)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Total pages: %d\n", totalPages)
    fmt.Printf("Total contacts (approx): %d\n", totalPages*250)
}
```

### Alternative: Iterate Until No Next Link

You can also iterate through pages using `_links.next` without finding total pages upfront:

```go
func fetchAllContacts(client *amocrm.Client) ([]amocrm.Contact, error) {
    ctx := context.Background()
    allContacts := []amocrm.Contact{}
    
    filter := &amocrm.ContactsFilter{Limit: 250, Page: 1}
    
    for {
        resp, err := client.Contacts.ListWithResponse(ctx, filter)
        if err != nil {
            return nil, err
        }
        
        allContacts = append(allContacts, resp.Embedded.Contacts...)
        
        // Check if there's a next page using _links.next
        if !resp.Links.HasNext() {
            break
        }
        
        filter.Page++
    }
    
    return allContacts, nil
}
```

### Example - Leads with Filter

```go
// Find total pages for leads in specific pipeline
filter := &amocrm.LeadsFilter{
    Limit:      250,
    PipelineID: 123,
}

checker := client.Pagination.CreateLeadsPageChecker(filter)
totalPages, err := client.Pagination.FindTotalPages(ctx, checker, 10000)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total pages in pipeline: %d\n", totalPages)
```

### Example - Concurrent Version

```go
// Faster but uses more API calls
filter := &amocrm.ContactsFilter{Limit: 250}
checker := client.Pagination.CreateContactsPageChecker(filter)

totalPages, err := client.Pagination.FindTotalPagesConcurrent(ctx, checker, 10000)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total pages (concurrent): %d\n", totalPages)
```

### Example - Universal Custom Page Checker

The library now provides a universal `CreatePageChecker` that works with any entity by checking `_links` in the response:

```go
// Universal approach - works for any entity
checker := client.Pagination.CreatePageChecker(func(ctx context.Context, page int) (amocrm.Links, error) {
    // For incomplete tasks
    completed := false
    filter := &amocrm.TasksFilter{
        Page:        page,
        Limit:       1,
        IsCompleted: &completed,
    }

    resp, err := client.Tasks.ListWithResponse(ctx, filter)
    if err != nil {
        return amocrm.Links{}, err
    }

    return resp.Links, nil
})

totalPages, err := client.Pagination.FindTotalPages(ctx, checker, 10000)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total pages of incomplete tasks: %d\n", totalPages)
```

### Example - Legacy Style (Still Supported)

```go
// Old style - also works, checks if data exists
customChecker := func(ctx context.Context, page int) (bool, error) {
    completed := false
    filter := &amocrm.TasksFilter{
        Page:        page,
        Limit:       1,
        IsCompleted: &completed,
    }

    tasks, err := client.Tasks.List(ctx, filter)
    if err != nil {
        return false, err
    }

    return len(tasks) > 0, nil
}

totalPages, err := client.Pagination.FindTotalPages(ctx, customChecker, 10000)
```

## Performance

### Time Complexity

- **FindTotalPages**: O(log n) where n is total pages
- **FindTotalPagesConcurrent**: O(log n) but parallelized

### API Calls

For a dataset with N pages:

| Total Pages | Sequential Iterations | Binary Search Calls | Improvement |
|-------------|----------------------|---------------------|-------------|
| 10          | 10                   | ~9                  | 1.1x        |
| 100         | 100                  | ~14                 | 7x          |
| 1,000       | 1,000                | ~20                 | 50x         |
| 5,000       | 5,000                | ~23                 | 217x        |
| 10,000      | 10,000               | ~24                 | 416x        |

### Benchmark Results

```
BenchmarkFindTotalPages/10_pages-8          500000    2.5 µs/op
BenchmarkFindTotalPages/100_pages-8         300000    4.2 µs/op
BenchmarkFindTotalPages/1000_pages-8        200000    6.1 µs/op
BenchmarkFindTotalPages/5000_pages-8        150000    7.3 µs/op
```

## Built-in Page Checkers

The library provides page checkers for common entities:

- `CreateContactsPageChecker(filter *ContactsFilter)`
- `CreateLeadsPageChecker(filter *LeadsFilter)`
- `CreateCompaniesPageChecker(filter *CompaniesFilter)`
- `CreateTasksPageChecker(filter *TasksFilter)`

## API Reference

### FindTotalPages

```go
func (s *PaginationService) FindTotalPages(
    ctx context.Context,
    checker PageChecker,
    maxPage int,
) (int, error)
```

Finds total number of pages using sequential binary search.

**Parameters:**
- `ctx` - context for cancellation and timeout
- `checker` - function that checks if a page has data
- `maxPage` - maximum page limit (use 0 for default 100000)

**Returns:**
- Total number of pages (last page with data)
- Error if any

### FindTotalPagesConcurrent

```go
func (s *PaginationService) FindTotalPagesConcurrent(
    ctx context.Context,
    checker PageChecker,
    maxPage int,
) (int, error)
```

Finds total number of pages using concurrent binary search. Faster but uses more API calls.

**Parameters:**
- Same as FindTotalPages

**Returns:**
- Same as FindTotalPages

### PageChecker

```go
type PageChecker func(ctx context.Context, page int) (hasData bool, err error)
```

Function type that checks if a specific page has data.

**Parameters:**
- `ctx` - context
- `page` - page number to check

**Returns:**
- `true` if page has data, `false` if empty
- Error if request fails

## Best Practices

1. **Use appropriate limit**: Set `Limit` to max value (250) to minimize total pages
2. **Set reasonable maxPage**: Prevents infinite loops on edge cases
3. **Handle context timeout**: Use context with timeout for long operations
4. **Choose right algorithm**:
   - Use `FindTotalPages` for most cases (fewer API calls)
   - Use `FindTotalPagesConcurrent` when speed is critical
5. **Reuse checkers**: Create checker once, use multiple times

## Error Handling

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

totalPages, err := client.Pagination.FindTotalPages(ctx, checker, 10000)
if err != nil {
    if err == context.DeadlineExceeded {
        log.Printf("Timeout: operation took too long")
    } else if apiErr, ok := err.(*amocrm.APIError); ok {
        log.Printf("API error: %v", apiErr)
    } else {
        log.Printf("Error: %v", err)
    }
    return
}
```

## Integration Example

Complete example showing how to use pagination to fetch all data:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ALipckin/amocrm-go-v4/amocrm"
)

func fetchAllContacts(client *amocrm.Client) ([]amocrm.Contact, error) {
    ctx := context.Background()

    // Find total pages
    filter := &amocrm.ContactsFilter{Limit: 250}
    checker := client.Pagination.CreateContactsPageChecker(filter)
    
    totalPages, err := client.Pagination.FindTotalPages(ctx, checker, 10000)
    if err != nil {
        return nil, fmt.Errorf("failed to find total pages: %w", err)
    }

    fmt.Printf("Found %d pages to fetch\n", totalPages)

    // Fetch all pages
    allContacts := make([]amocrm.Contact, 0, totalPages*250)
    
    for page := 1; page <= totalPages; page++ {
        filter.Page = page
        contacts, err := client.Contacts.List(ctx, filter)
        if err != nil {
            return nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
        }
        
        allContacts = append(allContacts, contacts...)
        fmt.Printf("Fetched page %d/%d (%d contacts)\n", 
            page, totalPages, len(contacts))
    }

    return allContacts, nil
}

func main() {
    client := amocrm.NewClient(
        amocrm.WithSubdomain("your-subdomain"),
        amocrm.WithPermanentToken("your-token"),
    )

    contacts, err := fetchAllContacts(client)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Total contacts fetched: %d\n", len(contacts))
}
```

## Notes

- The algorithm assumes pages are 1-indexed
- Empty pages (no data) indicate the end of pagination
- Rate limiting is handled by the client automatically
- All page checkers use `Limit: 1` for efficiency