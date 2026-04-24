/*
Copyright © 2023 Sirrend

*/

package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/enescakir/emoji"
	"github.com/sirrend/terrap-cli/internal/commons"
	"github.com/sirrend/terrap-cli/internal/state"
	"github.com/sirrend/terrap-cli/internal/utils"
	"github.com/sirrend/terrap-cli/internal/utils/terraform"
	"github.com/spf13/cobra"
)

// terraformInit
/*
@brief:
	terraformInit performs the Terraform init command on the given folder
@params:
	dir - the folder to initialize
*/
func terraformInit(dir string) {
	_, err := os.Stat(path.Join(dir, ".terrap.json"))

	if err != nil {
		_, _ = commons.YELLOW.Println(emoji.Rocket, "Initializing directory...")
		mainWorkspace.ExecPath, mainWorkspace.IsTempProvider,
			mainWorkspace.TerraformVersion, err = terraform.TerraformInit(dir) // initiate new terraform tool in context

		if err != nil {
			fmt.Println()
			terraform.TerraformErrorPrettyPrint(err)
			os.Exit(1)

		}

		_, _ = commons.YELLOW.Print(emoji.Toolbox, " Looking for providers...")
		terraform.FindTfProviders(dir, &mainWorkspace) //find all providers and assign to mainWorkspace
		_, _ = commons.GREEN.Println(" Done!")

		_, _ = commons.YELLOW.Print(emoji.WavingHand, " Saving workspace...")
		saveInitData() //Save to configuration file
		_, _ = commons.GREEN.Println(" Done!")

	} else {
		// Note: changed exit message to be more actionable and mention the -u flag clearly
		_, _ = commons.YELLOW.Println("Folder already initialized.")
		_, _ = commons.YELLOW.Println("Run `terrap init -u` to upgrade your context, or `terrap init -d <dir>` for a different directory.")
		os.Exit(0)

	}
}

/*
@brief: saveInitData saves the configuration file of the initialized folder
*/
func saveInitData() {
	err := state.Save(path.Join(mainWorkspace.Location, ".terrap.json"), mainWorkspace)
	if err != nil {
		_, _ = commons.RED.Println("Terrap failed saving the current workspace.")
		os.Exit(1)
	}
}

// the init command declaration
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize directory",
	Long:  `Initialize directory for sirrend to have all needed files`,

	Run: //check which flags are set and run the appropriate init
	func(cmd *cobra.Command, args []string) {
		if cmd.Flag("upgrade").Changed {
			var directory string

			if cmd.Flag("directory").Changed {
				directory, _ = filepath.Abs(cmd.Flag("directory").Value.String())
			} else {
				directory, _ = os.Getwd()
			}

			deleteInitData(directory)
			mainWorkspace.Location = directory
			terraformInit(directory)

			fmt.Println()
			_, _ = commons.SIRREND.Println(emoji.BeerMug, "Terrap directory upgraded!")

		} else if cmd.Flag("directory").Changed {
			if utils.IsDir(cmd.Flag("directory").Value.String()) {
				directory, _ := filepath.Abs(cmd.Flag("directory").Value.String())
				mainWorkspace.Location = directory
				terraformInit(directory)
				_, _ = commons.SIRREND.Println("\nTerrap Initialized Successfully!")

			} else {
				_, _ = commons.YELLOW.Println("The given path is not a directory")
				os.Exit(0)
			}

		} else {
			location, err := os.Getwd() // get current directory
			if err != nil