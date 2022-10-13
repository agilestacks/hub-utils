// Copyright (c) 2022 EPAM Systems, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	stateCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm <Stack ID>",
	Short: "Removes stack state from the project",
	Run:   rm,
	Args:  cobra.ExactValidArgs(1),
}

func rm(cmd *cobra.Command, args []string) {
	if Project == "" {
		altProjectSources()
	}
	fmt.Printf("Do you really want to remove the remote state of the stack [%s]? (type Y or Yes) ", args[0])
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	if text == "Y" || text == "Yes" {
		req, err := NewRequest()
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		resp, err := req.
			Delete(fmt.Sprintf("%s/%s", baseURL(), args[0]))
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		if resp.IsSuccess() {
			fmt.Printf("State \"%s\" is removed", args[0])
			return
		}
		if resp.StatusCode() == 404 {
			fmt.Printf("Error: State \"%s\" not found", args[0])
			return
		}
		fmt.Printf("Error: %s", resp.Status())
		return
	}
}
