package cmd

import (
	"fmt"
	"github.com/renatomdg/rodai/pkg"
	"github.com/spf13/cobra"
)

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
	flowAddCmd.Flags().StringArrayP("commands", "c", []string{}, "Commands to execute.")
	flowAddCmd.Flags().BoolP("parallel", "p", false, "Set execution to be in parallel.")
	flowAddCmd.Flags().Bool("store-results", false, "Store the stdout/err in a file")
	flowAddCmd.Flags().BoolP("serial", "s", true, "Set execution to be serial.")
	flowAddCmd.Flags().BoolP("use-bastion", "u", false, "Use a bastion to execute commands.")
	flowAddCmd.Flags().StringP("bastion-name", "b", "", "Specify the bastion to run commands over.")
	flowAddCmd.MarkFlagRequired("name")
	flowAddCmd.MarkFlagRequired("commands")

	flowDeleteCmd.Flags().StringP("name", "n", "", "Name of the flow to delete.")
	flowDeleteCmd.MarkFlagRequired("name")

	flowUpdateCmd.Flags().StringP("flow", "f", "", "Name of the  flow to update")
	flowUpdateCmd.Flags().StringP("name", "n", "", "New name of the command flow.")
	flowUpdateCmd.Flags().StringArrayP("commands", "c", []string{}, "Commands to execute.")
	flowUpdateCmd.Flags().BoolP("parallel", "p", false, "Set execution to be in parallel.")
	flowUpdateCmd.Flags().Bool("store-results", false, "Store the stdout/err in a file")
	flowUpdateCmd.Flags().BoolP("serial", "s", true, "Set execution to be serial.")
	flowUpdateCmd.Flags().BoolP("use-bastion", "u", false, "Use a bastion to execute commands.")
	flowUpdateCmd.Flags().StringP("bastion-name", "b", "", "Specify the bastion to run commands over.")
	flowUpdateCmd.MarkFlagRequired("flow")
	flowUpdateCmd.MarkFlagRequired("name")

}

func flowAdd(cmd *cobra.Command, args []string) {

	var flow pkg.Flow

	flow.Name, _ = cmd.Flags().GetString("name")
	flow.Commands, _ = cmd.Flags().GetStringArray("commands")
	flow.Parallel, _ = cmd.Flags().GetBool("parallel")
	flow.Serial, _ = cmd.Flags().GetBool("serial")
	flow.UseBastion, _ = cmd.Flags().GetBool("use-bastion")
	flow.BastionName, _ = cmd.Flags().GetString("bastion-name")
	flow.StoreResults, _ = cmd.Flags().GetBool("store-results")

	result, code := pkg.AddFlow(flow)
	fmt.Println(pkg.ColoredStatus(code), result)
}

func flowUpdate(cmd *cobra.Command, args []string) {

	var flow pkg.Flow

	flow.Name, _ = cmd.Flags().GetString("name")
	flow.Commands, _ = cmd.Flags().GetStringArray("commands")
	flow.Parallel, _ = cmd.Flags().GetBool("parallel")
	flow.Serial, _ = cmd.Flags().GetBool("serial")
	flow.UseBastion, _ = cmd.Flags().GetBool("use-bastion")
	flow.BastionName, _ = cmd.Flags().GetString("bastion-name")
	flow.StoreResults, _ = cmd.Flags().GetBool("store-results")
	toUpdate, _ := cmd.Flags().GetString("flow")

	result, code := pkg.UpdateFlow(toUpdate, flow)
	fmt.Println(pkg.ColoredStatus(code), result)
}

func flowDelete(cmd *cobra.Command, args []string) {
	deleteFlow, _ := cmd.Flags().GetString("name")
	result, code := pkg.DeleteFlow(deleteFlow)
	fmt.Println(pkg.ColoredStatus(code), result)
}
