package response

// ShoeboxListItemsResponse defines the reponse object returned by the `sfomuseum.you.shoebox.listItems` API method.
type ShoeboxListItemsResponse struct {
	// Zero or `ShoeboxListItem` instances.
	Items []*ShoeboxListItem `json:"items"`
}

// ShoeboxListItem defines an individual (shoebox) item returned by the `sfomuseum.you.shoebox.listItems` API method.
type ShoeboxListItem struct {
	// The unique identifier for the shoebox item.
	Id int64 `json:"id"`
	// The unique identifier of the item added to a person's shoebox.
	ItemId int64 `json:"item_id"`
	// The numeric status code of the shoebox item.
	Status uint8 `json:"status"`
	// The Unix timestamp when the item was added the shoebox.
	Created int64 `json:"created"`
	// The Unix timestamp when the shoebox item was last modified.
	LastModified int64 `json:"lastmodified"`
	// TypeId is the numeric identifier of the type of item added to the shoebox.
	TypeId uint8 `json:"type_id"`
}

type ShoeboxTypesMapResponse struct {
	Types map[string]uint8
}
