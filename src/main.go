package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/alecthomas/kingpin.v2"
)

const dateRegex = `^\d{4}-\d{2}-\d{2}$`

var (
	app = kingpin.New("strava-bulk-edit", "Edit multiple Strava activities at once").Version("0.0.1").Author("Samuel Meuli")

	// Flags
	all  = app.Flag("all", "Update all Strava activities").Bool()
	from = app.Flag("from", "Update all Strava activities after the specified date. Date format: YYYY-MM-DD").String()
	to   = app.Flag("to", "Update all Strava activities before the specified date. Date format: YYYY-MM-DD").String()

	// Commands and args
	title          = app.Command("title", "Update the activities' titles")
	titleArg       = title.Arg("title", "New title for the activities").Required().String()
	description    = app.Command("description", "Update the activities' descriptions")
	descriptionArg = description.Arg("description", "New description for the activities").Required().String()
	commute        = app.Command("commute", "Set the activities as commute/non-commute activities")
	commuteArg     = commute.Arg("commute", "true or false indicating whether activities should be marked commute").Required().Bool()
	sport          = app.Command("sport", "Update the activities' sports type")
	sportArg       = sport.Arg("sport", "New sports type for the activities").Required().Enum(optionsSport...)
	visibility     = app.Command("visibility", "Update the activities' visibility")
	visibilityArg  = visibility.Arg("visibility", "New visibility for the activities").Required().Enum(optionsVisibility...)
)

func main() {
	// Parse and validate args
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	validateFlags()

	// Ask for Strava credentials
	email, password := promptCredentials()

	// Update the specified attribute
	switch cmd {
	case title.FullCommand():
		update(email, password, Activity{Title: *titleArg})
	case description.FullCommand():
		update(email, password, Activity{Description: *descriptionArg})
	case commute.FullCommand():
		update(email, password, Activity{Commute: *commuteArg})
	case sport.FullCommand():
		update(email, password, Activity{Sport: *sportArg})
	case visibility.FullCommand():
		update(email, password, Activity{Visibility: *visibilityArg})
	default:
		log.Fatal("No command specified")
	}
}

func promptCredentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	// Prompt for email
	fmt.Print("Enter Strava email: ")
	email, _ := reader.ReadString('\n')

	// Prompt for password
	fmt.Print("Enter Strava password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal("Please enter a valid password.")
	}
	password := string(bytePassword)
	fmt.Println()

	return strings.TrimSpace(email), strings.TrimSpace(password)
}

func validateFlags() {
	// Validate that only compatible flags are used
	allDefined := *all == true
	fromDefined := *from != ""
	toDefined := *to != ""
	if !allDefined && !fromDefined && !toDefined {
		log.Fatal("Please use a flag (--all, --from, and/or --to) to specify the activities to update")
	}
	if allDefined && (fromDefined || toDefined) {
		log.Fatal("The --all flag cannot be used together with the --from and --to flags")
	}

	// Validate date format
	if fromDefined {
		match, _ := regexp.MatchString(dateRegex, *from)
		if !match {
			log.Fatal("Invalid --from date. Please use the format YYYY-MM-DD")
		}
	}
	if toDefined {
		match, _ := regexp.MatchString(dateRegex, *to)
		if !match {
			log.Fatal("Invalid --to date. Please use the format YYYY-MM-DD")
		}
	}
}
