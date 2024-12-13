package flag_manager

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"password_generator/repository"

	genPass "password_generator/package/password_generator"
)

var showTableContent = flag.Bool("showme", false, "Flag to show table content")
var dropTableContent = flag.Bool("drop", false, "Flag to drop table content")
var deletePassword = flag.String("delete", "", "Flag to delete password by id")
var deleteLastPasswords = flag.String("delete-last", "", "Flag to delete multiple passwords from the end of the list")
var findService = flag.String("find-service", "", "Flag to find a group of passwords by service name")
var findPassword = flag.Int("find-password", 0, "Flag to find a password by its ID number")
var noSave = flag.Bool("no-save", false, "Flag to create a password without saving it in the database")

func ManageFlags(db *sql.DB) {
	// Initialising created flags
	flag.Parse()

	// Getting a service name from the console
	service := getServiceArg()

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		panic(fmt.Sprintf("Error parsing flags: %v", err))
	}

	switch {
	case *showTableContent:
		repository.PrintTableContent(db)
	case *dropTableContent:
		err := repository.ClearTable(db)
		if err != nil {
			panic(fmt.Sprintf("Error clearing table: %v", err))
		}
	case *deletePassword != "":
		err := repository.DeletePassword(db, *deletePassword)
		if err != nil {
			panic(fmt.Sprintf("Error deleting password: %v", err))
		}
	case *deleteLastPasswords != "":
		err := repository.DeleteLastPasswords(db, *deleteLastPasswords)
		if err != nil {
			panic(fmt.Sprintf("Error deleting password group: %v", err))
		}
	case *findService != "":
		err := repository.GetPasswordsByService(db, *findService)
		if err != nil {
			panic(fmt.Sprintf("Error getting passwords by service: %v", err))
		}
	case *findPassword > 0:
		err := repository.GetPasswordById(db, *findPassword)
		if err != nil {
			panic(fmt.Sprintf("Error getting a password by ID: %v", err))
		}
	case *noSave:
		fmt.Printf("New password: \n\n%s", genPass.GenerateRandomPasssword())
	default:
		randPass := genPass.GenerateRandomPasssword()
		repository.InsertIntoTable(db, randPass, service)
		fmt.Printf("New password: \n\n%s \n\nService: %s", randPass, service)
	}
}

func getServiceArg() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return "undefined"
}
