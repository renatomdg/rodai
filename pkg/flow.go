package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"

	"gopkg.in/yaml.v3"
)

type Flow struct {
	Name         string   `yaml:"Name"`
	Parallel     bool     `yaml:"Parallel"`
	Serial       bool     `yaml:"Serial"`
	Commands     []string `yaml:"Commands"`
	UseBastion   bool     `yaml:"UseBastion"`
	StoreResults bool     `yaml:"StoreResults"`
	BastionName  string   `yaml:"BastionName"`
}

type Flows struct {
	Flows []Flow `yaml:"Flows"`
}

var (
	ConfigFilePath string
	FlowsFilePath  string
	ConfigPath     string
	home           string
)

func init() {
	home, _ = homedir.Dir()
	ConfigPath = home + "/.config/rodai"
	ConfigFilePath = home + "/.config/rodai/config.yml"
	FlowsFilePath = home + "/.config/rodai/flows.yml"

}

func AddFlow(flow Flow) (string, int) {
	flows := GetFlows(true)

	for _, f := range flows.Flows {
		if f.Name == flow.Name {
			return "Failed creating flow " + flow.Name + " (Already exists).", 1
		}
	}

	flows.Flows = append(flows.Flows, flow)

	CommitFlow(flows)
	return "Success creating flow " + flow.Name, 0
}

func UpdateFlow(flowName string, toUpdate Flow) (string, int) {
	flows := GetFlows(false)
	for _, flow := range flows.Flows {
		if flow.Name == flowName {
			DeleteFlow(flow.Name)
			flow.Name = toUpdate.Name
			flow.Commands = toUpdate.Commands
			flow.Serial = toUpdate.Serial
			if toUpdate.Serial && toUpdate.Parallel {
				flow.Serial = false
				flow.Parallel = true
			}
			flow.StoreResults = toUpdate.StoreResults
			flow.Parallel = toUpdate.Parallel
			AddFlow(flow)
			return fmt.Sprintf("Update flow %s (now %s)", flowName, toUpdate.Name), 0
		}
	}

	return "Failed updating flow " + flowName + " (Not found)", 1
}

func DeleteFlow(toDelete string) (string, int) {
	flows := GetFlows(false)
	for i, flow := range flows.Flows {
		if flow.Name == toDelete {
			flows.Flows = append(flows.Flows[:i], flows.Flows[i+1:]...)
			CommitFlow(flows)
			return "Deleting flow " + toDelete, 0
		}
	}
	return "Failed deleting flow " + toDelete + " (Not found)", 1
}

func GetFlowParameters(name string) (Flow, error) {
	var flows Flows

	_, err := os.OpenFile(FlowsFilePath, os.O_RDONLY|os.O_CREATE, 0644)
	CheckErr(err)

	flowContent := FileContentAsByteArray(FlowsFilePath)
	CheckErr(err)

	err = yaml.Unmarshal(flowContent, &flows)
	CheckErr(err)

	for _, flow := range flows.Flows {
		if flow.Name == name {
			return flow, nil
		}
	}

	return Flow{}, errors.New("flow not found")
}

func GetFlows(isAdd bool) Flows {

	var flows Flows

	err := os.MkdirAll(ConfigPath, os.ModePerm)
	CheckErr(err)

	_, err = os.OpenFile(FlowsFilePath, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	CheckErr(err)

	flowContent := FileContentAsByteArray(FlowsFilePath)
	CheckErr(err)

	if len(flowContent) == 0 && !isAdd {
		fmt.Println(ColoredStatus(16), "No Flows yet registered.")
		os.Exit(0)
	}

	err = yaml.Unmarshal(flowContent, &flows)
	CheckErr(err)

	return flows
}

func CommitFlow(flows Flows) {
	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)

	err := encoder.Encode(&flows)
	CheckErr(err)

	file, err := os.OpenFile(FlowsFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	CheckErr(err)
	defer file.Close()

	_, err = file.WriteString(b.String())
	CheckErr(err)
}
