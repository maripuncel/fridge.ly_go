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
type UsedItem struct {
	ItemId   int64
	Name     string `json:"name" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
	Unit     string `json:"unit" binding:"required"`
}

type CreateUsageRequest struct {
	UsedItems []UsedItem `json:"used_items" binding:"required"`
}

type CreateUsageResponse struct {
	OK bool `json:"ok"`
}

func (api API) CreateUsageHandler(c *gin.Context) {
	var createUsageRequest CreateUsageRequest
	err := c.BindJSON(&createUsageRequest)
	if err != nil {
		log.Println(err)
	}

	tx, err := api.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %s", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	usedItemsTableName := "used_items"
	itemsTableName := "items"
	for _, usedItem := range createUsageRequest.UsedItems {
		err := tx.QueryRow(
			fmt.Sprintf("SELECT id FROM %s WHERE name = $1", pq.QuoteIdentifier(itemsTableName)),
			usedItem.Name,
		).Scan(&usedItem.ItemId)

		if err == sql.ErrNoRows {
			err = tx.QueryRow(
				fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", pq.QuoteIdentifier(itemsTableName)),
				usedItem.Name,
			).Scan(&usedItem.ItemId)

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
			fmt.Sprintf("INSERT INTO %s (item_id, quantity, unit) VALUES ($1, $2, $3)", pq.QuoteIdentifier(usedItemsTableName)),
			usedItem.ItemId,
			usedItem.Quantity,
			usedItem.Unit,
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
