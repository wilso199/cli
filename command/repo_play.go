package command

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/cli/cli/git"
	"github.com/cli/cli/tone"
	"github.com/spf13/cobra"
)

func repoPlay(cmd *cobra.Command, args []string) error {
	commits, err := git.LastXCommits(100)
	if err != nil {
		return fmt.Errorf("could get commits for repo: %w", err)
	}

	var notes []tone.Note
	for _, commit := range commits {
		fRegex, _ := regexp.Compile("\\s+(\\d+) file")
		filesChanges, _ := strconv.Atoi(fRegex.FindStringSubmatch(commit.Info)[1])

		fDeletions, _ := regexp.Compile("(\\d+) deletion")
		deletions := fDeletions.FindStringSubmatch(commit.Info)
		var x = 0
		if deletions != nil {
			x, _ = strconv.Atoi(deletions[1])
		}

		index := convertToInt(commit.Sha[0:1])
		length := filesChanges*50 + 50
		delay := x * 2

		freq := 220 * math.Pow(1.5, float64(index))
		output := fmt.Sprintf("%v : %v", commit.Sha[0:7], commit.Title)
		notes = append(notes, tone.Note{freq, length, delay, output})

	}
	tone.Play(notes)

	return nil
}

func convertToInt(hex string) int {
	value, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse commit sha: %w", err)
	}
	return int(value)
}

type key struct {
	Name string
	Freq float64
}
