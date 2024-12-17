package response

type ShoeboxListItemsResponse struct {
	Items []*ShoeboxListItem `json:"items"`
}

type ShoeboxListItem struct {
	Id           int64 `json:"id"`
	ItemId       int64 `json:"item_id"`
	Status       uint8 `json:"status"`
	Created      int64 `json:"created"`
	LastModified int64 `json:"lastmodified"`
	TypeId       uint8 `json:"type_id"`
}
