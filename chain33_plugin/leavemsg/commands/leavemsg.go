package commands

import "github.com/spf13/cobra"

// 本执行器的命令行初始化总入口
func LeaveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "echo",
		Short: "echo commandline interface",
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
		SendCmd(),  // 查询消息记录
		// 如果有其它命令，在这里加入
	)
	return cmd
}

// 本执行器的命令行初始化总入口
func SendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "echo",
		Short: "echo commandline interface",
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
		QueryCmd(),  // 查询消息记录
		// 如果有其它命令，在这里加入
	)
	return cmd
}