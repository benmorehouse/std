package repo

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// Repo mocks basic responsibilities of the database that holds
// the functionalities of the list.
type Repo interface {
	List() []string
	Put(bucketName, newValue string) error
	Create(bucketName string) error
	Remove(key string) error
	Get(key string) string
}

// defaultRepo will implement the repo interface
type defaultRepo struct {
	db *bolt.DB
}

// List will list what is in the repo
func (d *defaultRepo) List() (list []string) {
	d.db.Update(func(tx *bolt.Tx) error {
		list = []string{}
		cursor := tx.Cursor()
		for name, _ := cursor.First(); name != nil; name, _ = cursor.Next() {
			list = append(list, string(name))
		}
		return nil
	})
	return
}

// Put will put whatever is in the repo
func (d *defaultRepo) Put(bucketName, newValue string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket_not_exists")
		}
		return bucket.Put([]byte(bucketName), []byte(newValue))
	})
}

// Create will create a list and drop it in our repo
func (d *defaultRepo) Create(bucketName string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(bucketName),
			[]byte(fmt.Sprintf("# %s \n\n\n\n\n # Close Vim window when you are finished", bucketName)))
	})
}

// Remove will remove the bucket
func (d *defaultRepo) Remove(key string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(key))
	})
}

// Get will go and get the value at the given key
func (d *defaultRepo) Get(key string) (value string) {
	d.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(key))
		if bucket == nil {
			return fmt.Errorf("bucket_not_found")
		}
		fmt.Println(key)
		value = string(bucket.Get([]byte(key)))
		return nil
	})
	return
}
