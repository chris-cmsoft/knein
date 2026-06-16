package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/chris-cmsoft/knein/internal/kubecontexts"
	"github.com/chris-cmsoft/knein/internal/picker"
	"github.com/spf13/cobra"
)

type rootOptions struct {
	kubeconfig string
	limit      int
}

// Execute runs the root command.
func Execute() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	opts := rootOptions{
		limit: 9,
	}

	cmd := &cobra.Command{
		Use:   "knein",
		Short: "Open k9s for a Kubernetes context",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.limit < 1 {
				return errors.New("limit must be greater than 0")
			}

			contexts, err := kubecontexts.Load(opts.kubeconfig)
			if err != nil {
				return err
			}
			if len(contexts) == 0 {
				return errors.New("no Kubernetes contexts found")
			}

			selected, err := picker.SelectContext(contexts, opts.limit)
			if err != nil {
				return err
			}
			if selected == "" {
				return nil
			}

			return picker.OpenK9s(selected)
		},
	}

	cmd.Flags().StringVar(&opts.kubeconfig, "kubeconfig", "", "Path to kubeconfig file")
	cmd.Flags().IntVar(&opts.limit, "limit", opts.limit, "Maximum contexts to show")

	return cmd
}
