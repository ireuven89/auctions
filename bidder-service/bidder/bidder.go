package bidder

type Bidder struct {
	ID   string `json:"-"`
	Name string `json:"name"`
	Item string `json:"item"`
}

type BiddersRequest struct {
	Name string `json:"name"`
}
