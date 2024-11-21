package mysql

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
		`❤`,
	}

	for _, pattern := range patternSlice {
		re, _ := regexp.Compile(pattern)
		str = re.ReplaceAllString(str, "")
	}

	re3, _ := regexp.Compile(`\s+`)
	str = re3.ReplaceAllString(str, "")

	return str
}

func findATextAndHref(level int, liSel *goquery.Selection) (string, string) {
	aSel := liSel.ChildrenFiltered("div").ChildrenFiltered("div.docs-sidebar-nav-link").First().Find("a")
	if aSel.Length() == 0 {
		log.Fatal("a 第" + strconv.Itoa(level) + "层li标签下未发现a标签")
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
	err0 := os.Mkdir(newDirPath, 0755)
	if err0 != nil {
		log.Fatal("创建目录：" + newDirPath + " 失败")
	}

	createMdFileAndWrite(withSlashJoinStr(newDirPath, fileNameWithoutExt), headContent)
}

func httpGetContentAndWriteToHtmlFile(cmd *cobra.Command, link, newFilePathWithoutExt string) {
	selector, _ := cmd.Flags().GetString("opt-content-selector")
	baseUrl, _ := cmd.Flags().GetString("base-url")
	url := link
	re, err := regexp.Compile(`^[(http\:\/\/)|(https:\/\/)]]`)
	if !re.MatchString(link) {
		url = strings.TrimSuffix(baseUrl, "/") + "/" + strings.TrimPrefix(link, "/")
	}

	file, err := os.Create(newFilePathWithoutExt + ".html")
	if err != nil {
		log.Fatal("创建文件：" + newFilePathWithoutExt + ".html" + " 失败")
	}

	time.Sleep(time.Second)

	client := http.Client{
		Timeout: 500 * time.Second, // 设置超时时间为5秒
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal("获取链接："+url+" 的html内容出现错误", err)
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	html, _ := doc.Find(selector).Html()
	file.WriteString(html)
	defer file.Close()
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
		err = os.Mkdir(dist, 0755)
		if err != nil {
			log.Fatal("创建目录：" + dist + " 失败")
		}
	}

	liSelX0s := doc.Find(selector).ChildrenFiltered("ul").ChildrenFiltered("li")

	menuText0, menuLink0, dirName0 := "", "", ""
	menuText1, menuLink1, dirName1 := "", "", ""
	menuText2, menuLink2, dirName2 := "", "", ""
	menuText3, menuLink3, dirName3 := "", "", ""
	menuText4, menuLink4, dirName4 := "", "", ""
	menuText5, menuLink5 := "", ""

	fmt.Println("获取html内容中。。。")
	liSelX0s.Each(func(liIndex0 int, liSel0 *goquery.Selection) {
		ulSelX1s := liSel0.ChildrenFiltered("div").ChildrenFiltered("ul")
		// 若第一层菜单下没有找到 ul 标签，则说明该菜单下没有子菜单
		if ulSelX1s.Length() == 0 {
			menuText0, menuLink0 = findATextAndHref(0, liSel0)
			fileName0 := fileNameFromMenuText(menuText0)
			createMdFileAndWrite(withSlashJoinStr(dist, fileName0), genNonIndexHeadStr(menuText0, baseUrl, menuLink0, liIndex0))

			// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
			httpGetContentAndWriteToHtmlFile(cmd, menuLink0, withSlashJoinStr(dist, fileName0))
			//fmt.Println("第0层", menuText0, " -> ", menuLink0)
		} else {
			menuText0, menuLink0 = findATextAndHref(0, liSel0)
			dirName0 = fileNameFromMenuText(menuText0)
			fileName0 := "_index"
			createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0), fileName0, genIndexHeadStr(menuText0, baseUrl, menuLink0, liIndex0))

			// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
			httpGetContentAndWriteToHtmlFile(cmd, menuLink0, withSlashJoinStr(dist, dirName0, fileName0))

			//fmt.Println("第0层", menuText0, " -> ", menuLink0)

			// 第一层菜单下有子菜单
			ulSelX1s.Each(func(ulIndex1 int, ulSel1 *goquery.Selection) {
				liSelX1s := ulSel1.ChildrenFiltered("li")

				liSelX1s.Each(func(liIndex1 int, liSel1 *goquery.Selection) {
					// ------------------------------------------------------------------------------------
					ulSelX2s := liSel1.ChildrenFiltered("div").ChildrenFiltered("ul")
					// ------------------------------------------------------------------------------------
					if ulSelX2s.Length() == 0 {
						menuText1, menuLink1 = findATextAndHref(1, liSel1)
						fileName1 := fileNameFromMenuText(menuText1)
						createMdFileAndWrite(withSlashJoinStr(dist, dirName0, fileName1), genNonIndexHeadStr(menuText1, baseUrl, menuLink1, liIndex1))

						// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
						httpGetContentAndWriteToHtmlFile(cmd, menuLink1, withSlashJoinStr(dist, dirName0, fileName1))

						//fmt.Println("第1层", menuText1, " -> ", menuLink1)
					} else {
						menuText1, menuLink1 = findATextAndHref(1, liSel1)
						dirName1 = fileNameFromMenuText(menuText1)
						fileName1 := "_index"
						createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1), fileName1, genNonIndexHeadStr(menuText1, baseUrl, menuLink1, liIndex1))

						// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
						httpGetContentAndWriteToHtmlFile(cmd, menuLink1, withSlashJoinStr(dist, dirName0, dirName1, fileName1))

						//fmt.Println("第1层", menuText1, " -> ", menuLink1)

						ulSelX2s.Each(func(ulIndex2 int, ulSel2 *goquery.Selection) {
							liSelX2s := ulSel2.ChildrenFiltered("li")

							liSelX2s.Each(func(liIndex2 int, liSel2 *goquery.Selection) {
								// ------------------------------------------------------------------------------------
								ulSelX3s := liSel2.ChildrenFiltered("div").ChildrenFiltered("ul")
								// ------------------------------------------------------------------------------------
								if ulSelX3s.Length() == 0 {
									menuText2, menuLink2 = findATextAndHref(2, liSel2)
									fileName2 := fileNameFromMenuText(menuText2)
									createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, fileName2), genNonIndexHeadStr(menuText2, baseUrl, menuLink2, liIndex2))

									// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
									httpGetContentAndWriteToHtmlFile(cmd, menuLink2, withSlashJoinStr(dist, dirName0, dirName1, fileName2))

									//fmt.Println("第2层", menuText2, " -> ", menuLink2)
								} else {
									menuText2, menuLink2 = findATextAndHref(2, liSel2)
									dirName2 = fileNameFromMenuText(menuText2)
									fileName2 := "_index"
									createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2), fileName2, genNonIndexHeadStr(menuText2, baseUrl, menuLink2, liIndex2))

									// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
									httpGetContentAndWriteToHtmlFile(cmd, menuLink2, withSlashJoinStr(dist, dirName0, dirName1, dirName2, fileName2))

									//fmt.Println("第2层", menuText2, " -> ", menuLink2)

									ulSelX3s.Each(func(ulIndex3 int, ulSel3 *goquery.Selection) {
										liSelX3s := ulSel3.ChildrenFiltered("li")

										liSelX3s.Each(func(liIndex3 int, liSel3 *goquery.Selection) {
											// ------------------------------------------------------------------------------------
											ulSelX4s := liSel3.ChildrenFiltered("div").ChildrenFiltered("ul")
											// ------------------------------------------------------------------------------------
											if ulSelX4s.Length() == 0 {
												menuText3, menuLink3 = findATextAndHref(3, liSel3)
												fileName3 := fileNameFromMenuText(menuText3)
												createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, fileName3), genNonIndexHeadStr(menuText3, baseUrl, menuLink3, liIndex3))

												// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
												httpGetContentAndWriteToHtmlFile(cmd, menuLink3, withSlashJoinStr(dist, dirName0, dirName1, dirName2, fileName3))

												//fmt.Println("第3层", menuText3, " -> ", menuLink3)
											} else {
												menuText3, menuLink3 = findATextAndHref(3, liSel3)

												dirName3 = fileNameFromMenuText(menuText3)
												fileName3 := "_index"
												createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3), fileName3, genNonIndexHeadStr(menuText3, baseUrl, menuLink3, liIndex3))

												// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
												httpGetContentAndWriteToHtmlFile(cmd, menuLink3, withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, fileName3))

												fmt.Println("第3层", menuText3, " -> ", menuLink3)

												ulSelX4s.Each(func(ulIndex4 int, ulSel4 *goquery.Selection) {
													liSelX4s := ulSel4.ChildrenFiltered("li")

													liSelX4s.Each(func(liIndex4 int, liSel4 *goquery.Selection) {
														// ------------------------------------------------------------------------------------
														ulSelX5s := liSel4.ChildrenFiltered("div").ChildrenFiltered("ul")
														// ------------------------------------------------------------------------------------
														if ulSelX5s.Length() == 0 {
															menuText4, menuLink4 = findATextAndHref(4, liSel4)
															fileName4 := fileNameFromMenuText(menuText4)
															createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, fileName4), genNonIndexHeadStr(menuText4, baseUrl, menuLink4, liIndex4))

															// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
															httpGetContentAndWriteToHtmlFile(cmd, menuLink4, withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, fileName4))

															//fmt.Println("第4层", menuText3, " -> ", menuLink3)
														} else {
															menuText4, menuLink4 = findATextAndHref(4, liSel4)
															dirName4 = fileNameFromMenuText(menuText4)
															fileName4 := "_index"
															createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4), fileName4, genNonIndexHeadStr(menuText4, baseUrl, menuLink4, liIndex4))

															// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
															httpGetContentAndWriteToHtmlFile(cmd, menuLink4, withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, fileName4))

															//fmt.Println("第4层", menuText4, " -> ", menuLink4)

															ulSelX5s.Each(func(ulIndex5 int, ulSel5 *goquery.Selection) {
																liSelX5s := ulSel5.ChildrenFiltered("li")

																liSelX5s.Each(func(liIndex5 int, liSel5 *goquery.Selection) {
																	// ------------------------------------------------------------------------------------
																	ulSelX6s := liSel5.ChildrenFiltered("div").ChildrenFiltered("ul")
																	// ------------------------------------------------------------------------------------
																	if ulSelX6s.Length() == 0 {
																		menuText5, menuLink5 = findATextAndHref(5, liSel5)
																		fileName5 := fileNameFromMenuText(menuText5)
																		createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, fileName5), genNonIndexHeadStr(menuText5, baseUrl, menuLink5, liIndex5))

																		// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
																		httpGetContentAndWriteToHtmlFile(cmd, menuLink5, withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, fileName5))

																		//fmt.Println("第5层", menuText3, " -> ", menuLink3)
																	} else {
																		fmt.Println("当前只支持5层菜单")
																		//log.Fatal("当前只支持5层菜单")
																		//menuText5, menuLink5 = findATextAndHref(5, liSel5)
																		//fmt.Println("第5层", menuText5, " -> ", menuLink5)

																	}
																})
															})

														}
													})
												})

											}
										})
									})

								}
							})
						})
					}
				})
			})
		}
	})

	//fmt.Println("获取html内容中。。。")
	fmt.Println("已全部获取所有html的内容！")
}
