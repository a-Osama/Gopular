package main

import (
	"context"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

type Response struct {
	Search struct {
		Nodes []struct {
			NameWithOwner  string `json:"nameWithOwner"`
			StargazerCount int    `json:"stargazerCount"`
			CreatedAt      string `json:"createdAt"`
		} `json:"nodes"`
	} `json:"search"`
}

const query = `query PopularRepos {
	search(query: "stars:>1", type: REPOSITORY, first: 5) {
	nodes {
		... on Repository {
		  
		  nameWithOwner
		  stargazerCount
		  createdAt
		}
	  }
	}
  }`

func main() {
	client := graphql.NewClient("https://api.github.com/graphql")
	// make a request to GitHub API
	req := graphql.NewRequest(query)

	var GithubToken = os.Getenv("GITHUB_TOKEN")
	req.Header.Add("Authorization", "bearer "+GithubToken)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response

	var respData Response
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	log.Println(respData.Search.Nodes)
}
