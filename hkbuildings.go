package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var cookiesMap = map[string]string{}
var names = map[string]string{}

func getCookie(resp *http.Response) string {
	for _, cookie := range resp.Cookies() {
		cookiesMap[cookie.Name] = cookie.Value
	}
	var c []string
	for key, value := range cookiesMap {
		c = append(c, key+"="+value)
	}
	return strings.Join(c, "; ")
}

func getViewState(b []byte) (pageViewState string) {
	re := regexp.MustCompile("javax.faces.ViewState.+?value=\"([^\"]+?)\"")
	m := re.FindAllSubmatch(b, -1)
	pageViewState = string(m[0][1])
	return
}

func getPage() (pageCookie, pageViewState string) {
	req, _ := http.NewRequest("GET", "https://bmis2.buildingmgt.gov.hk/bd_hadbiex/content/searchbuilding/building_search.jsf?renderedValue=true", nil)
	req.Header.Set("Accept-Language", "zh-TW;q=1.0")
	client := &http.Client{}
	resp, _ := client.Do(req)
	pageCookie = getCookie(resp)
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	pageViewState = getViewState(b)
	return
}

func getDocument(pageCookie, pageViewState string, start int) (string, error) {
	data := url.Values{}

	data.Set("javax.faces.partial.ajax", "true")
	data.Set("javax.faces.source", "_bld_result_frm:_result_tbl")
	data.Set("javax.faces.partial.execute", "_bld_result_frm:_result_tbl")
	data.Set("javax.faces.partial.render", "_bld_result_frm:_result_tbl")
	data.Set("javax.faces.behavior.event", "page")
	data.Set("javax.faces.partial.event", "page")
	data.Set("_bld_result_frm:_result_tbl_pagination", "true")
	data.Set("_bld_result_frm:_result_tbl_first", fmt.Sprintf("%d", start))
	data.Set("_bld_result_frm:_result_tbl_rows", "10")
	data.Set("_bld_result_frm:_result_tbl_encodeFeature", "true")
	data.Set("_bld_result_frm:_result_tbl_columnOrder", "_bld_result_frm:_result_tbl:j_id_4c,_bld_result_frm:_result_tbl:j_id_4i,_bld_result_frm:_result_tbl:j_id_4o")
	data.Set("_bld_result_frm_SUBMIT", "1")
	data.Set("autoScroll", "")
	data.Set("javax.faces.ViewState", pageViewState)

	req, err := http.NewRequest("POST", "https://bmis2.buildingmgt.gov.hk/bd_hadbiex/content/searchbuilding/building_search.jsf", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Cookie", pageCookie)
	req.Header.Set("Accept-Language", "zh-TW;q=1.0")
	req.Header.Set("Faces-Request", "partial/ajax")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://bmis2.buildingmgt.gov.hk")
	req.Header.Set("Referer", "https://bmis2.buildingmgt.gov.hk/bd_hadbiex/content/searchbuilding/building_search.jsf?renderedValue=true")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	ss := bytes.Index(b, []byte("<![CDATA[<tr")) + 9
	ee := bytes.Index(b, []byte("</tr>]]>")) + 5
	html := string(b[ss:ee])

	r1 := regexp.MustCompile("<tr.+?</tr>")
	r2 := regexp.MustCompile("<a([^>]+?)>([^>]+?)</a>")
	r3 := regexp.MustCompile("id=\"[^\"]+?:(\\d+)[^\"]+?\"")

	m1 := r1.FindAllString(html, -1)
	for _, match := range m1 {
		m2 := r2.FindAllStringSubmatch(match, -1)
		building := m2[0]
		names[r3.FindStringSubmatch(building[1])[1]] = building[2]
	}

	pageCookie = getCookie(resp)

	return "", nil
}

func foo(pageCookie, pageViewState string, index int) *http.Response {
	data := url.Values{}
	data.Set("_bld_result_frm:_result_tbl_columnOrder", "_bld_result_frm:_result_tbl:j_id_4c,_bld_result_frm:_result_tbl:j_id_4i,_bld_result_frm:_result_tbl:j_id_4o")
	data.Set("_bld_result_frm_SUBMIT", "1")
	data.Set("autoScroll", "")
	data.Set("javax.faces.ViewState", pageViewState)
	key := fmt.Sprintf("_bld_result_frm:_result_tbl:%d:j_id_4e", index)
	data.Set(key, key)
	req, _ := http.NewRequest("POST", "https://bmis2.buildingmgt.gov.hk/bd_hadbiex/content/searchbuilding/building_search.jsf", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Cookie", pageCookie)
	req.Header.Set("Accept-Language", "zh-TW;q=1.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://bmis2.buildingmgt.gov.hk")
	req.Header.Set("Referer", "https://bmis2.buildingmgt.gov.hk/bd_hadbiex/content/searchbuilding/building_search.jsf?renderedValue=true")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

func get(pageCookie, pageViewState string, index int) map[string]string {
	resp := foo(pageCookie, pageViewState, index)
	pageCookie = getCookie(resp)
	url := resp.Header.Get("Location")
	// fmt.Fprintln(os.Stderr, url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", pageCookie)
	req.Header.Set("Accept-Language", "zh-TW;q=1.0")
	req.Header.Set("Referer", "https://bmis2.buildingmgt.gov.hk/bd_hadbiex/content/searchbuilding/building_search.jsf?renderedValue=true")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	pageCookie = getCookie(resp)
	doc, _ := goquery.NewDocumentFromResponse(resp)
	idx := fmt.Sprintf("%d", index)
	var da = map[string]string{
		"ID":   idx,
		"大廈名稱": names[idx],
	}
	doc.Find(".ax-content").Each(func(_ int, s *goquery.Selection) {
		s.Find("[role=grid]").Each(func(_ int, s2 *goquery.Selection) {
			col := strings.TrimSpace(s2.Find("[role=columnheader]").Text())
			value := s2.Find("[role=gridcell]").Map(func(_ int, s3 *goquery.Selection) string {
				return clean(s3.Text())
			})
			da[col] = strings.Join(value, "; ")
		})
		s.Find("[data-widget=widget__detail_form_j_id_2k] .col").Each(func(_ int, s2 *goquery.Selection) {
			col := strings.TrimSpace(s2.Find(".label").Text())
			if col == "大廈名稱" {
				return
			}
			value := s2.Find(".text").Map(func(_ int, s3 *goquery.Selection) string {
				return clean(s3.Text())
			})
			da[col] = strings.Join(value, "; ")
		})
		var org []string
		s.Find("[data-widget=widget__detail_form_j_id_3r] .field").Each(func(i int, s2 *goquery.Selection) {
			if i%2 == 1 {
				org = append(org, "=>")
			} else if i > 0 {
				org = append(org, "; ")
			}
			org = append(org, clean(s2.Text()))
		})
		if len(org) > 0 {
			da["大廈組織"] = strings.Join(org, "")

		}

	})
	return da
}

func clean(i string) string {
	return strings.Replace(strings.Replace(strings.TrimSpace(i), "\n", "", -1), "\t", "", -1)
}

func main() {
	a, _ := strconv.Atoi(os.Args[1])
	b, _ := strconv.Atoi(os.Args[2])

	fmt.Fprintf(os.Stderr, "Getting %d to %d\n", a, b)

	var pageCookie, pageViewState string
	for i := a; i < b; i++ {
		if (i-a)%10 == 0 {
			fmt.Fprintln(os.Stderr, "Get Page")
			pageCookie, pageViewState = getPage()
			getDocument(pageCookie, pageViewState, i)
		}

		t, _ := json.Marshal(get(pageCookie, pageViewState, i))
		fmt.Println(string(t))
	}
}
