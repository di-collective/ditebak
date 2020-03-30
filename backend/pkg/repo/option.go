package repo

// FindOptions ...
type FindOptions struct {
	Page           int
	Size           int
	IncludeRemoved bool
	Sort           *SortOption
	Params         map[string]interface{}
}

// SortOption ...SortOption
type SortOption struct {
	Field      string
	Descending bool
}

// Skip calculate skip based on page and size
func (o FindOptions) Skip() int {
	p := o.Page
	if p < 1 {
		p = 1
	}

	return (p - 1) * o.Size
}
