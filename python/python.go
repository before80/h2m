package python

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func fileNameFromMenuText(str string) string {
	re1, _ := regexp.Compile(`\.`)
	str = re1.ReplaceAllString(str, "_")

	re2, _ := regexp.Compile(` — `)
	str = re2.ReplaceAllString(str, "_")

	patternSlice := []string{
		`·`,
		`!`, `！`, `=`, `@`, `#`, `\$`,
		`\^`, `&`, `\*`,
		`\(`, `\)`,
		`（`, `）`,
		`\+`, `:`, `：`, `;`, `；`,
		`'`, `"`,
		`,`, `，`,
		`<`, `>`,
		`《`, `》`,
		`\?`, `？`,
		`\/`, `\|`,
		`—`,
	}

	for _, pattern := range patternSlice {
		re, _ := regexp.Compile(pattern)
		str = re.ReplaceAllString(str, "")
	}

	re3, _ := regexp.Compile(`\s+`)
	str = re3.ReplaceAllString(str, "")

	return str
}

func fileNameFromMenuLink(str string) string {
	re1, _ := regexp.Compile(`\.`)
	locS := re1.FindAllStringIndex(str, -1)
	// 找到最后一个点的所在位置
	if locS != nil {
		str = str[0:locS[len(locS)-1][0]]
		if re1.MatchString(str) {
			return strings.Replace(str, ".", "_", -1)
		}
		return str
	}
	return str
}

func removeSpace(str string) string {
	// 若有换行的情况，则去掉换行符
	re1, _ := regexp.Compile(`\n`)
	if re1.MatchString(str) {
		str = strings.Replace(str, "\n", "", -1)
	}

	// 若有空格的情况，则去掉空格
	re2, _ := regexp.Compile(`\s+`)
	str = re2.ReplaceAllString(str, " ")
	str = strings.TrimSpace(str)
	return str
}

// 获取li标签下第一个a标签的文本
func findATextAndHref(level int, liSel *goquery.Selection) (string, string) {
	aSel := liSel.ChildrenFiltered("a").First()
	if aSel.Length() == 0 {
		aSel = liSel.ChildrenFiltered("p").ChildrenFiltered("a").First()
	}

	if aSel.Length() == 0 {
		//fmt.Println(liSel.Html())
		log.Fatal("第" + strconv.Itoa(level) + "层li标签下未发现a标签")
	}
	menuLink, existAttr := aSel.Attr("href")
	if !existAttr {
		log.Fatal("第" + strconv.Itoa(level) + "层li标签下未发现a标签存在href属性")
	}
	return removeSpace(aSel.Text()), menuLink
}

// 获取li标签下第一个button.select-none的文本
func findButtonText(level int, liSel *goquery.Selection) string {
	buttonSel := liSel.ChildrenFiltered("div").First().ChildrenFiltered("div").First().ChildrenFiltered("button.select-none").First()
	if buttonSel.Length() == 0 {
		log.Fatal("第" + strconv.Itoa(level) + "层li标签下未发现button.select-none标签")
	}
	return removeSpace(buttonSel.Text())
}

var mdHeadStr = `+++
title = "%s"
date = %s
weight = %d
type = "docs"
description = ""
isCJKLanguage = true
draft = false
+++

> 原文: [%s](%s)
>
> 收录该文档的时间：%s
`

func genDate() string {
	//secondsEastOfUTC := int((8 * time.Hour).Seconds())
	//beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	//now := time.Now()
	return time.Now().Format(time.RFC3339)
}

// 生成 Markdown 文件的头部内容
func genMdHeadStr(title, baseUrl, link string, order int) string {
	finalUrl := strings.TrimSuffix(baseUrl, "/") + "/" + strings.TrimPrefix(link, "/")
	weight := 1
	if order != 0 {
		weight = order * 10
	}
	t := genDate()
	return fmt.Sprintf(mdHeadStr, title, t, weight, finalUrl, finalUrl, "`"+t+"`")
}

// 生成Markdown文件，并往该文件中写入头部内容
func createMdFileAndWrite(filePathWithoutExt, headContent string) {
	file, err := os.Create(filePathWithoutExt + ".md")
	if err != nil {
		log.Fatal("创建文件：" + filePathWithoutExt + ".md" + " 失败")
	}
	defer file.Close()
	_, err = file.WriteString(headContent)
	if err != nil {
		log.Fatal("写入文件：" + filePathWithoutExt + ".md" + " 失败")
	}
}

// 先创建菜单对应的目录文件，再在该目录下生成_index.md的Markdown文件，并往该文件中写入头部内容
func createDirAndMdFileAndWrite(newDirPath, fileNameWithoutExt, headContent string) {
	// 创建目录
	err0 := os.Mkdir(newDirPath, 0755)
	if err0 != nil {
		log.Fatal("创建目录：" + newDirPath + " 失败")
	}

	createMdFileAndWrite(withSlashJoinStr(newDirPath, fileNameWithoutExt), headContent)
}

// 创建HTML文件，（目前不会使用到，故直接return）
func createHtml(newFilePathWithoutExt string) {
	return
}

// 通过创建的http客户端请求对应url的HTML内容，并将请求到的HTML内容写入到新创建的HTML文件中 目前不会使用到，故直接return
func httpGetContentAndWriteToHtmlFile(cmd *cobra.Command, link, newFilePathWithoutExt string) {
	return
}

func withSlashJoinStr(sl ...string) string {
	str := ""
	for i, iSl := range sl {
		if i > 0 {
			str = str + "/" + iSl
		} else {
			str = iSl
		}
	}
	return str
}

// 判断是否存在class="select-none"的元素 ？
func judgeExistClassEqSelectNone(liSel *goquery.Selection) bool {
	el := liSel.ChildrenFiltered("div").First().ChildrenFiltered("div").First().ChildrenFiltered(".select-none").First()
	return el.Length() > 0
}

// 判断class="select-none"的元素是button、div还是a标签
// 1：button 2：a
func judgeEleIsDivOrA(liSel *goquery.Selection) int {
	if liSel.ChildrenFiltered("div").First().ChildrenFiltered("div").First().ChildrenFiltered(".select-none").First().Is("button") {
		return 1
	}
	if liSel.ChildrenFiltered("div").First().ChildrenFiltered("div").First().ChildrenFiltered(".select-none").First().Is("a") {
		return 2
	}
	panic("class=\"select-none\"的元素不是button、a标签")
}

// 获取菜单文本和链接
func getMenuTextAndLink(level, liIndex int, liSel *goquery.Selection) (string, string) {
	_ = liIndex
	menuText, menuLink := "", ""
	// 判断有class="select-none"的元素是 button 还是 div？
	if judgeExistClassEqSelectNone(liSel) {
		// 判断是 button 还是 a
		bda := judgeEleIsDivOrA(liSel)
		if bda == 1 { // 是button
			//fmt.Println("1")
			menuText = findButtonText(level, liSel)
		} else if bda == 2 { // 是a
			//fmt.Println("2")
			menuText, menuLink = findATextAndHref(level, liSel)
		}
	} else { // 没有class="select-none"的元素的情况， 默认情况下不出现这种情况
		menuText, menuLink = findATextAndHref(level, liSel)
	}
	return menuText, menuLink
}

func Html2md(cmd *cobra.Command, args []string) {
	url, _ := cmd.Flags().GetString("opt-url")
	selector, _ := cmd.Flags().GetString("opt-nav-selector")
	menuname, _ := cmd.Flags().GetString("menuname")
	dist, _ := cmd.Flags().GetString("dist")
	baseUrl, _ := cmd.Flags().GetString("base-url")

	client := http.Client{
		Timeout: 500 * time.Second, // 设置超时时间为5秒
	}

	fmt.Println(url)
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	fmt.Println(doc.Html())

	if err != nil {
		resp.Body.Close()
		log.Fatal(err)
	}
	resp.Body.Close()

	_, err = os.Stat(dist)
	if os.IsNotExist(err) {
		err = os.Mkdir(dist, 0755)
		if err != nil {
			log.Fatal("创建目录：" + dist + " 失败")
		}
	}

	_, err = os.Stat(filepath.Join(dist, menuname))
	if os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(dist, menuname), 0755)
		if err != nil {
			log.Fatal("创建目录：" + filepath.Join(dist, menuname) + " 失败")
		}
	}

	dist = filepath.Join(dist, menuname)
	fmt.Println("dist=", dist)

	menuText0, menuLink0 := "", ""
	ulSelXs := doc.Find(selector).ChildrenFiltered("ul")
	if ulSelXs.Length() != 0 {
		fmt.Println("获取html内容中。。。")
		ulSelXs.Each(func(ulIndex0 int, ulSel0 *goquery.Selection) {
			liSelX0s := ulSel0.ChildrenFiltered("li")
			fmt.Println("存在第一级菜单个数：", liSelX0s.Length())
			liSelX0s.Each(func(liIndex0 int, liSel0 *goquery.Selection) {
				menuText0, menuLink0 = findATextAndHref(0, liSel0)
				fileName0 := fileNameFromMenuLink(menuLink0)
				createMdFileAndWrite(withSlashJoinStr(dist, fileName0), genMdHeadStr(menuText0, baseUrl, menuLink0, ulIndex0*1000+liIndex0))
				// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
				httpGetContentAndWriteToHtmlFile(cmd, menuLink0, withSlashJoinStr(dist, fileName0))
			})
		})
	}
	fmt.Println("已全部获取所有html的内容！")
}
