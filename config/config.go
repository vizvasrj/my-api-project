package config

import "database/sql"

type Conf struct {
	DB *sql.DB
}
