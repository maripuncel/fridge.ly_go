package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Item struct {
	Name     string   `json:"name"`
	Quantity int      `json:"quantity"`
	Unit     string   `json:"unit"`
}

type ViewInventoryResponse struct {
	Items []Item `json:"purchased_items"`
}

func (api API) ViewInventoryHandler(c *gin.Context) {
	purchasedItemsTableName := "purchased_items"
	rows, err := api.DB.Query(
		fmt.Sprintf("SELECT name, quantity, unit FROM %s", pq.QuoteIdentifier(purchasedItemsTableName)),
		)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	result := ViewInventoryResponse{}
	for rows.Next() {
		item := Item{}
		rows.Scan(&item.Name, &item.Quantity, &item.Unit)
		result.Items = append(result.Items, item)
	}

	if !c.IsAborted() {
		c.JSON(http.StatusOK, result)
	} else{
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
