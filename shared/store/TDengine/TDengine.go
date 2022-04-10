package TDengine

import (
	"database/sql"
	_ "github.com/i-Things/driver-go/v2/taosRestful"
)

type Td struct {
	Dsn string
	*sql.DB
}

func NewTDengine(DataSource string) (*Td, error) {
	taos, err := sql.Open("taosRestful", DataSource)
	if err != nil {
		return nil, err
	}
	return &Td{
		Dsn: DataSource,
		DB:  taos,
	}, nil
}
