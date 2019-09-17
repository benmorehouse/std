package main

import (
	"github.com/spf13/cobra"
	"github.com/boltdb/bolt"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"strings"
	"time"
)

//a function that will manage the backlog repo

var backlog struct{
	created bool
	path string // the path to the github repo from 
}


var updateBacklog{
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
				log.Fatal("Not able to open the bucket with all the lists:")
			}

			content := bucket.Get([]byte("backlog")) // this will return what is in the list

			if content == nil{ // this means that backlog is not a key within the lists bucket
				for content == nil{
					fmt.Println("backlog not made yet... \nwhere do you wish to place backlog repo\nnote you can change this at anytime within within markdown file")
					chosen_directory := ""
					fmt.Scan(&chosen_directory)

					var mybacklog := backlog{
						created: true,
						path: (string(chosen_directory))
					}

					new_list_content := "# backlog \n\n\n\n\n"
					new_list_content += "# " + "Hit MQ when you are finished\n\n\ncurrent git location:" + chosen_directory
					err := bucket.Put([]byte("backlog"),[]byte(new_list_content))//creates the new bucket with nothing in it 
					if err != nil{
						log.Println("Unable to add new backlog bucket in backlog command")
					}

					content = bucket.Get([]byte("backlog"))
				}
			}else{
				path := string.Fields(string(content))
				if path[len(path)-4] + path[len(path)-3 + path[len(path)-2] != "current git location:"{
					fmt.Println("Error: not at current git location in file")
					return nil
				}

				chosen_directory := path[len(path)-1] // this will give it the directory 

				var mybacklog := backlog{
					created: true,
					path: (string(chosen_directory))
				}
			}

			// loop through the file and put it into one big ass string. Then push that string to the bucket
			//first we need to write whats in the key to the file
			// then we let the user manipulate
			// then we .put it back in

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

			err = bucket.Put([]byte("backlog",content)

			for err != nil{
				log.Println("error in write command",err)
			}
			// need to go and see if ioutil can fetch the file and if not we must create it 
			err = ioutil.WriteFile(mybacklog.path,content, 0644) // this will create the file if it doesnt already exist
			if !mybacklog.created{
				// this means we need to create the directory that it needs to go in 
				//first we will mkdir in current directory, then move to new file path
				filename:=(strings.Fields(string(content)))
				filename = filename[len(filename)-1]
				var temp string
				for i:=len(filename)-1;i=0;i--{
					if filename[i] == '/'{
						break
					}else{
						temp+=filename[i]
					}
				}
				os.Mkdir(

			return nil
		})
		if err != nil {
			log.Fatal("error in write command on line 82:",err) // this will return if the database isnt open?
		}
	},
}
