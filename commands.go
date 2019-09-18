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

var bucketName = []byte("Lists")

var writeList = &cobra.Command{ // appends to the end of the bucket
// store some data that is written by the user
	Use: "open",
	Short:"Open the current list",
	Run: func(cmd *cobra.Command, args []string){ // args is gonna be what we pass through 
		// open tmp, let user input, then read line for line and add into bucket
		db, err := bolt.Open("mainDatabase.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close() // will close at end of function run

		var temp string
		if len(args) != 1{
			show_lists(db)
			fmt.Scan(&temp)
		}else{
			temp = args[0]
		}
		temp = strings.ToLower(strings.TrimSpace(temp))
		chosen_list_key := []byte(temp)
		err = db.Update(func(tx *bolt.Tx) error{
			bucket := tx.Bucket(bucketName)
			if bucket == nil {
				log.Fatal("Not able to open the bucket with all the lists:")
			}

			content := bucket.Get([]byte(chosen_list_key)) // this will return what is in the list

			if content == nil{ // this means that chose_list_key is not a key within the lists bucket
				for content == nil{
					fmt.Println("Not a valid list, please enter in an existing list")
					fmt.Scan(&temp)
					temp = strings.ToLower(strings.TrimSpace(temp))
					chosen_list_key = []byte(temp)
					content = bucket.Get([]byte(chosen_list_key))
				}
			}

			// loop through the file and put it into one big ass string. Then push that string to the bucket
			//first we need to write whats in the key to the file
			// then we let the user manipulate
			// then we .put it back in
			file , err := os.Create("buffer.md")
			if err != nil{
				log.Println("Error opening file in writelist:",err)
			}
			_ , err = file.Write(content)

			if err != nil{
				log.Println("Error writing file in writelist: ", err)
			}

			openFile() // this will open the file and let the user input 

			content, err = ioutil.ReadFile("buffer.md")

			err = bucket.Put(chosen_list_key,content)

			for err != nil{
				log.Println("error in write command",err)
			}
			return nil
		})
		if err != nil {
			log.Fatal("error in write command on line 82:",err) // this will return if the database isnt open?
		}
	},
}

var createList = &cobra.Command{
	Use: "create",
	Short: "create a list",
	Example: "./std create work",
	Run: func(cmd *cobra.Command, args []string){ // args is gonna be what we pass through 
		db, err := bolt.Open("mainDatabase.db", 0600, nil)
		if err != nil{
			log.Println("Error opening database at createlist command:",err)
		}
		defer db.Close()

		var temp string
		if len(args) != 1{ // this means that they didnt enter in any sort of list to add into the bucket
			fmt.Println("What is your new list called")
			fmt.Scan(&temp)
		}else{
			temp = args[0]
		}

		chosen_list_key := []byte(temp) // this is the new list name within the bucket lists

		err = db.Update(func(tx *bolt.Tx) error{ // error happening here 
			bucket , err := tx.CreateBucketIfNotExists(bucketName) // this is going into the database and accessing one bucket

			if err != nil {
				log.Println("blank bucket name or too long of a name")
			}

			checkKey := bucket.Get(chosen_list_key) //  if the list already exists then it will come up 
			// this will return nil if this doesnt exist which is what we want 
			for checkKey != nil{ // this means that they entered something that already exists
				fmt.Println("list already exists:")
				fmt.Scan(&temp)
				checkKey = bucket.Get([]byte(temp))
				chosen_list_key =[]byte(temp)
			} // at this point now we know that checkKey exists
			new_list_content := "# " + string(chosen_list_key) + "\n\n\n\n\n"
			new_list_content += "# " + "Hit MQ when you are finished"

			err = bucket.Put(chosen_list_key,[]byte(new_list_content))//creates the new bucket with nothing in it 

			if err != nil{
				log.Println("Unable to add new list in bucket in write command")
			}

			show_list_temp := bucket.Get([]byte("show_lists"))

			if show_list_temp == nil{
				// this means that show_list has yet to be created within the database
				bucket.Put([]byte("show_lists"),[]byte(chosen_list_key))
			}else{
				bucket.Put([]byte("show_lists"),[]byte(string(chosen_list_key) + "\n" +  string(show_list_temp)))
			}

			return nil
		})

		if err != nil {
			log.Fatal("error in create command:",err) // this will 
		}
	},
}

var deleteList = &cobra.Command{
	Use: "delete",
	Short: "Delete the list from the database",
	Example: "./std delete work",
	Run: func(cmd *cobra.Command, args []string){
		db, err := bolt.Open("mainDatabase.db", 0600, nil)
		if err != nil{
			log.Println("Error opening database at deletelist command:",err)
		}
		defer db.Close()

		var temp string
		if len(args) != 1{ // this means that they didnt enter in any sort of list to add into the bucket
			show_lists(db)
			fmt.Println("Which list do you want deleted?")
			fmt.Scan(&temp)
		}else{
			temp = args[0]
		}
		if temp == "backlog"{
			log.Fatal("You cannot delete the backlog")
		}
		chosen_list_key := []byte(temp) // this is the new list name within the bucket lists

		err = db.Update(func(tx *bolt.Tx) error{ // how can you run the show_lists function within db.Update?
			bucket , err := tx.CreateBucketIfNotExists(bucketName) // this is going into the database and accessing one bucket

			if err != nil {
				log.Println("blank bucket name or too long of a name")
			}

			checkKey := bucket.Get(chosen_list_key)

			for checkKey == nil{
				fmt.Println("list doesnt exist:")
				fmt.Scan(&temp)
				checkKey = bucket.Get([]byte(temp))
				chosen_list_key =[]byte(temp)
			}// chosen_list_key is what we are gonna get rid of

			err = bucket.Delete(chosen_list_key)

			if err != nil{
				log.Fatal("deleteList could not delete the chosen key")
			}

			show_list_temp := bucket.Get([]byte("show_lists"))

			if show_list_temp == nil{
				log.Fatal("show_list broken at delete list command")
			}

			temp_field := strings.Fields(string(show_list_temp))

			for i , val := range temp_field{
				if val == string(chosen_list_key){
					temp_field[i] = ""
				}
			}

			err = bucket.Put([]byte("show_lists"),[]byte(strings.Join(temp_field," ")))

			if err != nil{
				log.Fatal("could not renter show_lists in delete command")
			}

			return nil
		})

		if err != nil{
			log.Fatal("issue with end of delete command")
		}
	},
}

var renameList = &cobra.Command{
	Use:"rename",
	Short:"rename a list",
	Example:"./std rename or ./std rename current_list new_list",
	Run: func(cmd *cobra.Command, args []string){
		db, err := bolt.Open("mainDatabase.db", 0600, nil)
		if err != nil{
			log.Println("Error opening database at renamelist command:",err)
		}
		defer db.Close()
		show_lists(db)

		err = db.Update(func(tx *bolt.Tx) error{
			var chosen_list_temp string
			var new_list_temp string
			bucket := tx.Bucket(bucketName)
			if bucket == nil{
				log.Fatal("rename_list couldnt open the bucket")
			}

			if len(args) == 0{ // simply means they entered nothing
				fmt.Println("which list do you want renamed?")
				fmt.Scan(&chosen_list_temp)
			}else if len(args) == 1{ // means they entered in a list but not the new name of the list
				chosen_list_temp = (args[0])
			}else if len(args) == 2{
				chosen_list_temp = (args[0])
				new_list_temp = (args[0])
			}else{
				fmt.Println("To rename a list you must add in the current list followed by new list name")
				time.Sleep(time.Second * 1)
			}


			if chosen_list_temp == ""{
				fmt.Println("Which list do you want renamed?")
				fmt.Scan(&chosen_list_temp)
			}

			content := bucket.Get([]byte(chosen_list_temp)) // this will return what is in the list

			if content == nil{ // this means that chose_list_key is not a key within the lists bucket
				for content == nil{
					fmt.Println("Not a valid list, please enter in an existing list")
					var temp string
					fmt.Scan(&temp)
					temp = strings.ToLower(strings.TrimSpace(temp))
					chosen_list_temp = temp
					content = bucket.Get([]byte(chosen_list_temp))
				}
			}

			if new_list_temp == ""{
				fmt.Println("What should the new name be")
				fmt.Scan(&new_list_temp)
			}
			var new_list string
			new_list = strings.TrimSpace(strings.ToLower(new_list_temp))

			err := bucket.Delete([]byte(chosen_list_temp))
			if err != nil{
				log.Fatal("Couldnt delete list in rename function")
			}
// Need to come here, create a system that takes in content, adds marks wherever there are \n, then does fields and joins
			content = []byte(rc_content_manip(string(content), string(new_list))) // takes in content, puts new name and returns
			err = bucket.Put([]byte(new_list),[]byte(content))

			if err != nil{
				log.Fatal("Couldnt push new_list in rename function")
			}

			// now need to edit the show_lists key in the bucket

			show_list_temp := bucket.Get([]byte("show_lists"))

			if show_list_temp == nil{
				log.Fatal("show_list broken at rename_list command")
			}

			temp_field := strings.Fields(string(show_list_temp))

			for i , val := range temp_field{
				if val == string(chosen_list_temp){
					temp_field[i] = new_list
				}
			}

			err = bucket.Put([]byte("show_lists"),[]byte(strings.Join(temp_field," ")))

			if err != nil{
				log.Fatal("could not renter show_lists in delete command")
			}

			return nil
		})

		if err != nil{
			log.Fatal("Error at the end of the rename function")
		}
	},
}


var welcome = &cobra.Command{ // appends to the end of the bucket
// store some data that is written by the user
	Use: "welcome",
	Short:"welcome message",
	Run: func(cmd *cobra.Command, args []string){ // args is gonna be what we pass through 
		// open tmp, let user input, then read line for line and add into bucket
		fmt.Println("\n\n\t\t\t\t\t\t\t\tHELLO! WELCOME TO STD (SHIT TO DO)\n\n\t\t\t\t\tTHIS IS A CLI TASK MANAGER THAT CAN MULTITASK AS A NOTEPAD AND iMESSAGE HUB \n\n\t\t\t\t\t\t\t\t\tENJOY!")
	time.Sleep(3 * time.Second)
	},
}
