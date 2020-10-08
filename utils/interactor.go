package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/benmorehouse/std/configs"
	"github.com/benmorehouse/std/repo"
)

// Interactor explains what we expect from the user in order to
// generally use std
type Interactor interface {
	Edit(filename string) error
	Input() string
	RunLifeCycle(db repo.Repo, bucketName string, user Interactor, creatingNewBucket bool) error
}

type iterm struct {
	rdr io.Reader
}

// DefaultInteractor will return our standard std user interactor
func DefaultInteractor() Interactor {
	return &iterm{rdr: os.Stdin}
}

// DisplayBucketList will display the list of
// buckets in the root of the database
func DisplayBucketList(db repo.Repo) {
	fmt.Print("\tAVAILABLE LISTS\n______________________________\n\n")
	for _, val := range db.List() {
		fmt.Printf("- %s\n", val)
	}
	fmt.Print("______________________________\n\n")
}

// Edit will use vim to edit the given file
func (i *iterm) Edit(filename string) error {
	cmd := exec.Command("vim", "-o", filename)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Input will take in input
func (i *iterm) Input() string {
	reader := bufio.NewReader(i.rdr)
	input, _ := reader.ReadString('\n')
	return strings.TrimSuffix(input, "\n")
}

// RunLifeCycle will run the lifecycle of getting the contents of a bucket,
// putting them into a file, letting the user edit the file, then updating the
// bucket based on what the user has given us
func (i *iterm) RunLifeCycle(db repo.Repo, bucketName string, user Interactor, creatingNewBucket bool) error {
	var content string
	if !creatingNewBucket {
		content = db.Get(bucketName)
		for content == "" {
			fmt.Println("Not a valid list, please enter in an existing list.")
			bucketName = user.Input()
			content = db.Get(bucketName)
		}
	}

	file, err := os.Create(configs.STDConf.BufferMDFile)
	if err != nil {
		log.Println("Error opening file in writelist:", err)
	}

	if _, err = file.Write([]byte(content)); err != nil {
		log.Println("Error writing file in writelist: ", err)
	}

	if user.Edit(configs.STDConf.BufferMDFile); err != nil {
		return fmt.Errorf("user_edit_fail: %s", err.Error())
	}

	fileContent, err := ioutil.ReadFile(configs.STDConf.BufferMDFile)
	if err != nil {
		return fmt.Errorf("file_close_fail: %s", err.Error())
	}

	if err := db.Put(bucketName, string(fileContent)); err != nil {
		return fmt.Errorf("put_command_fail: %s", err.Error())
	}
	return nil
}
