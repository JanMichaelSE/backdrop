/*
Copyright © 2024 Jan Montalvo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"

	"github.com/janmichaelse/backdrop/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "backdrop",
	Short: "Backdrop is a command-line utility for setting, reverting, and organizing desktop wallpapers.",
	Long: `backdrop is a command-line utility for managing wallpapers on your desktop.
It allows you to set a new wallpaper, revert to a previous wallpaper, 
and specify the directory where your wallpaper images are stored.`,
	Version:      "2.0.0",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			cmd.Usage()
			return err
		}
		imageUrl, err := cmd.Flags().GetBool("url")
		if err != nil {
			cmd.Usage()
			return err
		}
		isSlideShow, err := cmd.Flags().GetBool("slideshow")
		if err != nil {
			cmd.Usage()
			return err
		}

		config := internal.NewConfig(path, imageUrl, isSlideShow)
		return internal.BackdropAction(os.Stdout, config, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.backdrop.yaml)")

	rootCmd.Flags().StringP("path", "p", "", "Set a custom path to find wallpaper images. If not provided, a default path will be used.")
	rootCmd.Flags().BoolP("slideshow", "s", false, "Will configure and set a custom slideshow of images you select with fzf.\nTo select multiple images hit 'Tab' on the images you desire to select, then hit 'Enter' to confirm.")
	rootCmd.Flags().BoolP("url", "u", false, `You will be prompted to provide an image url to be set as wallpaper. The image will be downloaded and previewed. 
    If confirmed, the image will be downloaded to the directory were all images are found (check "IMAGES" section). If image is NOT accepted by user, 
    the image gets deleted and previous wallpaper is set.`)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".backdrop" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".backdrop")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
