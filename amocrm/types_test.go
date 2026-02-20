package amocrm

import (
	"testing"
)

func TestLinks_HasNext(t *testing.T) {
	tests := []struct {
		name     string
		links    Links
		expected bool
	}{
		{
			name: "Has next link",
			links: Links{
				Self: Link{Href: "https://example.amocrm.ru/api/v4/contacts?page=1"},
				Next: Link{Href: "https://example.amocrm.ru/api/v4/contacts?page=2"},
			},
			expected: true,
		},
		{
			name: "No next link - empty string",
			links: Links{
				Self: Link{Href: "https://example.amocrm.ru/api/v4/contacts?page=5"},
				Next: Link{Href: ""},
			},
			expected: false,
		},
		{
			name: "No next link - last page",
			links: Links{
				Self: Link{Href: "https://example.amocrm.ru/api/v4/contacts?page=10"},
			},
			expected: false,
		},
		{
			name: "Has prev but no next",
			links: Links{
				Self: Link{Href: "https://example.amocrm.ru/api/v4/contacts?page=10"},
				Prev: Link{Href: "https://example.amocrm.ru/api/v4/contacts?page=9"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.links.HasNext()
			if result != tt.expected {
				t.Errorf("HasNext() = %v, want %v", result, tt.expected)
			}
		})
	}
}
