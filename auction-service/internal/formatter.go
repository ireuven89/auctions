package internal

import "github.com/ireuven89/auctions/auction-service/auction"

func formatAuction(auction *auction.Auction) map[string]interface{} {

	return map[string]interface{}{
		"id":          auction.ID,
		"name":        auction.Name,
		"description": auction.Description,
		"user_id":     auction.UserId,
		"active":      auction.Active,
		"end_time":    auction.EndTime,
		"created_at":  auction.CreatedAt,
		"updated_at":  auction.UpdatedAt,
	}
}
