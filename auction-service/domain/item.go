package domain

type Item struct {
	ID        string
	AuctionID string
	Name      string
}

type ItemPicture struct {
	ItemID       string
	AuctionID    string
	Name         string
	DownloadLink string
}
