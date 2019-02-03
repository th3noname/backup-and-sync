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

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Repositoies []Repository `mapstructure:"repositories"`
	Backup      []Backup     `mapstructure:"backups"`
	Forget      []Forget     `mapstructure:"forget"`
}

type Repository struct {
	Repository string `mapstructure:"repository"`
	Path       string `mapstructure:"path"`
	Password   string `mapstructure:"password"`
}

type Backup struct {
	Backup     string   `mapstructure:"backup"`
	Repository string   `mapstructure:"repository"`
	Source     string   `mapstructure:"source"`
	Exclude    []string `mapstructure:"no_backup"`
}

type Forget struct {
	Repository  string   `mapstructure:"repository"`
	Prune       bool     `mapstructure:"prune"`
	KeepLast    int      `mapstructure:"keep-last"`
	KeepHourly  int      `mapstructure:"keep-hourly"`
	KeepDaily   int      `mapstructure:"keep-daily"`
	KeepWeekly  int      `mapstructure:"keep-weekly"`
	KeepMonthly int      `mapstructure:"keep-monthly"`
	KeepYearly  int      `mapstructure:"keep-yearly"`
	KeepTag     []string `mapstructure:"keep-tag"`
}

type Restic struct {
	config *Config
}

func New(conf *Config) Restic {
	return Restic{config: conf}
}

func (r *Restic) Run() error {
	var err error

	for _, v := range r.config.Backup {
		err = r.runBackup(v)

		if err != nil {
			return errors.Wrap(err, "runBackup failed")
		}
	}

	for _, v := range r.config.Forget {
		err = r.runForget(v)

		if err != nil {
			return errors.Wrap(err, "runForget failed")
		}
	}

	return nil
}

func (r *Restic) runBackup(b Backup) error {
	repo, exists := r.repository(b.Repository)
	if !exists {
		return errors.New(fmt.Sprintf("repository \"%s\" does not exist", b.Repository))
	}

	args := []string{"backup"}

	for _, v := range b.Exclude {
		args = append(args, "--exclude", v)
	}

	args = append(args, "--repo", repo.Path)
	args = append(args, b.Source)

	err := r.execute(args, repo.Password)

	return errors.Wrap(err, "execute failed")
}

func (r *Restic) runForget(f Forget) error {

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

	return errors.Wrap(err, "Restic exec failed")
}
