package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Masterminds/vcs"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/cmd/helm/require"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/helmpath"
)

func main() {
	cmd := NewInitCmd(os.Stdout)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type initOptions struct {
	starter    string
	directory  string
	name       string
	starterDir string
}

func NewInitCmd(out io.Writer) *cobra.Command {
	o := &initOptions{}

	cmd := &cobra.Command{
		Use:          "init NAME",
		Short:        "Initialize new chart with the given name",
		Args:         require.ExactArgs(1),
		SilenceUsage: true,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				// Allow file completion when completing the argument for the name
				// which could be a path
				return nil, cobra.ShellCompDirectiveDefault
			}
			// No more completions, so disable file completion
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			o.name = args[0]
			o.starterDir = helmpath.DataPath("starters")
			return o.run(out)
		},
	}

	cmd.Flags().StringVarP(&o.starter, "starter", "p", "", "the name or url to Helm starter scaffold")
	cmd.Flags().StringVarP(&o.directory, "directory", "d", "scaffold", "Directory within repo that holds the Chart.yaml file")
	return cmd
}

func (o *initOptions) run(out io.Writer) error {
	fmt.Fprintf(out, "Creating %s\n", o.name)

	chartname := filepath.Base(o.name)
	cfile := &chart.Metadata{
		Name:        chartname,
		Description: "A Helm chart for Kubernetes",
		Type:        "application",
		Version:     "0.1.0",
		AppVersion:  "0.1.0",
		APIVersion:  chart.APIVersionV2,
	}

	if o.starter != "" {
		if isLocalReference(o.starter) {
			// Create from the starter
			lstarter := filepath.Join(o.starterDir, o.starter)
			// If path is absolute, we don't want to prefix it with helm starters folder
			if filepath.IsAbs(o.starter) {
				lstarter = o.starter
			}
			return chartutil.CreateFrom(cfile, filepath.Dir(o.name), lstarter)
		}

		repo, err := vcs.NewRepo(o.starter, o.starterDir)
		if err != nil {
			return err
		}

		if _, err := os.Stat(repo.LocalPath()); os.IsNotExist(err) {
			repo.Get()
		} else {
			repo.Update()
		}

		lstarter := filepath.Join(repo.LocalPath(), o.directory)

		if _, err := os.Stat(lstarter); os.IsNotExist(err) {
			return fmt.Errorf("%s directory is missing", o.directory)
		}

		return chartutil.CreateFrom(cfile, filepath.Dir(o.name), lstarter)
	}

	chartutil.Stderr = out
	_, err := chartutil.Create(chartname, filepath.Dir(o.name))

	return err
}

func isLocalReference(source string) bool {
	_, err := os.Stat(source)
	return err == nil
}
