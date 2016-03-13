package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Item struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
}

type ViewInventoryResponse struct {
	Items []Item `json:"items"`
}

func (api API) ViewInventoryHandler(c *gin.Context) {
	purchasedItemsTableName := "purchased_items"
	usedItemsTableName := "used_items"
	itemsTableName := "items"

	rows, err := api.DB.Query(
		fmt.Sprintf("select i.name, SUM(pu.quantity), pu.unit from (SELECT pi.item_id, COALESCE(SUM(pi.quantity), 0) as quantity, pi.unit from %s pi group by pi.item_id, pi.unit union select ui.item_id, COALESCE(-SUM(ui.quantity),0) as quantity, ui.unit from %s ui group by ui.item_id, ui.unit) pu join %s i on i.id = pu.item_id group by i.name, pu.unit",
			pq.QuoteIdentifier(purchasedItemsTableName),
			pq.QuoteIdentifier(usedItemsTableName),
			pq.QuoteIdentifier(itemsTableName),
		),
	)
	if err == sql.ErrNoRows {
		log.Println("no rows found")
		c.JSON(http.StatusOK, ViewInventoryResponse{})
		return
	}

	if err != nil {
		log.Println("Error in query: %s", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	result := ViewInventoryResponse{}
	log.Println(rows)
	defer rows.Close()
	for rows.Next() {
		item := Item{}
		rows.Scan(&item.Name, &item.Quantity, &item.Unit)
		result.Items = append(result.Items, item)
	}

	c.JSON(http.StatusOK, result)
}
