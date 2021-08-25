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

//GraphQL Response Struct contain the repo name, star count, and puplish date
type Response struct {
	Search struct {
		Nodes []struct {
			NameWithOwner  string `json:"nameWithOwner"`
			StargazerCount int    `json:"stargazerCount"`
			CreatedAt      string `json:"createdAt"`
		} `json:"nodes"`
	} `json:"search"`
}

type Request struct {
	PrgLanguage string
	Date        string
	Count       uint
}

var popularCmd = &cobra.Command{
	Use:   "popular",
	Short: "return most popular GitHub repositories",
	Long: `By default it fetches for the 10 most popular repositories on GitHub:

You could configure it to return public repos based on PROGRAMMING LANGUAGE, PUPLISH DATE, and how many RESULTS you want.`,
	Run: func(cmd *cobra.Command, args []string) {
		var req Request
		req.PrgLanguage, _ = cmd.Flags().GetString("p")
		req.Date, _ = cmd.Flags().GetString("d")
		req.Count, _ = cmd.Flags().GetUint("c")
		getPopularRepos(req)
	},
}

func init() {
	rootCmd.AddCommand(popularCmd)
	popularCmd.PersistentFlags().String("p", "Go", "Programming Language")
	popularCmd.PersistentFlags().String("d", "2014-01-01", "Date in format of yyyy-mm-dd")
	popularCmd.PersistentFlags().Uint("c", 10, "Count")

}

//GraphQL Query with two request variables.
var query = `query PopularRepos($count: Int = 10, $qry: String!) {
	search(query: $qry, type: REPOSITORY, first: $count) {
	  nodes {
		... on Repository {
		  nameWithOwner
		  stargazerCount
		  createdAt
		}
	  }
	}
  }`

//return the Envoiroment Variable GITHUB_TOKEN
func getGitHubToken() (string, error) {
	var GithubToken = os.Getenv("GITHUB_TOKEN")
	if GithubToken == "" {
		return "", fmt.Errorf("please set your Github token")
	}
	return GithubToken, nil
}

//takes request, programming language, date, and counter
func reqFormating(req *graphql.Request, r Request) {
	req.Var("count", r.Count)
	req.Var("qry", fmt.Sprintf("language:%s stars:>1 created:>%s", r.PrgLanguage, r.Date))
}

//takes the response and format it as a nice table
func outputFormating(response Response, count uint) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(writer, "Owner/Name\t Stars\t Publish Date")
	var i uint
	for i = 0; i < count; i++ {
		fmt.Fprintln(writer, response.Search.Nodes[i].NameWithOwner+"\t",
			fmt.Sprint(response.Search.Nodes[i].StargazerCount)+"\t",
			response.Search.Nodes[i].CreatedAt+"\t")
	}
	writer.Flush()
}

func getPopularRepos(r Request) Response {
	var response Response
	client := graphql.NewClient("https://api.github.com/graphql")
	// make a request to GitHub API
	req := graphql.NewRequest(query)
	// Add the request variable to the query using the values of the flags.
	reqFormating(req, r)
	//Get the github token from the enviroment variables
	gt, err := getGitHubToken()
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "bearer "+gt)
	// define a Context for the request
	ctx := context.Background()
	if err := client.Run(ctx, req, &response); err != nil {
		log.Fatal(err)
	}
	outputFormating(response, r.Count)
	return response
}
