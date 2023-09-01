package docker

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

func findSpanATextAndHref(level int, liSel *goquery.Selection) (string, string) {
	aSel := liSel.ChildrenFiltered("span").First().Find("a")
	if aSel.Length() == 0 {
		log.Fatal("第" + strconv.Itoa(level) + "层li标签下未发现span>a标签")
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
		log.Fatal("第" + strconv.Itoa(level) + "层li标签下未发现a标签存在href属性")
	}
	if menuLink == "" {
		fmt.Println("menuLink = ", menuLink, ", existAttr = ", existAttr)
		fmt.Println(aSel.Html())
	}

	return menuText, menuLink
}

func findSpanText(level int, liSel *goquery.Selection) string {
	spanSel := liSel.ChildrenFiltered("button").ChildrenFiltered("span").First()
	if spanSel.Length() == 0 {
		log.Fatal("第" + strconv.Itoa(level) + "层li标签下未发现button>span标签")
	}
	menuText := spanSel.Text()

	re1, _ := regexp.Compile(`\n`)
	if re1.MatchString(menuText) {
		menuText = strings.Replace(menuText, "\n", "", -1)
	}
	re2, _ := regexp.Compile(`\s+`)
	menuText = re2.ReplaceAllString(menuText, " ")
	menuText = strings.TrimSpace(menuText)

	return menuText
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
	finalUrl := strings.TrimSuffix(baseUrl, "/") + "/" + link
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

func createHtml(newFilePathWithoutExt string) {
	file, err := os.Create(newFilePathWithoutExt + ".html")
	if err != nil {
		log.Fatal("创建文件：" + newFilePathWithoutExt + ".html" + " 失败")
	}
	defer file.Close()
}

func httpGetContentAndWriteToHtmlFile(cmd *cobra.Command, link, newFilePathWithoutExt string) {
	selector, _ := cmd.Flags().GetString("opt-content-selector")
	baseUrl, _ := cmd.Flags().GetString("base-url")
	//urlpath := "/" + strings.TrimPrefix(link, "/")
	url := link
	re, err := regexp.Compile(`^[(http\:\/\/)|(https:\/\/)]]`)
	if !re.MatchString(link) {
		url = strings.TrimSuffix(baseUrl, "/") + "/" + strings.TrimPrefix(link, "/")
	}

	file, err := os.Create(newFilePathWithoutExt + ".html")
	if err != nil {
		log.Fatal("创建文件：" + newFilePathWithoutExt + ".html" + " 失败")
	}
	defer file.Close()

	time.Sleep(5 * time.Second)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	//req.Header.Set(":authority", "docs.docker.com")
	//req.Header.Set(":method", "GET")
	//req.Header.Set(":path", urlpath)
	//req.Header.Set(":scheme", "http")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Cookie", "")
	req.Header.Set("Pragma", "no-cache")
	//req.Header.Set("Sec-Ch-Ua", `"Chromium";v="116", "Not)A;Brand";v="24", "Google Chrome";v="116"`)
	//req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	//req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	//req.Header.Set("Sec-Fetch-Dest", "document")
	//req.Header.Set("Sec-Fetch-Mode", "navigate")
	//req.Header.Set("Sec-Fetch-Site", "none")
	//req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	client := &http.Client{
		Timeout: 500 * time.Second, // 设置超时时间为5秒
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("获取链接："+url+" 的html内容出现错误", err)
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	//doc.Find("#subnav").Remove()
	//doc.Find("#page-github-links").Remove()
	//doc.Find("#subnav").Remove()
	//doc.Find("#toc").Remove()

	// 处理 a href 和 img src
	doc.Find("a").Each(func(aIndex int, aSel *goquery.Selection) {
		href, isExist := aSel.Attr("href")
		if isExist {
			// 注意先去除最后的 /
			trueUrl := strings.TrimSuffix(url, "/")
			if strings.HasPrefix(href, "/") {
				href = strings.TrimSuffix(baseUrl, "/") + href
				aSel.SetAttr("href", href)
			} else {
				if strings.HasPrefix(href, ".././") {
					href = strings.Replace(href, ".././", "../", -1)
				}

				if strings.HasPrefix(href, "../") {
					href = strings.Replace(href, "../", "", -1)

					// 去除 trueUrl 中有 ? 和 #
					index := strings.LastIndex(trueUrl, "?")
					if index != -1 {
						trueUrl = trueUrl[:index]
					}

					index = strings.LastIndex(trueUrl, "#")
					if index != -1 {
						trueUrl = trueUrl[:index]
					}

					trueUrl = strings.TrimSuffix(trueUrl, "/")
					index = strings.LastIndex(trueUrl, "/")
					if index != -1 {
						trueUrl = trueUrl[:(index + 1)]
					}
					href = trueUrl + href
					aSel.SetAttr("href", href)
				}
			}
		}
	})

	doc.Find("img").Each(func(imgIndex int, imgSel *goquery.Selection) {
		src, isExist := imgSel.Attr("src")
		if isExist {
			// 注意先去除最后的 /
			trueUrl := strings.TrimSuffix(url, "/")
			if strings.HasPrefix(src, "/") {
				src = strings.TrimSuffix(baseUrl, "/") + src
				imgSel.SetAttr("src", src)
			} else {
				if strings.HasPrefix(src, ".././") {
					src = strings.Replace(src, ".././", "../", -1)
				}

				// 去除 trueUrl 中有 ? 和 #
				index := strings.LastIndex(trueUrl, "?")
				if index != -1 {
					trueUrl = trueUrl[:index]
				}

				index = strings.LastIndex(trueUrl, "#")
				if index != -1 {
					trueUrl = trueUrl[:index]
				}

				for strings.HasPrefix(src, "../") {
					src = strings.Replace(src, "../", "", 1)

					trueUrl = strings.TrimSuffix(trueUrl, "/")
					index = strings.LastIndex(trueUrl, "/")
					if index != -1 {
						trueUrl = trueUrl[:(index + 1)]
					}
				}
				src = trueUrl + src
				imgSel.SetAttr("src", src)
			}
		}
	})

	html, _ := doc.Find(selector).Html()

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
	menuname, _ := cmd.Flags().GetString("menuname")
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
	if os.IsNotExist(err) {
		err = os.Mkdir(dist, 755)
		if err != nil {
			log.Fatal("创建目录：" + dist + " 失败")
		}
	}

	_, err = os.Stat(filepath.Join(dist, menuname))
	if os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(dist, menuname), 755)
		if err != nil {
			log.Fatal("创建目录：" + filepath.Join(dist, menuname) + " 失败")
		}
	}

	dist = filepath.Join(dist, menuname)
	fmt.Println("dist=", dist)

	liSelX0s := doc.Find(selector).ChildrenFiltered("ul").ChildrenFiltered("li")

	menuText0, menuLink0, dirName0 := "", "", ""
	menuText1, menuLink1, dirName1 := "", "", ""
	menuText2, menuLink2, dirName2 := "", "", ""
	menuText3, menuLink3, dirName3 := "", "", ""
	menuText4, menuLink4, dirName4 := "", "", ""
	menuText5, menuLink5, dirName5 := "", "", ""
	menuText6, menuLink6 := "", ""

	fmt.Println("获取html内容中。。。")
	liSelX0s.Each(func(liIndex0 int, liSel0 *goquery.Selection) {
		// 判断li的子元素中是否有 button 标签
		// 若有 button 标签，则需创建一个子目录，并在其中创建 _index.md（为了对等也可创建_index.html）
		if liSel0.ChildrenFiltered("button").Length() == 0 {
			menuText0, menuLink0 = findSpanATextAndHref(0, liSel0)
			fileName0 := fileNameFromMenuText(menuText0)
			createMdFileAndWrite(withSlashJoinStr(dist, fileName0), genNonIndexHeadStr(menuText0, baseUrl, menuLink0, liIndex0))
			// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
			httpGetContentAndWriteToHtmlFile(cmd, menuLink0, withSlashJoinStr(dist, fileName0))
		} else {
			menuText0 = findSpanText(0, liSel0)
			dirName0 = fileNameFromMenuText(menuText0)
			fileName0 := "_index"
			createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0), fileName0, genNonIndexHeadStr(menuText0, "", "", liIndex0))
			// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
			createHtml(withSlashJoinStr(dist, dirName0, fileName0))

			liSelX1s := liSel0.ChildrenFiltered("ul").ChildrenFiltered("li")
			liSelX1s.Each(func(liIndex1 int, liSel1 *goquery.Selection) {
				if liSel1.ChildrenFiltered("button").Length() == 0 {
					menuText1, menuLink1 = findSpanATextAndHref(1, liSel1)
					fileName1 := fileNameFromMenuText(menuText1)
					createMdFileAndWrite(withSlashJoinStr(dist, dirName0, fileName1), genNonIndexHeadStr(menuText1, baseUrl, menuLink1, liIndex1))
					// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
					httpGetContentAndWriteToHtmlFile(cmd, menuLink1, withSlashJoinStr(dist, dirName0, fileName1))
				} else {
					menuText1 = findSpanText(1, liSel1)
					dirName1 = fileNameFromMenuText(menuText1)
					fileName1 := "_index"
					createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1), fileName1, genNonIndexHeadStr(menuText1, "", "", liIndex1))
					// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
					createHtml(withSlashJoinStr(dist, dirName0, dirName1, fileName1))

					liSelX2s := liSel1.ChildrenFiltered("ul").ChildrenFiltered("li")
					liSelX2s.Each(func(liIndex2 int, liSel2 *goquery.Selection) {

						if liSel2.ChildrenFiltered("button").Length() == 0 {
							menuText2, menuLink2 = findSpanATextAndHref(2, liSel2)
							fileName2 := fileNameFromMenuText(menuText2)
							createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, fileName2), genNonIndexHeadStr(menuText2, baseUrl, menuLink2, liIndex2))
							// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
							httpGetContentAndWriteToHtmlFile(cmd, menuLink2, withSlashJoinStr(dist, dirName0, dirName1, fileName2))
						} else {
							menuText2 = findSpanText(2, liSel2)
							dirName2 = fileNameFromMenuText(menuText2)
							fileName2 := "_index"
							createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2), fileName2, genNonIndexHeadStr(menuText2, "", "", liIndex2))
							// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
							createHtml(withSlashJoinStr(dist, dirName0, dirName1, dirName2, fileName2))

							liSelX3s := liSel2.ChildrenFiltered("ul").ChildrenFiltered("li")
							liSelX3s.Each(func(liIndex3 int, liSel3 *goquery.Selection) {

								if liSel3.ChildrenFiltered("button").Length() == 0 {
									menuText3, menuLink3 = findSpanATextAndHref(3, liSel3)
									fileName3 := fileNameFromMenuText(menuText3)
									createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, fileName3), genNonIndexHeadStr(menuText3, baseUrl, menuLink3, liIndex3))
									// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
									httpGetContentAndWriteToHtmlFile(cmd, menuLink3, withSlashJoinStr(dist, dirName0, dirName1, dirName2, fileName3))
								} else {
									menuText3 = findSpanText(3, liSel3)
									dirName3 = fileNameFromMenuText(menuText3)
									fileName3 := "_index"
									createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3), fileName3, genNonIndexHeadStr(menuText3, "", "", liIndex3))
									// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
									createHtml(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, fileName3))

									liSelX4s := liSel3.ChildrenFiltered("ul").ChildrenFiltered("li")
									liSelX4s.Each(func(liIndex4 int, liSel4 *goquery.Selection) {

										if liSel4.ChildrenFiltered("button").Length() == 0 {
											menuText4, menuLink4 = findSpanATextAndHref(4, liSel4)
											fileName4 := fileNameFromMenuText(menuText4)
											createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, fileName4), genNonIndexHeadStr(menuText4, baseUrl, menuLink4, liIndex4))
											// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
											httpGetContentAndWriteToHtmlFile(cmd, menuLink4, withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName4, fileName4))
										} else {
											menuText4 = findSpanText(4, liSel4)
											dirName4 = fileNameFromMenuText(menuText4)
											fileName4 := "_index"
											createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4), fileName4, genNonIndexHeadStr(menuText4, "", "", liIndex4))
											// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
											createHtml(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, fileName4))

											liSelX5s := liSel4.ChildrenFiltered("ul").ChildrenFiltered("li")
											liSelX5s.Each(func(liIndex5 int, liSel5 *goquery.Selection) {

												if liSel5.ChildrenFiltered("button").Length() == 0 {
													menuText5, menuLink5 = findSpanATextAndHref(5, liSel5)
													fileName5 := fileNameFromMenuText(menuText5)
													createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, fileName5), genNonIndexHeadStr(menuText5, baseUrl, menuLink5, liIndex5))
													// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
													httpGetContentAndWriteToHtmlFile(cmd, menuLink5, withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, fileName5))
												} else {
													menuText5 = findSpanText(5, liSel5)
													dirName5 = fileNameFromMenuText(menuText5)
													fileName5 := "_index"
													createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, dirName5), fileName5, genNonIndexHeadStr(menuText5, "", "", liIndex5))
													// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
													createHtml(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, dirName5, fileName5))

													liSelX6s := liSel5.ChildrenFiltered("ul").ChildrenFiltered("li")
													liSelX6s.Each(func(liIndex6 int, liSel6 *goquery.Selection) {

														if liSel6.ChildrenFiltered("button").Length() == 0 {
															menuText6, menuLink6 = findSpanATextAndHref(6, liSel6)
															fileName6 := fileNameFromMenuText(menuText6)
															createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, dirName5, fileName6), genNonIndexHeadStr(menuText6, baseUrl, menuLink6, liIndex6))
															// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
															httpGetContentAndWriteToHtmlFile(cmd, menuLink6, withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName4, dirName5, fileName6))
														} else {
															panic("没想到层数竟然达到了6层！")
														}
													})
												}
											})
										}
									})

								}
							})
						}
					})
				}

			})
		}
	})

	//fmt.Println("获取html内容中。。。")
	fmt.Println("已全部获取所有html的内容！")
}
