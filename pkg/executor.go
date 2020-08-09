package pkg

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"os/exec"
	"strings"
	"time"
)

func Executor(cmd *cobra.Command, args []string) {
	flows := GetFlows(false)
	var flowsList []string

	for _, flow := range flows.Flows {
		flowsList = append(flowsList, flow.Name)
	}

	prompt := promptui.Select{
		Label: "Select Flow",
		Items: flowsList,
	}

	_, result, err := prompt.Run()
	CheckErr(err)
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	selectedFlow, err := GetFlowParameters(result)
	CheckErr(err)

	RunFlowCommands(selectedFlow)

}

func RunFlowCommands(flow Flow) {

	results := make(map[string]string)
	var outFile string


	fmt.Print("\033[H\033[2J")
	for _, command := range flow.Commands {
		values := strings.Split(command, " ")
		bin := values[0]
		args := values[1:]
		s := spinner.New(spinner.CharSets[6], 100*time.Millisecond)
		s.Suffix = fmt.Sprintf(" Running command: %s", command)
		s.Color("cyan")
		s.Start()
		out, status := ExecCommand(bin, strings.Join(args, " "))
		results[command] = out
		fmt.Println(ColoredStatus(status))
		s.Stop()
	}

	outFile = GenerateFileName(flow.Name)
	for command, output := range results {
		WriteInFile(outFile, "Command: "+command)
		WriteInFile(outFile, output)
	}
	fmt.Printf(ColoredStatus(15) + " Commands stdout/err stored in ~/.config/rodai/runs/%s", outFile)
}

func ExecCommand(command string, args ...string) (string, int) {

	cmd := exec.Command(command, args...)
	output, _ := cmd.CombinedOutput()

	switch cmd.ProcessState.ExitCode() {
	case 0:
		return string(output), 0
	case 1:
		return string(output), 1
	default:
		return string(output), 3012
	}

	return string(output), 3012
}

