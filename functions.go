package main

import(
	"fmt"
	"os"
	"os/exec"
	"log"
	"github.com/boltdb/bolt"
	"strings"
	"errors"
)

func openFile(){
	cmd := exec.Command("vim","-o","buffer.md")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func show_lists(db *bolt.DB) error {
	if db == nil { // this means if DB is nil
		log.Fatal("show_lists was given a null database")
	}

	err := db.Update(func(tx *bolt.Tx) error{
		bucket := tx.Bucket(bucketName)
		if bucket == nil{
			log.Fatal("show_lists couldnt open the bucket")
		}

		get_list := bucket.Get([]byte("show_lists"))

		if get_list == nil{
			fmt.Println("No lists yet")
			return nil
		}
		temp_list := string(get_list)

		final_list := strings.Fields(temp_list)

		fmt.Println("\tAVAILABLE LISTS\n______________________________\n")
		for _ , val := range final_list {
			fmt.Print("- ")
			fmt.Println(val)
		}
		fmt.Println("______________________________\n")
		return nil
	})

	return err
}

func rc_content_manip(input, new_list string)(string){ // takes in content, puts new name on it and returns it
	var marker []int // used to keep place of where there ~
	marker_temp := []byte(input)
	for i , val := range input{
		if string(val) == "\n"{
			marker_temp[i] = '~' // right here is not registering it to input, only to value
			marker = append(marker,i + 3)
		}
	}
	input = string(marker_temp)
	temp_content := strings.Fields(input)
	temp_content[0] = new_list + "\n\n"
	input_temp := []byte(strings.Join(temp_content, " "))
	for _ , val := range marker{ // marker is not working well 
		if input_temp[val] != '~'{
			continue
		}else if val > len(input_temp){
			break // prevents seg fault
		}else{
			input_temp[val] = '\n'   // have to change this to a byte slice
		}
	}
	input = string(input_temp)
	return input
}// it's not replacing the ~ with \n. Also not getting rid of the field shit and taking out the initial test piece like i thouhgt

/*
	This function is for the manipulation of the characters for the readme for 
	backlog
*/

func backlog_content_manip(content []byte)(error,[]byte){
	input := string(content)
	if input == ""{
		return errors.New("Content is completely empty"),content
	}
	temp_input := strings.Split(input,string('\n'))
	for i:=1;i<len(temp_input)-1;i++{
		if len(temp_input[i]) <=1{
			continue
		}else{
			temp := " - " + temp_input[i]
			temp_input[i] = temp
		}
	}
	return nil,[]byte(strings.Join(temp_input, "\n"))
}

