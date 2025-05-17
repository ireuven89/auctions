package bidder

type Bidder struct {
	ID   string `json:"-"`
	Name string `json:"name"`
}

type BiddersRequest struct {
	Name string `json:"name"`
}
