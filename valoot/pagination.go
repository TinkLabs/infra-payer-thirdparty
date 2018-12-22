package valoot

type Pagination struct {
	Total       int   `json:"total"`
	Count       int   `json:"count"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	Links       Links `json:"links"`
}

type Links struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

func (p Pagination) HasNext() bool {
	return p.CurrentPage < p.TotalPages
}

func (p Pagination) GetNextPage() int {
	return p.CurrentPage + 1
}

func (p Pagination) HasPrevious() bool {
	return p.CurrentPage > 1
}

func (p Pagination) GetPreviousPage() int {
	return p.CurrentPage - 1
}
