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
				log.Fatal("Not able to open the bucket with all the lists:")
			}

			content := bucket.Get([]byte("backlog")) // this will return what is in the list
// we assume initially that backlog has been created
			path := strings.Fields(string(content))
			if path[len(path)-4] + path[len(path)-3] + path[len(path)-2] != "current git location:"{
				fmt.Println("Error: not at current git location in file")
				return nil
			}

			chosen_directory := path[len(path)-1] // this will give it the directory 

			var mybacklog = backlog{
				created: true,
				path: (string(chosen_directory)),
			}

			if content == nil{ // this means that backlog is not a key within the lists bucket yet
				for content == nil{
					fmt.Println("backlog not made yet... \nwhere do you wish to place backlog repo\nnote you can change this at anytime within within markdown file")
					fmt.Println("Chosen Directory:")
					chosen_directory := ""
					fmt.Scan(&chosen_directory)

					mybacklog.create = false
					mybacklog.path = string(chose_directory)

					new_list_content := "# backlog \n\n\n\n\n"
					new_list_content += "# " + "Hit MQ when you are finished\n\n\ncurrent git location: " + chosen_directory
					err := bucket.Put([]byte("backlog"),[]byte(new_list_content))//creates the new bucket with nothing in it 
					if err != nil{
						log.Println("Unable to add new backlog bucket in backlog command")
					}

					content = bucket.Get([]byte("backlog"))
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

			err = bucket.Put([]byte("backlog"),content)

			for err != nil{
				log.Println("error in backlog command putting content back into bucket",err)
			}
			// need to go and see if ioutil can fetch the file and if not we must create it 
			if !mybacklog.created{
				// this means we need to create the directory that it needs to go in 
				filename_field := strings.Fields(string(content))
				filename := filename_field[len(filename_field)-1] // this have some sort of string like /users/benmorehouse
				var temp string
				for i:=0;i<len(filename);i++{
					temp+=string(filename[i])
				}
				// at this point temp is the filepath
				mybacklog.path = temp
				temp = ""
				err = os.Chdir(mybacklog.path)
				if err != nil{
					log.Println("error: mybacklog.created = false and we couldnt switch into the new directky:",mybacklog.path)
					return err
				}

				err = ioutil.Writefile("backlog.md", content, 0644)
				if err != nil{
					log.Println("error: couldnt do writefile")
					return err
				}

			}else{
				err = os.Chdir(mybacklog.path)
				if err != nil{
					log.Println("error: mybacklog.created = true and we couldnt switch into the new directory:",mybacklog.path)
					return err
				}

				err = ioutil.WriteFile(mybacklog.path+"/backlog.md", content , 0644) // this will create the file if it doesnt already exist
				if err != nil{
					log.Println("error: mybacklog.create = true but couldnt write to backlog.md")
					return err
				}
			}
			// at this point we have a markdown file. 
			return nil
		})
		if err != nil {
			log.Fatal("error in write command on line 82:",err) // this will return if the database isnt open?
		}
	},
}


