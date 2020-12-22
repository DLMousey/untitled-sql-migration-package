package main

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"io/ioutil"
	"log"
	"untitled-agnostic-sql-migration-package/database"
	"untitled-agnostic-sql-migration-package/structs"
)

func main() {
	/**
		Take the config.hcl file at the root of the project, parse it and convert the objects
		in that file into structs that we can work with
	 */
	var config structs.Projects
	contents, err := ioutil.ReadFile("config.hcl")
	if err != nil {
		log.Fatal(err)
	}

	err = hcl.Decode(&config, string(contents))
	if err != nil {
		log.Fatal(err)
	}

	// Make a channel for goroutine fun times of type DatabaseState
	channel := make(chan structs.DatabaseState)

	/**
		Iterate over every project that was found in the HCL config, pass it's config data
		to the CheckInitialisationStatus method and wait for a result back on the channel
	 */
	for _, project := range config.Projects {
		go database.CheckInitialisationStatus(&project, channel)

		x := <-channel
		if x.Initialised {
			fmt.Println("Project " + project.Name + " is initialised")
		} else {
			database.InitialiseDatabase(&project)
			fmt.Println("Project " + project.Name + " is not initialised")
		}

		database.Migrate(&project)
	}

}
