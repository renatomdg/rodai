package cmd

import (
	"fmt"
	"github.com/renatomdg/rodai/pkg"
	"github.com/spf13/cobra"
)
// FlowAddParameters should be used as a struct to use validator to avoid bad behaviors :)


var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Manage command flows adding, updating or deleting.",
}

var flowAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new command flow",
	Run:   flowAdd,
}

var flowUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates an existent command flow.",
	Run:   flowUpdate,
}

var flowDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an existent command flow.",
	Run:   flowDelete,
}

func init() {
	rootCmd.AddCommand(flowCmd)
	flowCmd.AddCommand(flowAddCmd, flowUpdateCmd, flowDeleteCmd)

	flowAddCmd.Flags().StringP("name", "n", "", "Name of the command flow.")
	flowAddCmd.Flags().BoolP("validate", "v", false, "Validate the command before store it in the flow.")
	flowAddCmd.Flags().StringP("import", "i", "", "Import commands from a file.")
	flowAddCmd.Flags().StringArrayP("commands", "c", []string{}, "Commands to execute.")
	flowAddCmd.Flags().BoolP("parallel", "p", false, "Set execution to be in parallel.")
	flowAddCmd.Flags().BoolP("serial", "s", true, "Set execution to be serial.")
	flowAddCmd.MarkFlagRequired("name")
	flowAddCmd.MarkFlagRequired("commands")

	flowDeleteCmd.Flags().StringP("name", "n", "", "Name of the flow to delete.")
	flowDeleteCmd.MarkFlagRequired("name")

	flowUpdateCmd.Flags().StringP("flow", "f", "", "Name of the  flow to update")
	flowUpdateCmd.Flags().StringP("name", "n", "", "New name of the command flow.")
	flowUpdateCmd.Flags().BoolP("validate", "v", false, "Validate the command before store it in the flow.")
	flowUpdateCmd.Flags().StringP("import", "i", "", "Import commands from a file.")
	flowUpdateCmd.Flags().StringArrayP("commands", "c", []string{}, "Commands to execute.")
	flowUpdateCmd.Flags().BoolP("parallel", "p", false, "Set execution to be in parallel.")
	flowUpdateCmd.Flags().BoolP("serial", "s", true, "Set execution to be serial.")
	flowUpdateCmd.MarkFlagRequired("flow")

}

func flowAdd(cmd *cobra.Command, args []string) {

	var newFlow pkg.Flow

	newFlow.Name, _ = cmd.Flags().GetString("name")
	newFlow.Validate, _ = cmd.Flags().GetBool("validate")
	newFlow.Import, _ = cmd.Flags().GetString("import")
	newFlow.Commands, _ = cmd.Flags().GetStringArray("commands")
	newFlow.Parallel, _ = cmd.Flags().GetBool("parallel")
	newFlow.Serial, _ = cmd.Flags().GetBool("serial")

	result, code := pkg.AddFlow(newFlow)
	fmt.Println(pkg.ColoredStatus(code), result)
}

func flowUpdate(cmd *cobra.Command, args []string) {

	var flow pkg.Flow

	flow.Name, _ = cmd.Flags().GetString("name")
	flow.Validate, _ = cmd.Flags().GetBool("validate")
	flow.Import, _ = cmd.Flags().GetString("import")
	flow.Commands, _ = cmd.Flags().GetStringArray("commands")
	flow.Parallel, _ = cmd.Flags().GetBool("parallel")
	flow.Serial, _ = cmd.Flags().GetBool("serial")
	toUpdate, _ := cmd.Flags().GetString("flow")

	result, code := pkg.UpdateFlow(toUpdate, flow)
	fmt.Println(pkg.ColoredStatus(code), result)
}

func flowDelete(cmd *cobra.Command, args []string) {
	deleteFlow, _ := cmd.Flags().GetString("name")
	result, code := pkg.DeleteFlow(deleteFlow)
	fmt.Println(pkg.ColoredStatus(code), result)
}
