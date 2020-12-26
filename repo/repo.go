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
	Remove(key string) error
	Get(key string) string
}

// defaultRepo will implement the repo interface
type listRepo struct {
	db *bolt.DB
}

type passwordRepo struct {
	db *bolt.DB
}

// List will list what is in the repo
func (d *listRepo) List() (list []string) {
	d.db.View(func(tx *bolt.Tx) error {
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
func (d *listRepo) Put(bucketName, newValue string) error {
	return d.db.Update(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			if bucket, err = tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
				return fmt.Errorf("unable_to_create_entry: %s", err.Error())
			}
		}

		if err = bucket.Put([]byte(bucketName), []byte(newValue)); err != nil {
			return fmt.Errorf("unable_to_put_entry: %s", err.Error())
		}
		return nil
	})
}

// Remove will remove the bucket
func (d *listRepo) Remove(key string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(key))
	})
}

// Get will go and get the value at the given key
func (d *listRepo) Get(key string) (value string) {
	d.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(key))
		if bucket == nil {
			return fmt.Errorf("bucket_not_found")
		}
		value = string(bucket.Get([]byte(key)))
		return nil
	})
	return
}

// ############################################################
// ###################### password repo #######################

const passwordBucketKey = "password"

// List will list what is in the repo
func (d *passwordRepo) List() []string {
	list := []string{}
	d.db.View(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(passwordBucketKey))
		if bucket == nil {
			if bucket, err = tx.CreateBucketIfNotExists([]byte(passwordBucketKey)); err != nil {
				return fmt.Errorf("unable_to_list_password_entry: %s", err.Error())
			}
		}

		cursor := bucket.Cursor()
		for name, _ := cursor.First(); name != nil; name, _ = cursor.Next() {
			list = append(list, string(name))
		}
		return nil
	})
	return list
}

// Put will put whatever is in the repo
func (d *passwordRepo) Put(key, newValue string) error {
	return d.db.Update(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(passwordBucketKey))
		if bucket == nil {
			if bucket, err = tx.CreateBucketIfNotExists([]byte(passwordBucketKey)); err != nil {
				return fmt.Errorf("unable_to_create_password_entry: %s", err.Error())
			}
		}

		if err = bucket.Put([]byte(key), []byte(newValue)); err != nil {
			return fmt.Errorf("unable_to_put_entry: %s", err.Error())
		}
		return nil
	})
}

// Remove will remove the bucket
func (d *passwordRepo) Remove(key string) error {
	return d.db.Update(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(passwordBucketKey))
		if bucket == nil {
			if bucket, err = tx.CreateBucketIfNotExists([]byte(passwordBucketKey)); err != nil {
				return fmt.Errorf("unable_to_create_password_entry: %s", err.Error())
			}
		}

		if err = bucket.Delete([]byte(key)); err != nil {
			return fmt.Errorf("unable_to_delete_password: %s", err.Error())
		}
		return nil
	})
}

// Get will go and get the value at the given key
func (d *passwordRepo) Get(key string) (value string) {
	d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(passwordBucketKey))
		if bucket == nil {
			return fmt.Errorf("password_bucket_not_found")
		}
		value = string(bucket.Get([]byte(key)))
		return nil
	})
	return
}
