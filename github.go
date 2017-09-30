package gostuff

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

func GetAllClosedCommits() {

	client := github.NewClient(nil)

	context := context.Background()
	opts := &github.SearchOptions{Sort: "created", Order: "asc"}

	page := 1
	query := "repo:jonpchin/gochess"

	results, _, err := client.Search.Issues(context, query, opts)
	if err != nil {
		fmt.Println("Unable to search issues", err)
		return
	}
	fmt.Println("Total issues are ", *results.Total)

	for _, value := range results.Issues {
		fmt.Println(value)
	}

	*results.Total -= 30
	for *results.Total > 0 {
		*results.Total -= 30
		page += 1
		query := "repo:jonpchin/gochess"

		results, _, err := client.Search.Issues(context, query, opts)
		if err != nil {
			fmt.Println("Unable to search issues", err)
			return
		}

		for _, value := range results.Issues {
			fmt.Println(value)
		}
	}
}
