package notes

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/github"
)

func writeCommitsToJSONfile() {
	client := github.NewClient(nil)

	context := context.Background()
	var listOptions github.ListOptions
	listOptions.Page = 1

	opts := &github.SearchOptions{Sort: "created", Order: "asc", ListOptions: listOptions}

	query := "repo:jonpchin/gochess state:closed"

	results, _, err := client.Search.Issues(context, query, opts)
	if err != nil {
		fmt.Println("Unable to search issues", err)
		return
	}
	fmt.Println("Total issues are ", *results.Total)

	f, err := os.OpenFile("release_notes.json", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Can't open release notes file", err)
	}

	defer f.Close()

	for _, value := range results.Issues {
		releaseNotes, _ := json.MarshalIndent(value, "", "    ")
		_, err = f.Write([]byte(releaseNotes))
		if err != nil {
			fmt.Println("Couldn't write to release notes file", err)
		}
	}

	*results.Total -= 30
	for *results.Total > 0 {
		*results.Total -= 30
		listOptions.Page += 1
		opts = &github.SearchOptions{Sort: "created", Order: "asc", ListOptions: listOptions}

		results, _, err := client.Search.Issues(context, query, opts)
		if err != nil {
			fmt.Println("Unable to search issues", err)
			return
		}

		for _, value := range results.Issues {
			releaseNotes, _ := json.MarshalIndent(value, "", "    ")
			_, err = f.Write([]byte(releaseNotes))
			if err != nil {
				fmt.Println("Couldn't write to release notes file", err)
			}
		}
	}
}

func GetAllClosedCommits() {

	client := github.NewClient(nil)

	context := context.Background()
	var listOptions github.ListOptions
	listOptions.Page = 1
	opts := &github.SearchOptions{Sort: "created", Order: "asc", ListOptions: listOptions}

	query := "repo:jonpchin/gochess state:closed"

	results, _, err := client.Search.Issues(context, query, opts)
	if err != nil {
		fmt.Println("Unable to search issues", err)
		return
	}

	fmt.Println("Total issues are ", *results.Total)
	release_notes := "release_notes.txt"
	os.Remove(release_notes)

	notes, err := os.OpenFile(release_notes, os.O_CREATE|os.O_APPEND, 0644)
	defer notes.Close()
	if err != nil {
		fmt.Println("Can't open release notes text file", err)
	}

	t := time.Now()
	notes.WriteString(t.Format("2006-01-02") + "\n")

	for _, value := range results.Issues {
		notes.WriteString(*value.Title + "\n")
	}

	*results.Total -= 30
	for *results.Total > 0 {
		*results.Total -= 30
		listOptions.Page += 1
		opts = &github.SearchOptions{Sort: "created", Order: "asc", ListOptions: listOptions}

		results, _, err := client.Search.Issues(context, query, opts)
		if err != nil {
			fmt.Println("Unable to search issues", err)
			return
		}

		for _, value := range results.Issues {
			notes.WriteString(*value.Title + "\n")
		}
	}
}
