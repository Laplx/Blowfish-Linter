package cmd

import (
	"fmt"
	"os"

	"github.com/Laplx/Blowfish-Linter/config"
	"github.com/Laplx/Blowfish-Linter/formatter"

	"github.com/spf13/cobra"
)

var (
	within     bool
	force      bool
	inspect    bool
	defaultDir string
)

var rootCmd = &cobra.Command{
	Use:   "bflint [directory]",
	Short: "Blowfish Markdown 格式检查与转换工具",
	Long: `Blowfish Markdown 格式检查与转换工具

支持 Front Matter 自动补全或强制覆盖、LaTeX 公式格式转换，以及子文件夹结构规范。

Notice:
  - 配置文件控制具体转换规则，请在运行前确认 $USERPROFILE/.config/bflint/ 或当前目录下的 config.yaml 配置正确
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetDir := args[0]
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "错误：目录 %s 不存在\n", targetDir)
			os.Exit(1)
		}
		config.LoadConfig()
		formatter.ProcessDir(targetDir, within, force, inspect, defaultDir)
	},
}

func Execute() {
	rootCmd.Flags().BoolVarP(&within, "within", "w", false, "仅处理已有 Front Matter 的文件")
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "已存在字段强制覆盖")
	rootCmd.Flags().BoolVarP(&inspect, "inspect", "i", false, "检查并整理子文件夹结构，规范为同名文件夹下 index.md 格式")
	rootCmd.Flags().StringVarP(&defaultDir, "default", "d", "", "指定默认文件来源文件夹，配合 --inspect 使用")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}
