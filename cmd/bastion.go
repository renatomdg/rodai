package cmd

import (
	"fmt"

	"github.com/renatomdg/rodai/pkg"
	"github.com/spf13/cobra"
)

var bastionCmd = &cobra.Command{
	Use:   "bastion",
	Short: "Manage bastion configs.",
}

var bastionAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new bastion config.",
	Run:   bastionAdd,
}

var bastionUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates an existent bastion config.",
	Run:   bastionUpdate,
}

var bastionDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an existent bastion config.",
	Run:   bastionDelete,
}

func init() {
	rootCmd.AddCommand(bastionCmd)
	bastionCmd.AddCommand(bastionAddCmd, bastionUpdateCmd, bastionDeleteCmd)

	bastionAddCmd.Flags().StringP("name", "n", "", "Set the name for bastion config.")
	bastionAddCmd.Flags().StringP("user", "u", "root", "Set the bastion ssh user.")
	bastionAddCmd.Flags().String("host", "c", "Set host (DNS or IP).")
	bastionAddCmd.Flags().IntP("port", "p", 22, "Set port.")
	bastionAddCmd.Flags().StringP("key", "k", "", "Specify the SSH Key to be used.")
	bastionAddCmd.MarkFlagRequired("name")
	bastionAddCmd.MarkFlagRequired("user")
	bastionAddCmd.MarkFlagRequired("host")
	bastionAddCmd.MarkFlagRequired("port")
	bastionAddCmd.MarkFlagRequired("key")

	bastionUpdateCmd.Flags().StringP("bastion", "b", "", "Specify the bastion config do update.")
	bastionUpdateCmd.Flags().StringP("name", "n", "", "Set the name for bastion config.")
	bastionUpdateCmd.Flags().StringP("user", "u", "root", "Set the bastion ssh user.")
	bastionUpdateCmd.Flags().String("host", "c", "Set host (DNS or IP).")
	bastionUpdateCmd.Flags().IntP("port", "p", 22, "Set port.")
	bastionUpdateCmd.Flags().StringP("key", "k", "", "Specify the SSH Key to be used.")
	bastionUpdateCmd.MarkFlagRequired("bastion")
	bastionUpdateCmd.MarkFlagRequired("name")
	bastionUpdateCmd.MarkFlagRequired("user")
	bastionUpdateCmd.MarkFlagRequired("port")
	bastionUpdateCmd.MarkFlagRequired("host")
	bastionUpdateCmd.MarkFlagRequired("key")

	bastionDeleteCmd.Flags().StringP("name", "n", "", "Specify the bastion config to delete.")
	bastionDeleteCmd.MarkFlagRequired("name")
}

func bastionAdd(cmd *cobra.Command, args []string) {

	var bastion pkg.Bastion

	bastion.Name, _ = cmd.Flags().GetString("name")
	bastion.Username, _ = cmd.Flags().GetString("user")
	bastion.Host, _ = cmd.Flags().GetString("host")
	bastion.Port, _ = cmd.Flags().GetInt("port")
	bastion.Key, _ = cmd.Flags().GetString("key")

	result, code := pkg.AddBastion(bastion)
	fmt.Println(pkg.ColoredStatus(code), result)
}

func bastionUpdate(cmd *cobra.Command, args []string) {

	var bastion pkg.Bastion
	toUpdate, _ := cmd.Flags().GetString("bastion")
	bastion.Name, _ = cmd.Flags().GetString("name")
	bastion.Username, _ = cmd.Flags().GetString("user")
	bastion.Host, _ = cmd.Flags().GetString("host")
	bastion.Port, _ = cmd.Flags().GetInt("port")
	bastion.Key, _ = cmd.Flags().GetString("key")

	result, code := pkg.UpdateBastion(bastion, toUpdate)
	fmt.Println(pkg.ColoredStatus(code), result)
}

func bastionDelete(cmd *cobra.Command, args []string) {
	deleteFlow, _ := cmd.Flags().GetString("name")
	result, code := pkg.DeleteFlow(deleteFlow)
	fmt.Println(pkg.ColoredStatus(code), result)
}
