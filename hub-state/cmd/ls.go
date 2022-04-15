package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/xeonx/timeago"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

var stackFilter []string

func init() {
	lsCmd.Flags().StringSliceVar(&stackFilter, "filter", []string{}, "Filter by name, status or initiator. Example: --filter \"name=GKE,status=incomplete\"")
	stateCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all stacks within a project",
	Run:   ls,
}

func ls(cmd *cobra.Command, args []string) {
	base := BaseURL(StateAPILocation, Project)
	filterQuery := make(map[string]string, 0)
	for _, param := range stackFilter {
		vals := strings.Split(param, "=")
		if len(vals) == 2 {
			key := strings.ToLower(vals[0])
			if key == "initiator" {
				key = "latestOperation.initiator"
			}
			filterQuery[key] = vals[1]
		}
	}
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(filterQuery).
		Get(base)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	if Out == JsonO {
		var pretty bytes.Buffer
		err = json.Indent(&pretty, resp.Body(), "", "\t")
		if err != nil {
			fmt.Println("Nothing has been found")
			return
		}
		fmt.Println(pretty.String())
		return
	}
	var states []State
	json.Unmarshal(resp.Body(), &states)
	limit := 20
	if len(states) < limit {
		limit = len(states)
	}
	if len(states) > 0 {
		tableFmtStates(states[:limit])
		return
	}
	fmt.Println("Nothing has been found")
}

func tableFmtStates(states []State) {
	fmt.Printf("Listing Stacks in [%s] GCP project", Project)
	fmt.Println()
	fmt.Println()
	headerFmt := color.New().SprintfFunc()
	tbl := table.New("ID", "NAME", "STATUS", "UPDATED", "INITIATOR")
	tbl.WithHeaderFormatter(headerFmt)
	for _, widget := range states {
		tbl.AddRow(widget.ID, widget.Name, widget.Status, timeago.English.Format(widget.LatestOP.Timestamp), widget.LatestOP.Initiator)
	}
	tbl.Print()
}
