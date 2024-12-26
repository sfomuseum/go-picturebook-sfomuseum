package response

type InstagramPostResponse struct {
	Post *InstagramPost
}

type InstagramPost struct {
	Caption         *InstagramPostCaption `json:"caption"`
	MediaId         string                `json:"media_id"`
	Path            string                `json:"path"`
	PerceptualHash  string                `json:"perceptual_hash"`
	Taken           int64                 `json:"taken"`
	TakenAt         string                `json:"taken_at"`
	WhosOnFirstId   int64                 `json:"wof:id"`
	WhosOnFirstRepo string                `json:"wof:repo"`
	SFOMuseumImage  string                `json:"sfomuseum:image"` // This should not be considered stable yet and may be replaced/removed
}

type InstagramPostCaption struct {
	Body     string   `json:"body"`
	Excerpt  string   `json:"excerpt"`
	HashTags []string `json:"hashtags"`
	Users    []string `json:"users"`
}
