package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// Note: the "binding" tags here don't actually work, probably because
// they're nested
type PurchasedItem struct {
	Name     string   `json:"name" binding:"required"`
	Quantity int      `json:"quantity" binding:"required"`
	Unit     string   `json:"unit" binding:"required"`
	Cost     *float64 `json:"cost" binding:"required"`
}

type CreateBasketRequest struct {
	PurchasedItems []PurchasedItem `json:"purchased_items" binding:"required"`
}

type CreateBasketResponse struct {
	OK bool `json:"ok"`
}

func (api API) CreateBasketHandler(c *gin.Context) {
	var submitPurchaseRequest CreateBasketRequest
	if c.BindJSON(&submitPurchaseRequest) == nil {
		tx, err := api.DB.Begin()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		purchasedItemsTableName := "purchased_items"
		for _, purchasedItem := range submitPurchaseRequest.PurchasedItems {
			_, err := tx.Exec(
				fmt.Sprintf("INSERT INTO %s (name, quantity, unit, cost) VALUES ($1, $2, $3, $4)", pq.QuoteIdentifier(purchasedItemsTableName)),
				purchasedItem.Name,
				purchasedItem.Quantity,
				purchasedItem.Unit,
				purchasedItem.Cost,
			)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				err := tx.Rollback()
				if err != nil {
					c.Error(err)
				}
				return
			}
		}

		err = tx.Commit()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var resp CreateBasketResponse
		resp.OK = true
		c.JSON(http.StatusOK, resp)
	}
}
