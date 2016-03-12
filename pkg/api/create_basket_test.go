package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func apiWithLocalDB(t *testing.T) API {
	db, err := sql.Open("postgres", "postgres://localhost/fridgely_test?sslmode=disable")
	if err != nil {
		t.Fatalf("Unable to connect to test database: %s", err)
	}

	_, err = db.Exec("truncate purchased_items;")
	if err != nil {
		t.Fatalf("Unable to truncate purchased_items table in test database: %s", err)
	}

	return API{
		DB: db,
	}
}

func TestCreateBasketHandler(t *testing.T) {
	cost := 1.28
	testItem := PurchasedItem{
		Name:     "eggs",
		Quantity: 12,
		Unit:     "count",
		Cost:     &cost,
	}
	request := CreateBasketRequest{
		PurchasedItems: []PurchasedItem{testItem},
	}
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Unable to marshal request as json: %s", err)
	}
	body := bytes.NewReader(bodyBytes)

	context, recorder, _ := gin.CreateTestContext()
	context.Request = &http.Request{
		Body: ioutil.NopCloser(body),
	}
	if err != nil {
		t.Fatalf("Unable to create test http request: %s", err)
	}

	api := apiWithLocalDB(t)

	api.CreateBasketHandler(context)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected %d response, got %d", http.StatusOK, recorder.Code)
	}

	var resp CreateBasketResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("Unable to unmarshal response: %s", err)
	}

	if !resp.OK {
		t.Errorf("Expected response OK to be true, was false")
	}

	if len(context.Errors) != 0 {
		t.Errorf("Got some errors: %s", context.Errors)
	}

	rows, err := api.DB.Query("select name, quantity, unit, cost from purchased_items;")
	if err != nil {
		t.Fatalf("Unable to query database to make sure insert worked: %s", err)
	}

	var count int
	item := &PurchasedItem{}
	for rows.Next() {
		rows.Scan(&item.Name, &item.Quantity, &item.Unit, &item.Cost)
		if item.Name != testItem.Name {
			t.Errorf("Item name didn't match: wanted %s got %s", testItem.Name, item.Name)
		}

		if item.Quantity != testItem.Quantity {
			t.Errorf("Item name didn't match: wanted %d got %d", testItem.Quantity, item.Quantity)
		}

		if item.Unit != testItem.Unit {
			t.Errorf("Item name didn't match: wanted %s got %s", testItem.Unit, item.Unit)
		}

		if *item.Cost != *testItem.Cost {
			t.Errorf("Item name didn't match: wanted %s got %s", testItem.Cost, item.Cost)
		}
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 row in purchased_items table, there were %d", count)
	}
}
