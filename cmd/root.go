/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/before80/h2m/dart"
	"github.com/before80/h2m/docker"
	"github.com/before80/h2m/fiber"
	"github.com/before80/h2m/grpc"
	"github.com/before80/h2m/mysql"
	"github.com/before80/h2m/npmjs"
	"github.com/before80/h2m/postgresql"
	"github.com/before80/h2m/protocolBuffers"
	"github.com/before80/h2m/python"
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
		// 用法：h2m -m "get-started" -e "docker" -u http://ft.cn/docker_getstarted_241023.html -t "dist/docker_docs" -n "#sectiontree" -c "article" -b ""
		// 用法：h2m -m "manuals" -e "docker" -u http://ft.cn/docker_manuals_241023.html -t "dist/docker_docs" -n "#sectiontree" -c "article" -b ""
		// 用法：h2m -m "reference" -e "docker" -u http://ft.cn/docker_reference_241023.html -t "dist/docker_docs" -n "#sectiontree" -c "article" -b ""
		// 用法：h2m -m "samp" -e "docker" -u http://ft.cn/docker_samples_240404.html -t "dist/docker_docs" -n "#sectiontree" -c "article" -b ""
		case "probuf":
			protocolBuffers.Html2md(cmd, args)
		// 用法：h2m -m "docs" -e "probuf" -u https://protobuf.dev/ -t "dist/probuf_docs" -n "#td-section-nav" -c "div.td-content" -b "https://protobuf.dev/"

		case "npmjs":
			npmjs.Html2md(cmd, args)
			// 用法：h2m -e "npmjs" -u http://ft.cn/npmjs_230922.html -t "dist/npmjs_docs" -n "#menus" -c "#skip-nav" -b https://docs.npmjs.com/
		case "fiber":
			fiber.Html2md(cmd, args)
			// 用法：h2m -e "fiber" -u http://ft.cn/fiber_docs_240205.html -t "dist/Fiber_docs" -n "nav#nav" -c "div.body" -b https://docs.gofiber.io/
		case "python":
			python.Html2md(cmd, args)
			// 用法：h2m -m "tutorial" -e "python" -u https://docs.python.org/zh-cn/3.13/tutorial/index.html -t "dist/python_docs" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/tutorial/"
			// 用法：h2m -m "reference" -e "python" -u https://docs.python.org/zh-cn/3.13/reference/index.html -t "dist/python_docs" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/reference/"
			// 用法：h2m -m "library" -e "python" -u https://docs.python.org/zh-cn/3.13/library/index.html -t "dist/python_docs" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "howto" -e "python" -u https://docs.python.org/zh-cn/3.13/howto/index.html -t "dist/python_docs" -n "section" -c "section" -b "https://docs.python.org/zh-cn/3.13/howto/"
			// 用法：h2m -m "using" -e "python" -u https://docs.python.org/zh-cn/3.13/using/index.html -t "dist/python_docs" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/using/"
			// 用法：h2m -m "faq" -e "python" -u https://docs.python.org/zh-cn/3.13/faq/index.html -t "dist/python_docs" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/faq/"
			// 用法：h2m -m "c_api" -e "python" -u https://docs.python.org/zh-cn/3.13/c-api/index.html -t "dist/python_docs" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/c-api/"
			// 用法：h2m -m "text" -e "python" -u https://docs.python.org/zh-cn/3.13/library/text.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "binary" -e "python" -u https://docs.python.org/zh-cn/3.13/library/binary.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "datatypes" -e "python" -u https://docs.python.org/zh-cn/3.13/library/datatypes.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "numeric" -e "python" -u https://docs.python.org/zh-cn/3.13/library/numeric.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "functional" -e "python" -u https://docs.python.org/zh-cn/3.13/library/functional.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "filesys" -e "python" -u https://docs.python.org/zh-cn/3.13/library/filesys.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "persistence" -e "python" -u https://docs.python.org/zh-cn/3.13/library/persistence.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "archiving" -e "python" -u https://docs.python.org/zh-cn/3.13/library/archiving.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "fileformats" -e "python" -u https://docs.python.org/zh-cn/3.13/library/fileformats.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "crypto" -e "python" -u https://docs.python.org/zh-cn/3.13/library/crypto.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "allos" -e "python" -u https://docs.python.org/zh-cn/3.13/library/allos.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "concurrency" -e "python" -u https://docs.python.org/zh-cn/3.13/library/concurrency.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "ipc" -e "python" -u https://docs.python.org/zh-cn/3.13/library/ipc.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "netdata" -e "python" -u https://docs.python.org/zh-cn/3.13/library/netdata.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "markup" -e "python" -u https://docs.python.org/zh-cn/3.13/library/markup.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "internet" -e "python" -u https://docs.python.org/zh-cn/3.13/library/internet.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "mm" -e "python" -u https://docs.python.org/zh-cn/3.13/library/mm.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "i18n" -e "python" -u https://docs.python.org/zh-cn/3.13/library/i18n.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "frameworks" -e "python" -u https://docs.python.org/zh-cn/3.13/library/frameworks.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "tk" -e "python" -u https://docs.python.org/zh-cn/3.13/library/tk.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "development" -e "python" -u https://docs.python.org/zh-cn/3.13/library/development.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "debug" -e "python" -u https://docs.python.org/zh-cn/3.13/library/debug.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "distribution" -e "python" -u https://docs.python.org/zh-cn/3.13/library/distribution.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "python" -e "python" -u https://docs.python.org/zh-cn/3.13/library/python.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "custominterp" -e "python" -u https://docs.python.org/zh-cn/3.13/library/custominterp.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "modules" -e "python" -u https://docs.python.org/zh-cn/3.13/library/modules.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "language" -e "python" -u https://docs.python.org/zh-cn/3.13/library/language.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "windows" -e "python" -u https://docs.python.org/zh-cn/3.13/library/windows.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "unix" -e "python" -u https://docs.python.org/zh-cn/3.13/library/unix.html -t "dist/python_docs/library" -n "div.compound" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "cmdline" -e "python" -u https://docs.python.org/zh-cn/3.13/library/cmdline.html -t "dist/python_docs/library" -n "section" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
			// 用法：h2m -m "superseded" -e "python" -u https://docs.python.org/zh-cn/3.13/library/superseded.html -t "dist/python_docs/library" -n "section" -c "section" -b "https://docs.python.org/zh-cn/3.13/library/"
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
