// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"strings"

	ignTypes "github.com/coreos/ignition/config/v2_0/types"
	"github.com/coreos/ignition/config/validate"
	"github.com/coreos/ignition/config/validate/report"

	"github.com/coreos/container-linux-config-transpiler/config/templating"
)

type Systemd struct {
	Units []SystemdUnit `yaml:"units"`
}

type SystemdUnit struct {
	Name     string              `yaml:"name"`
	Enable   bool                `yaml:"enable"`
	Mask     bool                `yaml:"mask"`
	Contents string              `yaml:"contents"`
	DropIns  []SystemdUnitDropIn `yaml:"dropins"`
}

type SystemdUnitDropIn struct {
	Name     string `yaml:"name"`
	Contents string `yaml:"contents"`
}

func init() {
	register2_0(func(in Config, ast validate.AstNode, out ignTypes.Config, platform string) (ignTypes.Config, report.Report, validate.AstNode) {
		for _, unit := range in.Systemd.Units {
			if templating.HasTemplating([]string{unit.Contents}) {
				contents, err := templating.PerformTemplating(platform, []string{unit.Contents})
				if err != nil {
					return out, report.Report{}, ast
				}
				unit.Contents = strings.Join(contents, "\n")
			}
			newUnit := ignTypes.SystemdUnit{
				Name:     ignTypes.SystemdUnitName(unit.Name),
				Enable:   unit.Enable,
				Mask:     unit.Mask,
				Contents: unit.Contents,
			}

			for _, dropIn := range unit.DropIns {
				if templating.HasTemplating([]string{dropIn.Contents}) {
					contents, err := templating.PerformTemplating(platform, []string{dropIn.Contents})
					if err != nil {
						return out, report.Report{}, ast
					}
					dropIn.Contents = strings.Join(contents, "\n")
				}
				newUnit.DropIns = append(newUnit.DropIns, ignTypes.SystemdUnitDropIn{
					Name:     ignTypes.SystemdUnitDropInName(dropIn.Name),
					Contents: dropIn.Contents,
				})
			}

			out.Systemd.Units = append(out.Systemd.Units, newUnit)
		}
		return out, report.Report{}, ast
	})
}
