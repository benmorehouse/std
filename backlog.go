package main

import(
	"github.com/spf13/cobra"
	"github.com/boltdb/bolt"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"strings"
)

type backlog struct{
	created bool
	path string // the path to the github repo from 
}

var Backlog = &cobra.Command{
	Use: "backlog",
	Short:"updates the backlog",
	Run: func(cmd *cobra.Command, args []string){ // args is gonna be what we pass through 
		db, err := bolt.Open("mainDatabase.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close() // will close at end of function run

		err = db.Update(func(tx *bolt.Tx) error{
			bucket := tx.Bucket(bucketName)
			if bucket == nil {
				log.Fatal("Not able to open the bucket with all the lists")
			}

			content := bucket.Get([]byte("backlog")) // this will return what is in the list

			var mybacklog = backlog{
				created: true,
				path: "",
			}

			if content == nil{ // this means that backlog is not a key within the lists bucket yet
				for content == nil{
					fmt.Println("backlog not made yet... \nwhere do you wish to place backlog repo?\nnote you can change this at anytime within within markdown file")
					fmt.Println("Chosen Directory:")
					var chosen_directory string
					fmt.Scan(&chosen_directory)

					mybacklog.created = false
					mybacklog.path = string(chosen_directory)

					new_list_content := "# backlog \n\n\n\n\n"
					new_list_content += "# " + "Hit MQ when you are finished\n\n\ncurrent git location: " + chosen_directory
					err := bucket.Put([]byte("backlog"),[]byte(new_list_content))//creates the new bucket with nothing in it 
					if err != nil{
						log.Println("Unable to add new backlog bucket in backlog command")
					}
					content = bucket.Get([]byte("backlog"))
				}
			}

			path := strings.Fields(string(content))

			if mybacklog.path ==""{ // this means we already had a backlog file
				chosen_directory := path[len(path)-1] // this will give it the directory 
				mybacklog.path = string(chosen_directory)
			}

			file , err := os.Create("buffer.md")
			if err != nil{
				log.Println("Error opening file in backlog:",err)
			}

			_ , err = file.Write(content)

			if err != nil{
				log.Println("Error writing file in writelist: ", err)
			}

			openFile() // this will open the file and let the user input 

			content, err = ioutil.ReadFile("buffer.md")

			err = bucket.Put([]byte("backlog"),content)

			for err != nil{
				log.Println("error in backlog command putting content back into bucket",err)
			}
			// need to go and see if ioutil can fetch the file and if not we must create it 
			err = os.Chdir(mybacklog.path)
			if err != nil{
				log.Println("error: mybacklog.created = true and we couldnt switch into the new directory:",mybacklog.path)
				return err
			}

			err , writefile_content := backlog_content_manip(content)

			if err != nil{
				log.Println("Error with backlog_content_manip")
				return err
			}

			err = ioutil.WriteFile("README.md", writefile_content , 0644) // this will create the file if it doesnt already exist
			if err != nil{
				log.Println("error: mybacklog.create = true but couldnt write to backlog.md")
				return err
			}

			show_list_temp := bucket.Get([]byte("show_lists"))

			if show_list_temp == nil{
				// this means that show_list has yet to be created within the database
				return nil
			}else{
				database_names := strings.Fields(string(show_list_temp))
				exists := false
				for _ , val := range database_names{
					if val == "backlog"{
						exists = true
						break
					}
				}
				if exists{
					return nil
				}else{
					bucket.Put([]byte("show_lists"),[]byte("backlog" + "\n" +  string(show_list_temp)))
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal("error in write command on line 82:",err) // this will return if the database isnt op
		}
	},
}

