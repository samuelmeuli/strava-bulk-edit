package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"strings"
	"syscall"
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

	// Args
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
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	email, password := promptCredentials()

	switch cmd {
	case title.FullCommand():
		fmt.Println("Update title")
		fmt.Println(*titleArg)
	case description.FullCommand():
		fmt.Println("Update description")
		fmt.Println(*descriptionArg)
	case commute.FullCommand():
		fmt.Println("Update commute")
		fmt.Println(*commuteArg)
	case sport.FullCommand():
		fmt.Println("Update sport")
		fmt.Println(*sportArg)
	case visibility.FullCommand():
		fmt.Println("Update visibility")
		fmt.Println(*visibilityArg)
	default:
		log.Fatal("No command specified")
	}

	fmt.Println(email)
	fmt.Println(password)
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
