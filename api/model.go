package api

type CalculatePacksRequest struct {
	ItemsCount int `json:"items_count"`
}

type CalculatePacksResponse struct {
	Packs []Pack `json:"packs"`
}

type Pack struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}

type SetPackSizes struct {
	Sizes []int `json:"sizes"`
}
