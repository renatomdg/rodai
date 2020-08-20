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

type ParallelFlowExec struct {
	Command string
	Name    string
	Bastion string
}

type FlowCommandReturn struct {
	Flow    string
	Command string
	Status  int
	Output  string
}

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

	if selectedFlow.Parallel {
		ParallelExecutor(selectedFlow)
	} else if selectedFlow.Serial {
		SerialExecutor(selectedFlow)
	}

}

func ParallelExecutor(selectedFlow Flow) {

	var outFile string

	results := make(map[string]string)
	flowChan := make(chan ParallelFlowExec)
	resultsChan := make(chan FlowCommandReturn)

	for i := 0; i <= len(selectedFlow.Commands); i++ {
		go RunFlowCommandsParallel(flowChan, resultsChan)
	}

	for _, command := range selectedFlow.Commands {
		fmt.Println(ColoredStatus(16), "Running command:", command)
		flowChan <- ParallelFlowExec{
			Command: command,
			Name:    selectedFlow.Name,
			Bastion: selectedFlow.BastionName,
		}
	}

	breakExecIn := len(selectedFlow.Commands)
	currentExec := 0
	for {
		if currentExec == breakExecIn {
			break
		}

		select {
		case res := <-resultsChan:
			fmt.Println("Running command:", res.Command, ColoredStatus(res.Status))
			results[res.Command] = res.Output
			currentExec++
		default:
		}
	}

	outFile = GenerateFileName(selectedFlow.Name)
	for command, output := range results {
		if selectedFlow.StoreResults {
			WriteInFile(outFile, "Command: "+command)
			WriteInFile(outFile, output)
		} else {
			fmt.Println("Command: "+command)
			fmt.Println(output)
		}
	}

	if selectedFlow.StoreResults {
		fmt.Printf(ColoredStatus(15)+" Commands stdout/err stored in ~/.config/rodai/runs/%s", outFile)
	}
}

func SerialExecutor(selectedFlow Flow) {
	results := make(map[string]string)
	var outFile string

	for _, command := range selectedFlow.Commands {
		values := strings.Split(command, " ")
		bin := values[0]
		args := values[1:]
		s := spinner.New(spinner.CharSets[6], 100*time.Millisecond)
		s.Suffix = fmt.Sprintf(" Running command: %s", command)
		s.Color("cyan")
		s.Start()
		out, status := ExecCommand(selectedFlow.BastionName, bin, args)
		results[command] = out
		fmt.Println(ColoredStatus(status))
		s.Stop()
	}

	outFile = GenerateFileName(selectedFlow.Name)
	for command, output := range results {
		if selectedFlow.StoreResults {
			WriteInFile(outFile, "Command: "+command)
			WriteInFile(outFile, output)
		} else {
			fmt.Println("Command: "+command)
			fmt.Println(output)
		}
	}

	if selectedFlow.StoreResults {
		fmt.Printf(ColoredStatus(15)+" Commands stdout/err stored in ~/.config/rodai/runs/%s", outFile)
	}
}

func RunFlowCommandsParallel(flowChan <-chan ParallelFlowExec, results chan<- FlowCommandReturn) {

	for flow := range flowChan {

		values := strings.Split(flow.Command, " ")
		bin := values[0]
		args := values[1:]
		out, status := ExecCommand(flow.Bastion, bin, args)
		execReturn := FlowCommandReturn{}
		execReturn.Status = status
		execReturn.Output = out
		execReturn.Command = flow.Command
		execReturn.Flow = flow.Name

		results <- execReturn
	}
}

func ExecCommand(bastionName, command string, args []string) (string, int) {

	if len(bastionName) > 0 {
		bastion, err := GetBastionConfigDetails(bastionName)
		CheckErr(err)
		fullCommand := command + " " + strings.Join(args, " ")
		cmd := exec.Command("ssh", "-i",
			bastion.Key,
			"-oStrictHostKeyChecking=no",
			bastion.Host,
			fullCommand,
		)

		output, _ := cmd.CombinedOutput()
		switch cmd.ProcessState.ExitCode() {
		case 0:
			return string(output), 0
		case 1:
			return string(output), 1
		default:
			return string(output), 3012
		}
	}

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
