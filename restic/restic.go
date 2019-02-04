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

package restic

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Config contains the information on all actions that should be performed
type Config struct {
	Repositoies []Repository `mapstructure:"repositories"`
	Backup      []Backup     `mapstructure:"backups"`
	Forget      []Forget     `mapstructure:"forget"`
}

// Repository stores information on a restic repository
type Repository struct {
	Repository string `mapstructure:"repository"`
	Path       string `mapstructure:"path"`
	Password   string `mapstructure:"password"`
}

// Backup represents a single restic backup job
type Backup struct {
	Backup          string   `mapstructure:"backup"`
	Repository      string   `mapstructure:"repository"`
	Source          string   `mapstructure:"source"`
	Exclude         []string `mapstructure:"exclude"`
	ContinueOnError bool     `mapstructure:"continue-on-error"`
}

// Forget represents a single restic forget job
type Forget struct {
	Repository      string   `mapstructure:"repository"`
	Prune           bool     `mapstructure:"prune"`
	KeepLast        int      `mapstructure:"keep-last"`
	KeepHourly      int      `mapstructure:"keep-hourly"`
	KeepDaily       int      `mapstructure:"keep-daily"`
	KeepWeekly      int      `mapstructure:"keep-weekly"`
	KeepMonthly     int      `mapstructure:"keep-monthly"`
	KeepYearly      int      `mapstructure:"keep-yearly"`
	KeepTag         []string `mapstructure:"keep-tag"`
	Tag             []string `mapstructure:"tag"`
	Hostname        string   `mapstructure:"hostname"`
	ContinueOnError bool     `mapstructure:"continue-on-error"`
}

// Restic is a CLI wrapper
type Restic struct {
	config *Config
}

// New creates a Restic wrapper instance for the provided config
func New(conf *Config) Restic {
	return Restic{config: conf}
}

// Run the configured restic jobs
func (r *Restic) Run() error {
	var err error

	for _, v := range r.config.Backup {
		err = r.runBackup(v)

		if err != nil {
			if v.ContinueOnError {
				log.WithError(err).Warnf("Backup job \"%s\" failed. Continuing...", v.Backup)
			} else {
				return errors.Wrap(err, "runBackup failed")
			}
		}
	}

	for _, v := range r.config.Forget {
		err = r.runForget(v)

		if err != nil {
			if v.ContinueOnError {
				log.WithError(err).Warnf("Forget job \"%s\" failed. Continuing...", v.Repository)
			} else {
				return errors.Wrap(err, "runForget failed")
			}
		}
	}

	return nil
}

func (r *Restic) runBackup(b Backup) error {
	log.WithFields(log.Fields{
		"backup":     b.Backup,
		"repository": b.Repository,
		"source":     b.Source,
		"exclude":    b.Exclude,
	}).Info("start runBackup")

	repo, exists := r.repository(b.Repository)
	if !exists {
		return errors.Errorf("repository \"%s\" does not exist", b.Repository)
	}

	args := []string{"backup"}
	args = append(args, b.Source)
	args = append(args, "--repo", repo.Path)

	for _, v := range b.Exclude {
		args = append(args, "--exclude", v)
	}

	err := r.execute(args, repo.Password)

	log.Info("end runBackup")
	return errors.Wrap(err, "execute failed")
}

func (r *Restic) runForget(f Forget) error {
	repo, exists := r.repository(f.Repository)
	if !exists {
		return errors.Errorf("repository \"%s\" does not exist", f.Repository)
	}

	args := []string{"forget"}
	args = append(args, "--repo", repo.Path)

	if f.Hostname != "" {
		args = append(args, "--hostname", f.Hostname)
	}

	if f.KeepLast > 0 {
		args = append(args, "--keep-last", strconv.Itoa(f.KeepLast))
	}

	if f.KeepHourly > 0 {
		args = append(args, "--keep-hourly", strconv.Itoa(f.KeepHourly))
	}

	if f.KeepDaily > 0 {
		args = append(args, "--keep-daily", strconv.Itoa(f.KeepDaily))
	}

	if f.KeepWeekly > 0 {
		args = append(args, "--keep-weekly", strconv.Itoa(f.KeepWeekly))
	}

	if f.KeepMonthly > 0 {
		args = append(args, "--keep-monthly", strconv.Itoa(f.KeepMonthly))
	}

	if f.KeepYearly > 0 {
		args = append(args, "--keep-yearly", strconv.Itoa(f.KeepYearly))
	}

	if len(f.KeepTag) > 0 {
		args = append(args, "--keep-tag", strings.Join(f.KeepTag, ","))
	}

	if len(f.Tag) > 0 {
		args = append(args, "--tag", strings.Join(f.Tag, ","))
	}

	if f.Prune {
		args = append(args, "--prune")
	}

	err := r.execute(args, repo.Password)
	return errors.Wrap(err, "execute failed")
}

func (r *Restic) repository(key string) (repo Repository, exists bool) {
	for _, v := range r.config.Repositoies {
		if v.Repository == key {
			return v, true
		}
	}

	return Repository{}, false
}

func (r *Restic) execute(arguments []string, password string) error {
	log.WithField("arguments", arguments).Info("Executing restic command")

	command := exec.Command("restic", arguments...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = append(os.Environ(), fmt.Sprintf("RESTIC_PASSWORD=%s", password))
	err := command.Run()

	return errors.Wrap(err, "restic exec failed")
}
