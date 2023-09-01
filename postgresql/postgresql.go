package postgresql

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
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
		`!`, `！`, `=`, `@`,
		`\$`,
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

func findATextAndHref(level int, dtSel *goquery.Selection) (string, string) {
	aSel := dtSel.Find("a")
	if aSel.Length() == 0 {
		log.Fatal("第" + strconv.Itoa(level) + "层dt标签下未发现a标签")
	}
	menuText := aSel.Text()

	re1, _ := regexp.Compile(`\n`)
	if re1.MatchString(menuText) {
		menuText = strings.Replace(menuText, "\n", "", -1)
	}

	re2, _ := regexp.Compile(`\s+`)
	menuText = re2.ReplaceAllString(menuText, " ")
	menuText = strings.TrimSpace(menuText)

	menuLink, existAttr := aSel.Attr("href")
	if !existAttr {
		log.Fatal("b 第" + strconv.Itoa(level) + "层li标签下未发现a标签存在href属性")
	}
	if menuLink == "" {
		fmt.Println("menuLink = ", menuLink, ", existAttr = ", existAttr)
		fmt.Println(aSel.Html())
	}

	return menuText, menuLink
}

var indexHeadStr = `+++
title = "%s"
linkTitle = "%s"
date = %s
type = "docs"
description = ""
isCJKLanguage = true
draft = false
[menu.main]
    weight = %d
+++

> 原文: [%s](%s)
`

var nonIndexHeadStr = `+++
title = "%s"
date = %s
weight = %d
type = "docs"
description = ""
isCJKLanguage = true
draft = false
+++

> 原文: [%s](%s)
`

func genDate() string {
	//secondsEastOfUTC := int((8 * time.Hour).Seconds())
	//beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	//now := time.Now()
	return time.Now().Format("2006-01-02T15:04:05+08:00")
}

func genIndexHeadStr(title, baseUrl, link string, order int) string {
	finalUrl := strings.TrimSuffix(baseUrl, "/") + "/" + strings.TrimPrefix(link, "/")
	weight := 1
	if order != 0 {
		weight = order * 10
	}

	return fmt.Sprintf(indexHeadStr, title, title, genDate(), weight, finalUrl, finalUrl)
}

func genNonIndexHeadStr(title, baseUrl, link string, order int) string {
	finalUrl := strings.TrimSuffix(baseUrl, "/") + "/" + strings.TrimPrefix(link, "/")
	weight := 1
	if order != 0 {
		weight = order * 10
	}

	return fmt.Sprintf(nonIndexHeadStr, title, genDate(), weight, finalUrl, finalUrl)
}

func createMdFileAndWrite(filePathWithoutExt, headContent string) {
	file0, err0 := os.Create(filePathWithoutExt + ".md")
	if err0 != nil {
		log.Fatal("创建文件：" + filePathWithoutExt + ".md" + " 失败")
	}

	file0.WriteString(headContent)
	file0.Close()
}

func createDirAndMdFileAndWrite(newDirPath, fileNameWithoutExt, headContent string) {
	// 创建目录
	err0 := os.Mkdir(newDirPath, 755)
	if err0 != nil {
		log.Fatal("创建目录：" + newDirPath + " 失败")
	}

	createMdFileAndWrite(withSlashJoinStr(newDirPath, fileNameWithoutExt), headContent)
}

func HttpGetContent(url string) (resp *http.Response) {
	time.Sleep(time.Second)

	client := http.Client{
		Timeout: 500 * time.Second, // 设置超时时间为5秒
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal("获取链接："+url+" 的html内容出现错误", err)
	}
	defer resp.Body.Close()
	return resp
}

func HandleUrl(url, baseUrl string) string {
	re, err := regexp.Compile(`^[(http\:\/\/)|(https:\/\/)]]`)
	if err != nil {
		log.Fatal("在HandleUrl函数创建正则匹配时发生错误：", err)
	}
	if !re.MatchString(url) {
		url = strings.TrimSuffix(baseUrl, "/") + "/" + url
	}
	return url
}

func httpGetContentAndWriteToHtmlFile(cmd *cobra.Command, link, newFilePathWithoutExt string) {
	selector, _ := cmd.Flags().GetString("opt-content-selector")
	baseUrl, _ := cmd.Flags().GetString("base-url")
	url := HandleUrl(link, baseUrl)
	file, err := os.Create(newFilePathWithoutExt + ".html")
	if err != nil {
		log.Fatal("创建文件：" + newFilePathWithoutExt + ".html" + " 失败")
	}
	defer file.Close()

	client := http.Client{
		Timeout: 500 * time.Second, // 设置超时时间为5秒
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal("获取链接："+url+" 的html内容出现错误", err)
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	// 去除不必要的内容
	doc.Find("div.navheader").Remove()
	doc.Find("div.toc").Remove()
	doc.Find("a.indexterm").Remove()
	doc.Find("div.navfooter").Remove()

	html, err := doc.Find(selector).Html()
	if err != nil {
		log.Fatal("往html文件中写入内容发生错误：", err)
	}
	file.WriteString(html)
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

func Html2md(cmd *cobra.Command, args []string) {
	url, _ := cmd.Flags().GetString("opt-url")
	selector, _ := cmd.Flags().GetString("opt-nav-selector")
	dist, _ := cmd.Flags().GetString("dist")
	baseUrl, _ := cmd.Flags().GetString("base-url")

	client := http.Client{
		Timeout: 500 * time.Second, // 设置超时时间为5秒
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		resp.Body.Close()
		log.Fatal(err)
	}
	resp.Body.Close()

	_, err = os.Stat(dist)
	if !os.IsExist(err) {
		err = os.Mkdir(dist, 755)
		if err != nil {
			log.Fatal("创建目录：" + dist + " 失败")
		}
	}

	dtSelX0s := doc.Find(selector).ChildrenFiltered("dl").ChildrenFiltered("dt")

	menuText0, menuLink0, dirName0 := "", "", ""
	menuText1, menuLink1, dirName1 := "", "", ""
	menuText2, menuLink2 := "", ""

	fmt.Println("获取html内容中。。。")
	dtSelX0s.Each(func(dtIndex0 int, dtSel0 *goquery.Selection) {
		fmt.Println(dtIndex0, "---------------------------------")
		fmt.Println(dtSel0.Find("a").Text())

		dtSelX1s := dtSel0.Next().ChildrenFiltered("dl").ChildrenFiltered("dt")
		if dtSelX1s.Length() == 0 {
			menuText0, menuLink0 = findATextAndHref(1, dtSel0)
			fileName0 := fileNameFromMenuText(menuText0)
			createMdFileAndWrite(withSlashJoinStr(dist, fileName0), genNonIndexHeadStr(menuText0, baseUrl, menuLink0, dtIndex0))

			// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
			httpGetContentAndWriteToHtmlFile(cmd, menuLink0, withSlashJoinStr(dist, fileName0))
			fmt.Println("run here2")
		} else { // 存在子菜单
			menuText0, menuLink0 = findATextAndHref(0, dtSel0)
			dirName0 = fileNameFromMenuText(menuText0)
			fileName0 := "_index"
			createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0), fileName0, genIndexHeadStr(menuText0, baseUrl, menuLink0, dtIndex0))

			// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
			httpGetContentAndWriteToHtmlFile(cmd, menuLink0, withSlashJoinStr(dist, dirName0, fileName0))

			// 以下获取 子菜单中的内容
			dtSelX1s.Each(func(dtIndex1 int, dtSel1 *goquery.Selection) {
				// 先获取该子菜单的html内容， 再更据子菜单中是否有更进一步的子菜单，分情况进行创建目录和文件
				menuText1, menuLink1 = findATextAndHref(1, dtSel1)
				client := http.Client{
					Timeout: 500 * time.Second, // 设置超时时间为5秒
				}

				menuResp1, err := client.Get(HandleUrl(menuLink1, baseUrl))
				if err != nil {
					log.Fatal("获取链接："+url+" 的html内容出现错误", err)
				}
				defer menuResp1.Body.Close()

				menuDoc1, _ := goquery.NewDocumentFromReader(menuResp1.Body)

				ddSelX2s := menuDoc1.Find("div.toc").ChildrenFiltered("dl").ChildrenFiltered("dd")
				if ddSelX2s.Length() == 0 { // 可能没有子菜单的情况
					// 进一步判断是否当前页面的内容不足的问题
					// 若不存在 div.sect2 但却有 dt 的情况，则也是要当成有子菜单
					if menuDoc1.Find("div.toc").ChildrenFiltered("dl").ChildrenFiltered("dt").Length() != 0 && menuDoc1.Find("div.sect2").Length() == 0 {
						// 1 先创建 menuText1 子菜单文件夹
						// 2 将 menuDoc1 中的内容放入 _index.html 文件中 以及 创建 _index.md 文件
						dirName1 = fileNameFromMenuText(menuText1)
						fileName1 := "_index"
						createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1), fileName1, genNonIndexHeadStr(menuText1, baseUrl, menuLink1, dtIndex1))
						// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
						httpGetContentAndWriteToHtmlFile(cmd, menuLink1, withSlashJoinStr(dist, dirName0, dirName1, fileName1))

						dtSelX2s := menuDoc1.Find("div.toc").ChildrenFiltered("dl").ChildrenFiltered("dt")
						// 3 在将 dtSelX2s 下的所有对应子菜单的内容 放入各自的.html 和 .md 文件
						dtSelX2s.Each(func(dtIndex2 int, dtSel2 *goquery.Selection) {
							menuText2, menuLink2 = findATextAndHref(2, dtSel2)
							fileName2 := fileNameFromMenuText(menuText2)
							createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, fileName2), genNonIndexHeadStr(menuText2, baseUrl, menuLink2, dtIndex2))
							// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
							httpGetContentAndWriteToHtmlFile(cmd, menuLink2, withSlashJoinStr(dist, dirName0, dirName1, fileName2))
						})
					} else {
						fileName1 := fileNameFromMenuText(menuText1)
						createMdFileAndWrite(withSlashJoinStr(dist, dirName0, fileName1), genNonIndexHeadStr(menuText1, baseUrl, menuLink1, dtIndex1))
						// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
						httpGetContentAndWriteToHtmlFile(cmd, menuLink1, withSlashJoinStr(dist, dirName0, fileName1))
					}
				} else { // 有子菜单的情况
					// 1 先创建 menuText1 子菜单文件夹
					// 2 将 menuDoc1 中的内容放入 _index.html 文件中 以及 创建 _index.md 文件
					dirName1 = fileNameFromMenuText(menuText1)
					fileName1 := "_index"
					createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1), fileName1, genNonIndexHeadStr(menuText1, baseUrl, menuLink1, dtIndex1))
					// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
					httpGetContentAndWriteToHtmlFile(cmd, menuLink1, withSlashJoinStr(dist, dirName0, dirName1, fileName1))

					dtSelX2s := menuDoc1.Find("div.toc").ChildrenFiltered("dl").ChildrenFiltered("dt")
					// 3 在将 dtSelX2s 下的所有对应子菜单的内容 放入各自的.html 和 .md 文件
					dtSelX2s.Each(func(dtIndex2 int, dtSel2 *goquery.Selection) {
						menuText2, menuLink2 = findATextAndHref(2, dtSel2)
						fileName2 := fileNameFromMenuText(menuText2)
						createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, fileName2), genNonIndexHeadStr(menuText2, baseUrl, menuLink2, dtIndex2))
						// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
						httpGetContentAndWriteToHtmlFile(cmd, menuLink2, withSlashJoinStr(dist, dirName0, dirName1, fileName2))
					})
				}
			})
		}
	})

	//fmt.Println("获取html内容中。。。")
	fmt.Println("已全部获取所有html的内容！")
}
