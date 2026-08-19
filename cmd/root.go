package cmd

import (
	"fmt"
	"os"

	kubecontext "github.com/chris-cmsoft/gotool-kubecontext-picker"
	"github.com/chris-cmsoft/knein/internal/k9s"
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
		limit: kubecontext.DefaultLimit,
	}

	cmd := &cobra.Command{
		Use:           "knein",
		Short:         "Open k9s for a Kubernetes context",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.limit < 1 {
				return kubecontext.ErrInvalidLimit
			}

			contexts, err := kubecontext.Load(opts.kubeconfig)
			if err != nil {
				return err
			}

			selected, err := kubecontext.Select(contexts, opts.limit)
			if err != nil {
				return err
			}
			if selected == "" {
				return nil
			}

			return k9s.Open(selected)
		},
	}

	cmd.Flags().StringVar(&opts.kubeconfig, "kubeconfig", "", "Path to kubeconfig file")
	cmd.Flags().IntVar(&opts.limit, "limit", opts.limit, "Maximum contexts to show")

	return cmd
}
