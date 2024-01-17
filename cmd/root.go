/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/before80/h2m/dart"
	"github.com/before80/h2m/docker"
	"github.com/before80/h2m/grpc"
	"github.com/before80/h2m/mysql"
	"github.com/before80/h2m/npmjs"
	"github.com/before80/h2m/postgresql"
	"github.com/before80/h2m/vscode"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "h2m",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		t, _ := cmd.Flags().GetString("type")
		switch t {
		case "mysql":
			mysql.Html2md(cmd, args)
			// 用法：h2m -e "mysql" -u http://ft.cn/nav3.html -t "dist/workbench" -n "#doc-22" -c "#docs-body" -b https://dev.mysql.com/doc/workbench/en/
			// 用法：h2m -e "mysql" -u http://ft.cn/nav2.html -t "dist/MySQL57" -n "#doc-12" -c "#docs-body" -b https://dev.mysql.com/doc/refman/5.7/en/
		case "dart":
			dart.Html2md(cmd, args)
		// 用法：h2m -e "dart" -u http://ft.cn/dart_docs_230826.html -t "dist/Dart_docs" -n "div.site-sidebar" -c "div.content" -b https://dart.dev/
		case "vscode":
			vscode.Html2md(cmd, args)
		// 用法：h2m -e "vscode" -u http://ft.cn/vscode_docs_240112.html -t "dist/Vscode_docs" -n "nav#docs-navbar" -c "div.body" -b https://code.visualstudio.com/
		case "grpc":
			grpc.Html2md(cmd, args)
		// 用法：h2m -e "grpc" -u http://ft.cn/grpc_docs_240117.html -t "dist/grpc_docs" -n "nav#td-section-nav" -c "main" -b https://grpc.io/
		case "postgresql":
			postgresql.Html2md(cmd, args)
		// 用法：h2m -e "postgresql" -u http://ft.cn/postgresql_docs_15_4.html -t "dist/Postgresql_docs" -n "div.toc" -c "div#docContent" -b https://www.postgresql.org/docs/current/
		case "docker":
			docker.Html2md(cmd, args)
		// 用法：h2m -m "Guides" -e "docker" -u http://ft.cn/docker_guides_230830.html -t "dist/docker_docs" -n "#sectiontree" -c "article" -b https://docs.docker.com/
		// 用法：h2m -m "Manuals" -e "docker" -u http://ft.cn/docker_manuals_230830.html -t "dist/docker_docs" -n "#sectiontree" -c "article" -b https://docs.docker.com/
		// 用法：h2m -m "Reference" -e "docker" -u http://ft.cn/docker_reference_230830.html -t "dist/docker_docs" -n "#sectiontree" -c "article" -b https://docs.docker.com/
		case "npmjs":
			npmjs.Html2md(cmd, args)
			// 用法：h2m -e "npmjs" -u http://ft.cn/npmjs_230922.html -t "dist/npmjs_docs" -n "#menus" -c "#skip-nav" -b https://docs.npmjs.com/
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.h2m.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("opt-url", "u", "", "所需转换的html的所在URL")
	rootCmd.Flags().StringP("opt-nav-selector", "n", "", "html选择器1")
	rootCmd.Flags().StringP("opt-content-selector", "c", "", "html选择器2")
	rootCmd.Flags().StringP("dist", "t", "./dist", "目标路径")
	rootCmd.Flags().StringP("base-url", "b", "", "基本url")
	rootCmd.Flags().StringP("type", "e", "", "抓取页面类型")
	rootCmd.Flags().StringP("menuname", "m", "", "文档子菜单")

	//rootCmd.Flags().StringP("opt-heading-style", "t", "atx", "标题样式，可选值：setext 和 atx， 默认值: atx")
	//rootCmd.Flags().StringP("opt-horizontal-rule", "r", "", "水平行断句规则，可选值：任意值，默认值：为空")
	//rootCmd.Flags().StringP("opt-bullet-list-marker", "b", "-", "列表项标记，可选值：-、+、*，默认值：-")
	//rootCmd.Flags().StringP("opt-code-block-style", "c", "fenced", "列表项标记，可选值：indented 和 fenced，默认值：fenced")
	//rootCmd.Flags().StringP("opt-fence", "f", "```", "转换时表示代码框所需的字符，可选值：任意值，默认是```")
	//rootCmd.Flags().StringP("opt-em-delimiter", "e", "*", "转换时将<em>标签中的文本的左右加上所需的字符，可选值：任意值，默认值：*")
	//rootCmd.Flags().StringP("opt-strong-delimiter", "s", "**", "转换时将<strong>标签中的文本的左右加上所需的字符，可选值：任意值，默认值：**")
	//rootCmd.Flags().StringP("opt-link-style", "l", "inlined", "链接样式，可选值：inlined（即在行内）和 referenced（未知，待弄清楚），默认值：inlined")
	//rootCmd.Flags().StringP("opt-link-reference-style", "k", "full", "链接参考样式，可选值：full、collapsed、 shortcut，默认值：full")
	//rootCmd.Flags().StringP("opt-escape-mode", "m", "basic", "转义模式，可选值：basic 和 disabled，默认值：basic")

	rootCmd.MarkFlagRequired("opt-url")
	rootCmd.MarkFlagRequired("opt-nav-selector")
	rootCmd.MarkFlagRequired("type")
}
