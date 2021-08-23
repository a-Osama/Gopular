package cmd_test

import (
	"os"
	"testing"
)

const succeed = "\u2713"
const failed = "\u2717"

func TestEnviromentVariable(t *testing.T) {
	got := os.Getenv("GITHUB_TOKEN")
	if len(got) == 0 {
		t.Errorf("\t%s\tGitHub Token is empty", failed)
	}
}

func TestPopularCommand(t *testing.T) {
}
