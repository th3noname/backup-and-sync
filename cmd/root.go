// Copyright © 2019 Jan Arens
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/th3noname/backup-and-sync/rclone"
	"github.com/th3noname/backup-and-sync/restic"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "backup-and-sync",
	Short: "Backup directories using restic and sync folders using rclone",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.IsSet("restic") {
			var resticConf *restic.Config

			err := viper.UnmarshalKey("restic", &resticConf)
			if err != nil {
				log.WithError(err).Error("Unmarshal restic configuration failed")
				return
			}

			r := restic.New(resticConf)
			err = r.Run()
			if err != nil {
				log.WithError(err).Error("restic execution failed")
				return
			}
		}

		if viper.IsSet("rclone") {
			var rcloneConf *rclone.Config

			err := viper.UnmarshalKey("rclone", &rcloneConf)
			if err != nil {
				log.WithError(err).Error("Unmarshal rclone configuration failed")
				return
			}

			r := rclone.New(rcloneConf)
			err = r.Run()
			if err != nil {
				log.WithError(err).Error("rclone execution failed")
				return
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Config file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./backup.config)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.Info("Start reading config file")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./")
		viper.SetConfigName("backup.config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Error("Reading Config file failed")
		os.Exit(1)
		return
	}

	log.Info("Using config file: ", viper.ConfigFileUsed())
}
