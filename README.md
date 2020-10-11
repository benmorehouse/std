# STD

[![codecov](https://codecov.io/gh/benmorehouse/std/branch/develop/graph/badge.svg)](https://codecov.io/gh/benmorehouse/std)
[![Build Status](https://travis-ci.com/benmorehouse/std.svg?branch=develop)](https://travis-ci.com/benmorehouse/std)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens)
![version](https://img.shields.io/badge/version-1.0.0-blue)


## A simple task manager

Welcome to STD (or as I like to call, shit to do)! This is a simple task manager that can be used to 
organize tasks from work, school, groceries, or anything else you want to keep a list of. 

STD is useful not only to keep work held in your work environment, but also because it is your own way 
of organizing what you want to have on your computer!

# Using STD 

If you just wish to use the project as a tool, [then run the binary!](https://github.com/benmorehouse/std/releases/tag/v1.0)

After installing, ./std will show you the commands that are available for you to use. 

	backlog - A place for ideas that can also be placed in it's own repo for your own github!
	create - create a list
	delete - delete list
	help - display the help menu
	open - open the designated list
	rename - rename a given list to something else
	welcome - welcome message as well as version that is currently in use

Examples
	
	./std open 
	./std delete worklist
	./std backlog



# Contributing

If you wish to fork this project and work on it more, then by all means!

If you do not have homebrew, go here: https://brew.sh

Once you have homebrew, if you do not have go installed then entire
	
	brew install go

into your desired directory.
Then install the following two package:

	Boltdb: go get github.com/boltdb/bolt/...
        Cobra: go get -u github.com/spf13/cobra/cobra 

Once this is installed, pick the directory you wish to store std, then 
clone this repository!

Finally, do the following command within the std file you downloaded: 
        go build 

Then you are all set to run STD!
		
	./std


# Coming soon
	
Functions in development include
	share - share a certain list with any contact on your macbook via imessage
	directory - create a list that is a directory of more lists and so on and so forth
	Search Box - As you type, a search box with autofill functions will fill in for you


Contributions are welcome at all times! 

# Thanks and welcome to STD - A better way to manage your shit.
	
