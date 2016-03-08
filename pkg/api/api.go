package api

import (
	"database/sql"
)

type API struct {
	DB *sql.DB
}
