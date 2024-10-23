package npmjs

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

func findATextAndHref(level int, liSel *goquery.Selection) (string, string) {
	aSel := liSel.Find("a").First()
	if aSel.Length() == 0 {
		log.Fatal("第" + strconv.Itoa(level) + "层li标签下未发现a标签")
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

func findAText(level int, liSel *goquery.Selection) string {
	aSel := liSel.ChildrenFiltered("a").First()
	if aSel.Length() == 0 {
		log.Fatal("第" + strconv.Itoa(level) + "层li标签下未发现a标签")
	}
	menuText := aSel.Text()

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

	time.Sleep(6 * time.Second)

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
	//req.Header.Set("Cache-Control", "no-cache")
	//req.Header.Set("Cookie", "_ga_XJWPQMJYHQ=GS1.2.1693807684.22.1.1693808279.57.0.0; _gcl_au=1.1.147290426.1693613667; _mkto_trk=id:790-SSB-375&token:_mch-docker.com-1693613667387-81538; _fbp=fb.1.1693613667438.603613920; _gid=GA1.2.289936888.1693613670; ajs_anonymous_id=d45e3d0e-f037-4b39-938f-27f8420f86c0; userty.core.p.fe7522=__2VySWQiOiIyZDBlYTZmNjQyNzc3Njc5NmRiODAxNGQ3ZDkzZGQ2NCJ9eyJ1c; _hjSessionUser_3169877=eyJpZCI6IjQ2NWI3YjA0LTkzYjYtNWNmNS1hMzFhLTYwMjJiOWZlMTNmMSIsImNyZWF0ZWQiOjE2OTM2MTM2NjcwODYsImV4aXN0aW5nIjp0cnVlfQ==; OptanonAlertBoxClosed=2023-09-02T00:41:54.828Z; fullstoryStart=false; ln_or=eyIzNzY1MjEwIjoiZCJ9; _hjIncludedInSessionSample_3169877=0; _hjSession_3169877=eyJpZCI6Ijc0NWE2OTZiLWFjNGYtNGRlNS05ODkzLTI0MzU4NmRhOTBhNyIsImNyZWF0ZWQiOjE2OTM4MDc2ODM1NTcsImluU2FtcGxlIjpmYWxzZX0=; _hjAbsoluteSessionInProgress=0; _gali=sidebar; _ga_XJWPQMJYHQ=GS1.1.1693807684.22.1.1693808277.59.0.0; OptanonConsent=isGpcEnabled=0&datestamp=Mon+Sep+04+2023+14%3A17%3A58+GMT%2B0800+(%E4%B8%AD%E5%9B%BD%E6%A0%87%E5%87%86%E6%97%B6%E9%97%B4)&version=202208.1.0&isIABGlobal=false&hosts=&consentId=66597ecd-711d-4a84-91d4-a32f63c1fa8a&interactionCount=1&landingPath=NotLandingPage&groups=C0003%3A1%2CC0001%3A1%2CC0002%3A1%2CC0004%3A1&AwaitingReconsent=false&geolocation=JP%3B13; _hp2_id.4204607514=%7B%22userId%22%3A%22201417274194244%22%2C%22pageviewId%22%3A%226671800441422329%22%2C%22sessionId%22%3A%221277591024353903%22%2C%22identity%22%3Anull%2C%22trackerVersion%22%3A%224.0%22%7D; _ga=GA1.2.258527346.1693613667; _gat=1; _hp2_ses_props.4204607514=%7B%22r%22%3A%22https%3A%2F%2Fdocs.docker.com%2Fdevelop%2Fdevelop-images%2Fdockerfile_best-practices%2F%22%2C%22ts%22%3A1693808274927%2C%22d%22%3A%22docs.docker.com%22%2C%22h%22%3A%22%2Fdevelop%2Fdevelop-images%2Fdockerfile_best-practices%2F%22%7D; userty.core.s.fe7522=__WQiOiI3YjQyYzY5ZjEzODgwNDQxNDY2ODMwMTg2ZmNmNDU3YiIsInN0IjoxNjkzODA3Njk0OTkzLCJyZWFkeSI6dHJ1ZSwid3MiOiJ7XCJ3XCI6MTM1NCxcImhcIjoxNDgwfSIsInNlIjoxNjkzODEwMDg3NzMzLCJwdiI6MTJ9eyJza")
	//req.Header.Set("Pragma", "no-cache")
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
	doc.Find("div.Box-nv15kw-0.Position-fla1y3-0.layout___StyledPosition-sc-1qhwq3g-4.PBNzC.kzYPsQ.ffFUDf").Remove()
	doc.Find("div.Box-nv15kw-0.BorderBox-oreju5-0.lnsqHZ.eOFYr").Remove()

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

	fmt.Println("doc.Find(selector)", doc.Find(selector))
	fmt.Println("doc.Find(selector).ChildrenFiltered(`div[role=\"listitem\"]`)", doc.Find(selector).ChildrenFiltered(`div[role="listitem"]`))

	liSelX0s := doc.Find(selector).ChildrenFiltered(`div[role="listitem"]`).ChildrenFiltered(`div[display="flex"]`)
	fmt.Println("liSelX0s", liSelX0s)
	menuText0, menuLink0, dirName0 := "", "", ""
	menuText1, menuLink1, dirName1 := "", "", ""
	menuText2, menuLink2, dirName2 := "", "", ""
	menuText3, menuLink3, dirName3 := "", "", ""
	menuText4, menuLink4, dirName4 := "", "", ""
	menuText5, menuLink5, dirName5 := "", "", ""
	menuText6, menuLink6 := "", ""

	fmt.Println("获取html内容中。。。")
	liSelX0s.Each(func(liIndex0 int, liSel0 *goquery.Selection) {
		// 判断li的子元素中是否有 div[role="list"] 标签
		// 若有 div[role="list"] 标签，则需创建一个子目录，并在其中创建 _index.md（为了对等也可创建_index.html）
		if liSel0.ChildrenFiltered(`div[role="list"]`).Length() == 0 {
			menuText0, menuLink0 = findATextAndHref(0, liSel0)
			dirName0 = fileNameFromMenuText(menuText0)
			fileName0 := "_index"

			createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0), fileName0, genIndexHeadStr(menuText0, baseUrl, menuLink0, liIndex0))
			// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
			httpGetContentAndWriteToHtmlFile(cmd, menuLink0, withSlashJoinStr(dist, dirName0, fileName0))
		} else {
			menuText0 = findAText(0, liSel0)
			dirName0 = fileNameFromMenuText(menuText0)
			fileName0 := "_index"
			createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0), fileName0, genIndexHeadStr(menuText0, baseUrl, menuLink0, liIndex0))
			// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
			createHtml(withSlashJoinStr(dist, dirName0, fileName0))

			liSelX1s := liSel0.ChildrenFiltered(`div[role="list"]`).ChildrenFiltered(`div[role="listitem"]`)
			liSelX1s.Each(func(liIndex1 int, liSel1 *goquery.Selection) {
				if liSel1.ChildrenFiltered(`div[role="list"]`).Length() == 0 {
					menuText1, menuLink1 = findATextAndHref(1, liSel1)
					fileName1 := fileNameFromMenuText(menuText1)
					createMdFileAndWrite(withSlashJoinStr(dist, dirName0, fileName1), genNonIndexHeadStr(menuText1, baseUrl, menuLink1, liIndex1))
					// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
					httpGetContentAndWriteToHtmlFile(cmd, menuLink1, withSlashJoinStr(dist, dirName0, fileName1))
				} else {
					menuText1, menuLink1 = findATextAndHref(1, liSel1)
					dirName1 = fileNameFromMenuText(menuText1)
					fileName1 := "_index"
					createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1), fileName1, genNonIndexHeadStr(menuText1, baseUrl, menuLink1, liIndex1))
					// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
					createHtml(withSlashJoinStr(dist, dirName0, dirName1, fileName1))

					liSelX2s := liSel1.ChildrenFiltered(`div[role="list"]`).ChildrenFiltered(`div[role="listitem"]`)
					liSelX2s.Each(func(liIndex2 int, liSel2 *goquery.Selection) {

						if liSel2.ChildrenFiltered(`div[role="list"]`).Length() == 0 {
							menuText2, menuLink2 = findATextAndHref(2, liSel2)
							fileName2 := fileNameFromMenuText(menuText2)
							createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, fileName2), genNonIndexHeadStr(menuText2, baseUrl, menuLink2, liIndex2))
							// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
							httpGetContentAndWriteToHtmlFile(cmd, menuLink2, withSlashJoinStr(dist, dirName0, dirName1, fileName2))
						} else {
							menuText2, menuLink2 = findATextAndHref(2, liSel2)
							dirName2 = fileNameFromMenuText(menuText2)
							fileName2 := "_index"
							createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2), fileName2, genNonIndexHeadStr(menuText2, baseUrl, menuLink2, liIndex2))
							// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
							createHtml(withSlashJoinStr(dist, dirName0, dirName1, dirName2, fileName2))

							liSelX3s := liSel2.ChildrenFiltered(`div[role="list"]`).ChildrenFiltered(`div[role="listitem"]`)
							liSelX3s.Each(func(liIndex3 int, liSel3 *goquery.Selection) {

								if liSel3.ChildrenFiltered(`div[role="list"]`).Length() == 0 {
									menuText3, menuLink3 = findATextAndHref(3, liSel3)
									fileName3 := fileNameFromMenuText(menuText3)
									createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, fileName3), genNonIndexHeadStr(menuText3, baseUrl, menuLink3, liIndex3))
									// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
									httpGetContentAndWriteToHtmlFile(cmd, menuLink3, withSlashJoinStr(dist, dirName0, dirName1, dirName2, fileName3))
								} else {
									menuText3, menuLink3 = findATextAndHref(3, liSel3)
									dirName3 = fileNameFromMenuText(menuText3)
									fileName3 := "_index"
									createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3), fileName3, genNonIndexHeadStr(menuText3, baseUrl, menuLink3, liIndex3))
									// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
									createHtml(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, fileName3))

									liSelX4s := liSel3.ChildrenFiltered(`div[role="list"]`).ChildrenFiltered(`div[role="listitem"]`)
									liSelX4s.Each(func(liIndex4 int, liSel4 *goquery.Selection) {

										if liSel4.ChildrenFiltered(`div[role="list"]`).Length() == 0 {
											menuText4, menuLink4 = findATextAndHref(4, liSel4)
											fileName4 := fileNameFromMenuText(menuText4)
											createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, fileName4), genNonIndexHeadStr(menuText4, baseUrl, menuLink4, liIndex4))
											// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
											httpGetContentAndWriteToHtmlFile(cmd, menuLink4, withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName4, fileName4))
										} else {
											menuText4, menuLink4 = findATextAndHref(4, liSel4)
											dirName4 = fileNameFromMenuText(menuText4)
											fileName4 := "_index"
											createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4), fileName4, genNonIndexHeadStr(menuText4, baseUrl, menuLink4, liIndex4))
											// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
											createHtml(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, fileName4))

											liSelX5s := liSel4.ChildrenFiltered(`div[role="list"]`).ChildrenFiltered(`div[role="listitem"]`)
											liSelX5s.Each(func(liIndex5 int, liSel5 *goquery.Selection) {

												if liSel5.ChildrenFiltered(`div[role="list"]`).Length() == 0 {
													menuText5, menuLink5 = findATextAndHref(5, liSel5)
													fileName5 := fileNameFromMenuText(menuText5)
													createMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, fileName5), genNonIndexHeadStr(menuText5, baseUrl, menuLink5, liIndex5))
													// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
													httpGetContentAndWriteToHtmlFile(cmd, menuLink5, withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, fileName5))
												} else {
													menuText5, menuLink5 = findATextAndHref(5, liSel5)
													dirName5 = fileNameFromMenuText(menuText5)
													fileName5 := "_index"
													createDirAndMdFileAndWrite(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, dirName5), fileName5, genNonIndexHeadStr(menuText5, baseUrl, menuLink5, liIndex5))
													// 获取链接html指定选择器下的内容到新创建的以.html结尾的文件中
													createHtml(withSlashJoinStr(dist, dirName0, dirName1, dirName2, dirName3, dirName4, dirName5, fileName5))

													liSelX6s := liSel5.ChildrenFiltered(`div[role="list"]`).ChildrenFiltered(`div[role="listitem"]`)
													liSelX6s.Each(func(liIndex6 int, liSel6 *goquery.Selection) {

														if liSel6.ChildrenFiltered(`div[role="list"]`).Length() == 0 {
															menuText6, menuLink6 = findATextAndHref(6, liSel6)
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
