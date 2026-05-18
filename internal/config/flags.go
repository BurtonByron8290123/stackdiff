package config

import "flag"

// BindFlags registers all CLI flags onto the provided FlagSet and returns
// a Config that will be populated when the FlagSet is parsed.
func BindFlags(fs *flag.FlagSet) *Config {
	c := Default()

	fs.StringVar(&c.BaseFile, "base", "", "Path to the base Terraform state file (required)")
	fs.StringVar(&c.TargetFile, "target", "", "Path to the target Terraform state file (required)")
	fs.StringVar(&c.Format, "format", c.Format, "Output format: text, markdown, json")
	fs.StringVar(&c.OutputFile, "output", "", "Write report to this file instead of stdout")
	fs.BoolVar(&c.Overwrite, "overwrite", c.Overwrite, "Overwrite output file if it already exists")
	fs.StringVar(&c.FilterName, "filter-name", "", "Only include resources whose address contains this prefix")
	fs.BoolVar(&c.NoColor, "no-color", c.NoColor, "Disable ANSI color codes in text output")

	return c
}
