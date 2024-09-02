package models

type PageRequest struct {
	Page         int
	ItemsPerPage int
}

type PageData struct {
	CurrentPage int
	PageCount   int
	ItemCount   int
}

type PaginatedLinks struct {
	PageData
	Links []ShortLink
}

type PaginatedUsers struct {
	PageData
	Users []User
}
