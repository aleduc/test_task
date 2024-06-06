package pagination

type OrderDirection string

type Order struct {
	Field     string `schema:"sort_by" query:"sort_by" json:"sort_by"`
	Direction string `schema:"order_by" query:"order_by" json:"order_by"`
}
