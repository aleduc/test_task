package pagination

type Pagination struct {
	Number int `schema:"page" json:"page" query:"page"`
	Size   int `schema:"size" json:"size" query:"size"`
}

type PageQueryParams struct {
	Pagination
	Order
}

func CalculateOffset(page, size int) int {
	if page <= 0 || size <= 0 {
		return 0
	}
	return (page - 1) * size
}

func CalculateLastPage(count, pageSize int) int {
	if count <= pageSize || pageSize < 1 {
		return 1
	}
	if count%pageSize > 0 {
		return count/pageSize + 1
	}
	return count / pageSize
}
