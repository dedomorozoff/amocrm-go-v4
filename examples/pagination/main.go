package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ALipckin/amocrm-go-v4/amocrm"
)

func main() {
	subdomain := os.Getenv("AMOCRM_SUBDOMAIN")
	token := os.Getenv("AMOCRM_TOKEN")

	if subdomain == "" || token == "" {
		log.Fatal("AMOCRM_SUBDOMAIN and AMOCRM_TOKEN environment variables are required")
	}

	client := amocrm.NewClient(
		amocrm.WithSubdomain(subdomain),
		amocrm.WithPermanentToken(token),
		amocrm.WithDebug(true),
	)

	ctx := context.Background()

	fmt.Println("=== Example 1: Find total pages for contacts ===")
	exampleContacts(ctx, client)

	fmt.Println("\n=== Example 2: Find total pages for leads with filter ===")
	exampleLeadsWithFilter(ctx, client)

	fmt.Println("\n=== Example 3: Find total pages using concurrent algorithm ===")
	exampleConcurrent(ctx, client)

	fmt.Println("\n=== Example 4: Custom page checker ===")
	exampleCustomChecker(ctx, client)

	fmt.Println("\n=== Example 5: Iterate using _links.next ===")
	exampleIterateWithLinks(ctx, client)
}

func exampleIterateWithLinks(ctx context.Context, client *amocrm.Client) {
	start := time.Now()

	var totalContacts int
	var pageCount int

	filter := &amocrm.ContactsFilter{
		Limit: 250,
		Page:  1,
	}

	for {
		resp, err := client.Contacts.ListWithResponse(ctx, filter)
		if err != nil {
			log.Printf("Error fetching page %d: %v", filter.Page, err)
			return
		}

		contactsOnPage := len(resp.Embedded.Contacts)
		totalContacts += contactsOnPage
		pageCount++

		fmt.Printf("Fetched page %d: %d contacts\n", filter.Page, contactsOnPage)

		// Check if there's a next page using _links.next
		if !resp.Links.HasNext() {
			break
		}

		filter.Page++

		// Safety limit for example
		if pageCount >= 10 {
			fmt.Println("Reached safety limit of 10 pages for demo")
			break
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Total contacts fetched: %d across %d pages\n", totalContacts, pageCount)
	fmt.Printf("Time taken: %v\n", elapsed)
	fmt.Printf("Note: This approach doesn't need to find total pages upfront\n")
}

func exampleContacts(ctx context.Context, client *amocrm.Client) {
	start := time.Now()

	filter := &amocrm.ContactsFilter{
		Limit: 250,
	}

	checker := client.Pagination.CreateContactsPageChecker(filter)

	totalPages, err := client.Pagination.FindTotalPages(ctx, checker, 10000)
	if err != nil {
		log.Printf("Error finding total pages: %v", err)
		return
	}

	elapsed := time.Since(start)
	fmt.Printf("Total pages for contacts: %d\n", totalPages)
	fmt.Printf("Time taken: %v\n", elapsed)
	fmt.Printf("Estimated API calls: ~%d (log2(%d) ≈ %.0f)\n",
		int(float64(totalPages)/250)+10,
		totalPages,
		float64(int(float64(totalPages)/250)+10))
}

func exampleLeadsWithFilter(ctx context.Context, client *amocrm.Client) {
	start := time.Now()

	filter := &amocrm.LeadsFilter{
		Limit:      250,
		PipelineID: 123, // Replace with your pipeline ID
	}

	checker := client.Pagination.CreateLeadsPageChecker(filter)

	totalPages, err := client.Pagination.FindTotalPages(ctx, checker, 10000)
	if err != nil {
		log.Printf("Error finding total pages: %v", err)
		return
	}

	elapsed := time.Since(start)
	fmt.Printf("Total pages for leads in pipeline %d: %d\n", filter.PipelineID, totalPages)
	fmt.Printf("Time taken: %v\n", elapsed)
}

func exampleConcurrent(ctx context.Context, client *amocrm.Client) {
	start := time.Now()

	filter := &amocrm.ContactsFilter{
		Limit: 250,
	}

	checker := client.Pagination.CreateContactsPageChecker(filter)

	totalPages, err := client.Pagination.FindTotalPagesConcurrent(ctx, checker, 10000)
	if err != nil {
		log.Printf("Error finding total pages: %v", err)
		return
	}

	elapsed := time.Since(start)
	fmt.Printf("Total pages for contacts (concurrent): %d\n", totalPages)
	fmt.Printf("Time taken: %v\n", elapsed)
	fmt.Printf("Note: Concurrent version is faster but uses more API requests\n")
}

func exampleCustomChecker(ctx context.Context, client *amocrm.Client) {
	start := time.Now()

	// Universal approach using CreatePageChecker - checks _links in response
	customChecker := client.Pagination.CreatePageChecker(func(ctx context.Context, page int) (amocrm.Links, error) {
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

	totalPages, err := client.Pagination.FindTotalPages(ctx, customChecker, 10000)
	if err != nil {
		log.Printf("Error finding total pages: %v", err)
		return
	}

	elapsed := time.Since(start)
	fmt.Printf("Total pages for incomplete tasks: %d\n", totalPages)
	fmt.Printf("Time taken: %v\n", elapsed)
	fmt.Printf("Note: Universal approach works with any entity type\n")
}
