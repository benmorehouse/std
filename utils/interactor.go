package utils

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Interactor explains what we expect from the user in order to
// generally use std
type Interactor interface {
	Edit(filename string) error
	Input() string
}

type iterm struct {
	rdr io.Reader
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

// StdInteractor will return our standard std user interactor
func StdInteractor() Interactor {
	return &iterm{rdr: os.Stdin}
}
