package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/machinebox/graphql"
	"github.com/spf13/cobra"
)

// popularCmd represents the popular command
var popularCmd = &cobra.Command{
	Use:   "popular",
	Short: "return most popular GitHub repositories",
	Long: `By default it fetches for the 10 most popular repositories on GitHub:

You could configure it to return public repos based on PROGRAMMING LANGUAGE, PUPLISH DATE, and how many RESULTS you want.`,
	Run: func(cmd *cobra.Command, args []string) {
		progLang, _ := cmd.Flags().GetString("p")
		date, _ := cmd.Flags().GetString("d")
		count, _ := cmd.Flags().GetUint("c")
		getPopularRepos(progLang, date, count)
	},
}

func init() {
	rootCmd.AddCommand(popularCmd)
	popularCmd.PersistentFlags().String("p", "Go", "Programming Language")
	popularCmd.PersistentFlags().String("d", "2014-01-01", "Date in format of yyyy-mm-dd")
	popularCmd.PersistentFlags().Uint("c", 10, "Count")

}

type Response struct {
	Search struct {
		Nodes []struct {
			NameWithOwner  string `json:"nameWithOwner"`
			StargazerCount int    `json:"stargazerCount"`
			CreatedAt      string `json:"createdAt"`
		} `json:"nodes"`
	} `json:"search"`
}

//GraphQL Query with two request variables.
var query = `query PopularRepos($first: Int = 10, $query: String!) {
	search(query: $query, type: REPOSITORY, first: $first) {
	  nodes {
		... on Repository {
		  nameWithOwner
		  stargazerCount
		  createdAt
		}
	  }
	}
  }`

//Initiate the client to the GraphQL API, and bind the query to the request.
//Add the request variable to the query using the values of the flags.
//Run it and capture the response.
//Format the Response.
//TODO make a function to format the response
func getPopularRepos(progLang, date string, count uint) {

	client := graphql.NewClient("https://api.github.com/graphql")

	// make a request to GitHub API
	req := graphql.NewRequest(query)
	req.Var("first", count)
	req.Var("query", fmt.Sprintf("language:%s stars:>1 created:>%s", progLang, date))
	//Get the github token from the enviroment variables
	var GithubToken = os.Getenv("GITHUB_TOKEN")
	if GithubToken == "" {
		log.Fatalln("GithubToken is required")
	}
	req.Header.Add("Authorization", "bearer "+GithubToken)
	// define a Context for the request
	ctx := context.Background()

	var respData Response
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	formatOutput(respData, count)

}

func formatOutput(response Response, count uint) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(writer, "Owner/Name\t Stars\t Publish Date")
	var i uint
	for i = 0; i < count; i++ {
		fmt.Fprintln(writer, response.Search.Nodes[i].NameWithOwner+"\t", fmt.Sprint(response.Search.Nodes[i].StargazerCount)+"\t", response.Search.Nodes[i].CreatedAt+"\t")
	}
	writer.Flush()

}
