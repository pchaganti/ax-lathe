package cmd

import (
	"fmt"

	"github.com/devenjarvis/lathe/internal/skills"
	"github.com/spf13/cobra"
)

var skillsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the skills bundled into this binary",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, err := skills.All()
		if err != nil {
			return err
		}
		out := cmd.OutOrStdout()
		for _, s := range all {
			_, _ = fmt.Fprintf(out, "%-14s %s\n", s.Slug, s.Description)
		}
		return nil
	},
}

func init() {
	skillsCmd.AddCommand(skillsListCmd)
}
