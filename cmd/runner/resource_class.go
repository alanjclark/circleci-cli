package runner

import (
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/CircleCI-Public/circleci-cli/api/runner"
)

func newResourceClassCommand(o *runnerOpts, preRunE validator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource-class",
		Short: "Operate on runner resource-classes",
	}

	genToken := false
	createCmd := &cobra.Command{
		Use:     "create <resource-class> <description>",
		Short:   "Create a resource-class",
		Args:    cobra.ExactArgs(2),
		PreRunE: preRunE,
		RunE: func(_ *cobra.Command, args []string) error {
			rc, err := o.r.CreateResourceClass(args[0], args[1])
			if err != nil {
				return err
			}
			table := newResourceClassTable(cmd.OutOrStdout())
			defer table.Render()
			appendResourceClass(table, *rc)

			if !genToken {
				return nil
			}

			token, err := o.r.CreateToken(args[0], "default")
			if err != nil {
				return err
			}
			return generateConfig(*token, cmd.OutOrStdout())
		},
	}
	createCmd.PersistentFlags().BoolVar(&genToken, "generate-token", false,
		"Generate a default token")
	cmd.AddCommand(createCmd)

	cmd.AddCommand(&cobra.Command{
		Use:     "delete <resource-class>",
		Short:   "Delete a resource-class",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		PreRunE: preRunE,
		RunE: func(_ *cobra.Command, args []string) error {
			rc, err := o.r.GetResourceClassByName(args[0])
			if err != nil {
				return err
			}
			return o.r.DeleteResourceClass(rc.ID)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "list <namespace>",
		Short:   "List resource-classes for a namespace",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		PreRunE: preRunE,
		RunE: func(_ *cobra.Command, args []string) error {
			rcs, err := o.r.GetResourceClassesByNamespace(args[0])
			if err != nil {
				return err
			}

			table := newResourceClassTable(cmd.OutOrStdout())
			defer table.Render()
			for _, rc := range rcs {
				appendResourceClass(table, rc)
			}

			return nil
		},
	})

	return cmd
}

func newResourceClassTable(writer io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Resource Class", "Description"})
	return table
}

func appendResourceClass(table *tablewriter.Table, rc runner.ResourceClass) {
	table.Append([]string{rc.ResourceClass, rc.Description})
}
