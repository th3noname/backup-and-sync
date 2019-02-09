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

type job interface {
	run() error
	name() string
	continueOnError() bool
	logFields() log.Fields
}

// Config contains the information on all actions that should be performed
type Config struct {
	Copy []Copy `mapstructure:"copy"`
	Sync []Sync `mapstructure:"sync"`
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
		err := r.callJob(&v)

		if err != nil {
			return err
		}
	}

	for _, v := range r.config.Sync {
		err := r.callJob(&v)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Rclone) callJob(j job) error {
	err := j.run()

	if err != nil {
		if j.continueOnError() {
			log.WithError(err).WithFields(j.logFields()).Warnf("run %s job failed. Continuing...", j.name())
			return nil
		}

		return errors.Wrapf(err, "run %s job failed", j.name())
	}

	return nil
}

func execute(arguments []string) error {
	log.WithField("arguments", arguments).Info("Executing rclone command")

	w := log.StandardLogger().Writer()

	command := exec.Command("rclone", arguments...)
	command.Stdout = os.Stdout
	command.Stderr = w
	err := command.Run()

	if err == nil {
		log.Info("rclone exited with return code 0")
	}

	return errors.Wrap(err, "rclone exec failed")
}

// Copy represents a single rclone copy job
type Copy struct {
	Source          string `mapstructure:"source"`
	Destination     string `mapstructure:"destination"`
	BwLimit         string `mapstructure:"bw-limit"`
	ContinueOnError bool   `mapstructure:"continue-on-error"`
}

func (c *Copy) name() string {
	return "copy"
}

func (c *Copy) continueOnError() bool {
	return c.ContinueOnError
}

func (c *Copy) logFields() log.Fields {
	return log.Fields{
		"source":      c.Source,
		"destination": c.Destination,
	}
}

func (c *Copy) run() error {
	log.WithFields(c.logFields()).Infof("start run rclone %s", c.name())

	args := []string{c.name()}
	args = append(args, c.Source)
	args = append(args, c.Destination)
	args = append(args, "--stats-log-level", "NOTICE")
	args = append(args, "--stats", "1m")

	if c.BwLimit != "" {
		args = append(args, "--bwlimit", c.BwLimit)
	}

	err := execute(args)

	log.Infof("end run  rclone %s", c.name())
	return errors.Wrap(err, "execute failed")
}

// Sync represents a single rclone copy job
type Sync struct {
	Source          string `mapstructure:"source"`
	Destination     string `mapstructure:"destination"`
	BwLimit         string `mapstructure:"bw-limit"`
	ContinueOnError bool   `mapstructure:"continue-on-error"`
}

func (s *Sync) name() string {
	return "sync"
}

func (s *Sync) continueOnError() bool {
	return s.ContinueOnError
}

func (s *Sync) logFields() log.Fields {
	return log.Fields{
		"source":      s.Source,
		"destination": s.Destination,
	}
}

func (s *Sync) run() error {
	log.WithFields(s.logFields()).Infof("start run rclone %s", s.name())

	args := []string{s.name()}
	args = append(args, s.Source)
	args = append(args, s.Destination)
	args = append(args, "--stats-log-level", "NOTICE")
	args = append(args, "--stats", "1m")

	if s.BwLimit != "" {
		args = append(args, "--bwlimit", s.BwLimit)
	}

	err := execute(args)

	log.Infof("end run rclone %s", s.name())
	return errors.Wrap(err, "execute failed")
}
