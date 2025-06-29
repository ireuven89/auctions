package bidder

type Bidder struct {
	ID   string `json:"-"`
	Name string `json:"name"`
	Bid  string `json:"bid"`
}

type BiddersRequest struct {
	Name string `json:"name"`
}
