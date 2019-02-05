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

package rclone

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Config contains the information on all actions that should be performed
type Config struct {
	Copy []Copy `mapstructure:"copy"`
}

// Copy represents a single rclone copy job
type Copy struct {
	Source          string `mapstructure:"source"`
	Destination     string `mapstructure:"destination"`
	BwLimit         string `mapstructure:"bw-limit"`
	ContinueOnError bool   `mapstructure:"continue-on-error"`
}

// Rclone is a CLI wrapper
type Rclone struct {
	config *Config
}

// New creates a Rclone wrapper instance for the provided config
func New(conf *Config) Rclone {
	return Rclone{config: conf}
}

// Run the configured rclone jobs
func (r *Rclone) Run() error {
	for _, v := range r.config.Copy {
		err := r.runCopy(v)

		if err != nil {
			if v.ContinueOnError {
				log.WithError(err).WithFields(log.Fields{
					"source":      v.Source,
					"destination": v.Destination,
				}).Warn("copy job failed. Continuing...")
			} else {
				return errors.Wrap(err, "runCopy failed")
			}
		}
	}

	return nil
}

func (r *Rclone) runCopy(c Copy) error {
	args := []string{"copy"}
	args = append(args, c.Source)
	args = append(args, c.Destination)

	if c.BwLimit != "" {
		args = append(args, "--bwlimit", c.BwLimit)
	}

	err := r.execute(args)
	return errors.Wrap(err, "execute failed")
}

func (r *Rclone) execute(arguments []string) error {
	log.WithField("arguments", arguments).Info("Executing rclone command")

	command := exec.Command("rclone", arguments...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()

	return errors.Wrap(err, "rclone exec failed")
}
