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

// ListConnector will return a connector that will
// bridge to the embedded database
func ListConnector() Connector {
	return &listConnector{}
}

// PasswordConnector will create a connector for the password db storage
func PasswordConnector() Connector {
	return &passwordConnector{}
}

// listConnector will create a database connector for our lists
type listConnector struct {
	db *bolt.DB
}

// passwordConnector will create a connection to our separate password database
type passwordConnector struct {
	db *bolt.DB
}

// Connect will use the default connector to wrap a connection
// to the database
func (d *listConnector) Connect() (repo Repo, err error) {
	d.db, err = bolt.Open(configs.STDConf.ListDatabasePath, 0744, nil)
	if err != nil {
		log.Fatal("Error opening database connection", err)
	}
	return &listRepo{db: d.db}, nil
}

func (d *listConnector) Disconnect() error {
	return d.db.Close()
}

// Connect will use the default connector to wrap a connection
// to the database
func (d *passwordConnector) Connect() (repo Repo, err error) {
	d.db, err = bolt.Open(configs.STDConf.PasswordDatabasePath, 0744, nil)
	if err != nil {
		log.Fatal("Error opening database connection", err)
	}
	return &passwordRepo{db: d.db}, nil
}

func (d *passwordConnector) Disconnect() error {
	return d.db.Close()
}
