package amocrm

// EntityType represents the type of entity
type EntityType string

const (
	EntityTypeContact  EntityType = "contacts"
	EntityTypeCompany  EntityType = "companies"
	EntityTypeLead     EntityType = "leads"
	EntityTypeCustomer EntityType = "customers"
)

// CustomFieldValue represents a custom field value
type CustomFieldValue struct {
	FieldID   int          `json:"field_id"`
	FieldName string       `json:"field_name,omitempty"`
	FieldCode string       `json:"field_code,omitempty"`
	FieldType string       `json:"field_type,omitempty"`
	Values    []FieldValue `json:"values"`
}

// FieldValue represents a single value in a custom field
type FieldValue struct {
	Value    interface{} `json:"value"`
	EnumID   int         `json:"enum_id,omitempty"`
	EnumCode string      `json:"enum_code,omitempty"`
	Enum     string      `json:"enum,omitempty"`
}

// EmbeddedTags represents tags in embedded format
type EmbeddedTags struct {
	Tags []Tag `json:"tags,omitempty"`
}

// Tag represents a tag
type Tag struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

// Links represents entity links
type Links struct {
	Self Link `json:"self,omitempty"`
	Next Link `json:"next,omitempty"`
	Prev Link `json:"prev,omitempty"`
}

// Link represents a single link
type Link struct {
	Href string `json:"href,omitempty"`
}

// HasNext returns true if there is a next page
func (l Links) HasNext() bool {
	return l.Next.Href != ""
}

// Embedded represents common embedded data
type Embedded struct {
	Tags      []Tag       `json:"tags,omitempty"`
	Companies []Company   `json:"companies,omitempty"`
	Contacts  []Contact   `json:"contacts,omitempty"`
	Leads     []Lead      `json:"leads,omitempty"`
	Catalog   interface{} `json:"catalog_elements,omitempty"`
}
