package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/list"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/pkg/convert"
	"github.com/bengarrett/retrotxtgo/pkg/str"
	"github.com/spf13/cobra"
)

type Lists int

const (
	Codepages Lists = iota
	Examples
	Table
	Tables
)

func (l Lists) Command() *cobra.Command {
	switch l {
	case Codepages:
		return ListCodepages()
	case Examples:
		return ListExamples()
	case Table:
		return ListTable()
	case Tables:
		return ListTables()
	}
	return nil
}

func ListCodepages() *cobra.Command {
	return &cobra.Command{
		Use:     "codepages",
		Aliases: []string{"c", "cp"},
		Short: fmt.Sprintf("List the legacy codepages that %s can convert to UTF-8",
			meta.Name),
		Long: fmt.Sprintf("List the available legacy codepages that %s can convert to UTF-8.",
			meta.Name),
		GroupID: "codepages",
		RunE: func(cmd *cobra.Command, args []string) error {
			b := convert.List()
			fmt.Fprint(cmd.OutOrStdout(), b)
			return nil
		},
	}
}

func ListExamples() *cobra.Command {
	return &cobra.Command{
		Use:     "examples",
		Aliases: []string{"e", "samples"},
		GroupID: "exaCmds",
		Short: fmt.Sprintf("List builtin tester text files available for use with the %s and %s commands",
			str.Example("info"), str.Example("view")),
		Long: fmt.Sprintf("List builtin tester text art and documents available for use with the %s and %s commands.",
			str.Example("info"), str.Example("view")),
		Example: fmt.Sprint(example.ListExamples),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := list.Examples()
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), b)
			return nil
		},
	}
}

func ListTable() *cobra.Command {
	return &cobra.Command{
		Use:     "table [codepage names or aliases]",
		Aliases: []string{"t"},
		Short:   "Display one or more codepage tables showing all the characters in use",
		Long:    "Display one or more codepage tables showing all the characters in use.",
		Example: fmt.Sprint(example.ListTable),
		GroupID: "tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := flag.Help(cmd, args...); err != nil {
				return err
			}
			s, err := list.Table(args...)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), fmt.Sprintln(s))
			return nil
		},
	}
}

func ListTables() *cobra.Command {
	return &cobra.Command{
		Use:     "tables",
		Short:   "Display the characters of every codepage table in use",
		Long:    "Display the characters of every codepage table in use.",
		GroupID: "tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := list.Tables()
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), fmt.Sprintln(s))
			return nil
		},
	}
}
