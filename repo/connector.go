package repo

import (
	"github.com/benmorehouse/std/configs"
	"github.com/boltdb/bolt"

	"log"
)

// Connector will bridge a connection to the embedded
// database. It will report an error if conn fails
type Connector interface {
	Connect() (repo Repo, err error)
	Disconnect() (err error)
}

// DefaultConnector will return a connector that will
// bridge to the embedded database
func DefaultConnector() Connector {
	return &defaultConnector{}
}

// defaultConnector will implement the connector interface
type defaultConnector struct {
	db *bolt.DB
}

// Connect will use the default connector to wrap a connection
// to the database
func (d *defaultConnector) Connect() (repo Repo, err error) {
	d.db, err = bolt.Open(configs.STDConf.DatabasePath, 0744, nil)
	if err != nil {
		log.Println("Error opening database connection", err)
		return
	}
	return &defaultRepo{db: d.db}, nil
}

func (d *defaultConnector) Disconnect() error {
	return d.db.Close()
}
