package internal

import "github.com/ireuven89/auctions/auction-service/auction"

func formatAuction(auction *auction.Auction) map[string]interface{} {

	return map[string]interface{}{
		"id":        auction.ID,
		"name":      auction.Name,
		"bidder_id": auction.BidderId,
	}
}
