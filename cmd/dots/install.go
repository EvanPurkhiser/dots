package main

import (
	"github.com/spf13/cobra"

	"go.evanpurkhiser.com/dots/installer"
	"go.evanpurkhiser.com/dots/output"
	"go.evanpurkhiser.com/dots/resolver"
)

var installCmd = cobra.Command{
	Use:   "install [filter...]",
	Short: "Install and compile dotfiles from sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		forceReInstall, _ := cmd.Flags().GetBool("reinstall")
		verbose, _ := cmd.Flags().GetBool("verbose")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		info, _ := cmd.Flags().GetBool("info")

		dotfiles := resolver.ResolveDotfiles(*sourceConfig, *sourceLockfile).Filter(args)
		prepared := installer.PrepareDotfiles(dotfiles, *sourceConfig)

		installConfig := installer.InstallConfig{
			SourceConfig:   sourceConfig,
			SourceLockfile: sourceLockfile,
			ForceReinstall: forceReInstall,
		}

		installLogger := output.New(output.Config{
			SourceConfig:    *sourceConfig,
			InstallConfig:   installConfig,
			PreparedInstall: prepared,
			IsVerbose:       verbose,
			IsInfo:          info || verbose || dryRun,
		})

		if dryRun {
			installLogger.DryrunInstall()
			return nil
		}

		installed := installer.InstallDotfiles(prepared, installConfig)
		installer.RunInstallScripts(prepared, installConfig)
		installer.FinalizeInstall(installed, installConfig)

		// TODO Needs some error handling clenaup

		return nil
	},
	Args: cobra.ArbitraryArgs,
}

func init() {
	flags := installCmd.Flags()
	flags.SortFlags = false

	flags.BoolP("reinstall", "r", false, "forces execution of all installation scripts")
	flags.BoolP("info", "i", false, "prints install operation details")
	flags.BoolP("verbose", "v", false, "prints debug data, implies info")
	flags.BoolP("dry-run", "n", false, "do not mutate any dotfiles, implies info")
}
