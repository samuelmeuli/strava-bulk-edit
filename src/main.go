package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/alecthomas/kingpin.v2"
)

var valuesSport = []string{"AlpineSki", "BackcountrySki", "Canoeing", "Crossfit", "EBikeRide", "Elliptical", "Handcycle", "Hike", "IceSkate", "InlineSkate", "Kayaking", "Kitesurf", "NordicSki", "Ride", "RockClimbing", "RollerSki", "Rowing", "Run", "Snowboard", "Snowshoe", "StairStepper", "StandUpPaddling", "Surfing", "Swim", "VirtualRide", "VirtualRun", "Walk", "WeightTraining", "Wheelchair", "Windsurf", "Workout", "Yoga"}
var valuesVisibility = []string{"only_me", "followers_only", "everyone"}

var (
	app = kingpin.New("strava-bulk-edit", "Edit multiple Strava activities at once").Version("0.0.1").Author("Samuel Meuli")

	// Flags
	all  = app.Flag("all", "Update all Strava activities").Bool()
	from = app.Flag("from", "Update all Strava activities after the specified date").String()
	to   = app.Flag("to", "Update all Strava activities before the specified date").String()
	// TODO make exactly one flag required, make all and from/to flags mutually exclusive

	// Commands and args
	title          = app.Command("title", "Update the activities' titles")
	titleArg       = title.Arg("title", "New title for the activities").Required().String()
	description    = app.Command("description", "Update the activities' descriptions")
	descriptionArg = description.Arg("description", "New description for the activities").Required().String()
	commute        = app.Command("commute", "Set the activities as commute/non-commute activities")
	commuteArg     = commute.Arg("commute", "true or false indicating whether activities should be marked commute").Required().Bool()
	sport          = app.Command("sport", "Update the activities' sports type")
	sportArg       = sport.Arg("sport", "New sports type for the activities").Required().Enum(valuesSport...)
	visibility     = app.Command("visibility", "Update the activities' visibility")
	visibilityArg  = visibility.Arg("visibility", "New visibility for the activities").Required().Enum(valuesVisibility...)
)

func main() {
	// Parse and validate args
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

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
