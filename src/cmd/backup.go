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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/th3noname/backup-and-sync/src/rclone"
	"github.com/th3noname/backup-and-sync/src/restic"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "run a backup",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()

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

func init() {
	rootCmd.AddCommand(backupCmd)
}
