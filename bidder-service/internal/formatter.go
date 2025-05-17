package internal

import "github.com/ireuven89/auctions/bidder-service/bidder"

func formatBidder(bidder bidder.Bidder) map[string]interface{} {

	return map[string]interface{}{
		"id":   bidder.ID,
		"name": bidder.Name,
	}
}
