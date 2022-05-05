package cmd

import "fmt"

func BaseURL(StateAPILocation, Project string) string {
	return fmt.Sprintf("https://%s-%s.cloudfunctions.net/stacks", StateAPILocation, Project)
}
