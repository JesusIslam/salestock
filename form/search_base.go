package form

type SearchBase struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	OrderBy string `json:"order_by"`
}

func validateSearchBase(page, perPage int, orderBy string) (int, int, string, error) {
	var err error

	if perPage < 1 {
		perPage = 10
	}

	if page < 0 {
		page = 0
	}
	page = page * perPage

	if orderBy == "" {
		orderBy = "-_id"
	}

	return page, perPage, orderBy, err
}
