package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// Note: the "binding" tags here don't actually work, probably because
// they're nested
type PurchasedItem struct {
	ItemId   int64
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
	err := c.BindJSON(&submitPurchaseRequest)
	if err != nil {
		log.Println(err)
	}

	tx, err := api.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %s", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	purchasedItemsTableName := "purchased_items"
	itemsTableName := "items"
	for _, purchasedItem := range submitPurchaseRequest.PurchasedItems {
		err := tx.QueryRow(
			fmt.Sprintf("SELECT id FROM %s WHERE name = $1", pq.QuoteIdentifier(itemsTableName)),
			purchasedItem.Name,
		).Scan(&purchasedItem.ItemId)

		if err == sql.ErrNoRows {
			err = tx.QueryRow(
				fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", pq.QuoteIdentifier(itemsTableName)),
				purchasedItem.Name,
			).Scan(&purchasedItem.ItemId)

			if err != nil {
				log.Printf("Error with insert: %s", err)
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

		} else if err != nil {
			log.Printf("Error querying for item: %s", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		_, err = tx.Exec(
			fmt.Sprintf("INSERT INTO %s (item_id, quantity, unit, cost) VALUES ($1, $2, $3, $4)", pq.QuoteIdentifier(purchasedItemsTableName)),
			purchasedItem.ItemId,
			purchasedItem.Quantity,
			purchasedItem.Unit,
			purchasedItem.Cost,
		)
		if err != nil {
			log.Printf("Error with insert: %s", err)
			c.AbortWithError(http.StatusInternalServerError, err)
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
		log.Printf("Error committing transaction: %s", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, CreateBasketResponse{
		OK: true,
	})
}
