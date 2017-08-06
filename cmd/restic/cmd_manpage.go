package main

import (
	"os"

	"github.com/restic/restic/internal/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var cmdManpage = &cobra.Command{
	Use:   "manpage [command]",
	Short: "generate manual pages",
	Long: `
The "manpage" command generates a manual page for a single command. It can also
be used to write all manual pages to a directory. If the output directory is
set and no command is specified, all manpages are written to the directory.
`,
	RunE: runManpage,
}

var manpageOpts = struct {
	OutputDir string
}{}

func init() {
	cmdRoot.AddCommand(cmdManpage)
	fs := cmdManpage.Flags()
	fs.StringVar(&manpageOpts.OutputDir, "output-dir", "", "write man pages to this `directory`")
}

func runManpage(cmd *cobra.Command, args []string) error {
	header := &doc.GenManHeader{
		Title:   "restic backup",
		Section: "1",
		Source:  "generated by `restic manpage`",
	}

	dir := manpageOpts.OutputDir
	if dir != "" {
		Verbosef("writing man pages to directory %v\n", dir)
		return doc.GenManTree(cmdRoot, header, dir)
	}

	switch {
	case len(args) == 0:
		return errors.Fatalf("no command given")
	case len(args) > 1:
		return errors.Fatalf("more than one command given: %v", args)
	}

	name := args[0]

	for _, cmd := range cmdRoot.Commands() {
		if cmd.Name() == name {
			return doc.GenMan(cmd, header, os.Stdout)
		}
	}

	return errors.Fatalf("command %q is not known", args)
}
