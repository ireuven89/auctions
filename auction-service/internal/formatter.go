package internal

import (
	"encoding/json"
	"github.com/ireuven89/auctions/auction-service/domain"
)

func formatAuction(auction *domain.Auction) map[string]interface{} {
	var regions string
	json.Unmarshal(auction.Regions, &regions)

	return map[string]interface{}{
		"id":           auction.ID,
		"description":  auction.Description,
		"regions":      regions,
		"starting_bid": auction.InitialOffer,
		"currentOffer": auction.CurrentBid,
		"status":       auction.Status.String(),
		"created_at":   auction.CreatedAt,
		"updated_at":   auction.UpdatedAt,
	}
}

func formatItem(auction *domain.Item) map[string]interface{} {

	return map[string]interface{}{
		"id":          auction.ID,
		"description": auction.Description,
		"created_at":  auction.CreatedAt,
		"updated_at":  auction.UpdatedAt,
		"auction_id":  auction.AuctionID,
	}
}
