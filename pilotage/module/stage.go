/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package module

import (
	"fmt"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
)

const (
	// Stage Type
	StartStage  = "start"
	EndStage    = "end"
	NormalStage = "normal"
	PauseStage  = "pause"
)

// Stage is
type Stage struct {
	T          string   `json:"type" yaml:"type"`
	Name       string   `json:"name" yaml:"name"`
	Title      string   `json:"title" yaml:"title"`
	Sequencing string   `json:"sequencing,omitempty" yaml:"sequencing,omitempty"`
	Status     string   `json:"status,omitempty" yaml:"status,omitempty"`
	Logs       []string `json:"logs,omitempty" yaml:"logs,omitempty"`
	Actions    []Action `json:"actions,omitempty" yaml:"actions,omitempty"`
}

// TODO filter the log print with different color.
func (s *Stage) Log(log string, verbose, timestamp bool) {
	s.Logs = append(s.Logs, fmt.Sprintf("[%s] %s", time.Now().String(), log))

	if verbose == true {
		if timestamp == true {
			fmt.Println(Cyan(fmt.Sprintf("[%s] %s", time.Now().String(), strings.TrimSpace(log))))
		} else {
			fmt.Println(Cyan(log))
		}
	}
}

func (s *Stage) SequencingRun(verbose, timestamp bool, f *Flow) (string, error) {
	s.Status = Running

	s.Log(fmt.Sprintf("Stage [%s] status change to %s", s.Name, s.Status), false, timestamp)
	f.Log(fmt.Sprintf("Stage [%s] status change to %s", s.Name, s.Status), verbose, timestamp)

	for i, _ := range s.Actions {
		action := &s.Actions[i]

		s.Log(fmt.Sprintf("The Number [%d] action is running: %s", i, s.Title), false, timestamp)
		f.Log(fmt.Sprintf("The Number [%d] action is running: %s", i, s.Title), verbose, timestamp)

		if status, err := action.Run(verbose, timestamp, f); err != nil {
			s.Status = Failure

			s.Log(fmt.Sprintf("Action [%s] run error: %s", action.Name, err.Error()), false, timestamp)
			f.Log(fmt.Sprintf("Action [%s] run error: %s", action.Name, err.Error()), verbose, timestamp)

		} else {
			s.Status = status
		}

		if s.Status == Failure || s.Status == Cancel {
			break
		}
	}

	return s.Status, nil
}
