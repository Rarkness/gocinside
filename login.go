package gocinside

import (
	"errors"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	errNoHiddenInputFound = errors.New("no hidden input values found")
	errLogin              = errors.New("error logging in")
)

// (*DcClient).Login() 메소드와 동일한 역할을 하지만,
// 에러가 발생할 경우 에러를 반환하지 않습니다.
func (c *DcClient) LoginWithNoTrace() *DcClient {
	c.Login()
	return c
}

// 반드시 이 메소드를 통해 로그인 처리를 해야 코드가 정상 작동합니다.
// 에러 반환이 필요 없을 경우 LoginWithNoTrace() 메소드를 사용하세요.
func (c *DcClient) Login() error {
	res, err := c.getDcMain()
	if err != nil {
		return err
	}

	c.ci_c = c.getCookies(res)
	if c.guest {
		return nil
	}

	key, value := c.getHiddenInput(res)
	if key == "" || value == "" {
		return errNoHiddenInputFound
	}

	login, err := c.session.R().
		EnableTrace().
		SetHeader("Referer", "https://www.dcinside.com/").
		SetFormData(map[string]string{
			key:       value,
			"user_id": c.id,
			"pw":      c.pw,
			"s_url":   "//www.dcinside.com/",
			"ssl":     "Y",
			"ci_t":    c.ci_c,
		}).
		Post("https://sign.dcinside.com/login/member_check")

	if err != nil {
		return err
	}

	if login.IsError() {
		return errLogin
	}

	return nil
}

func (c *DcClient) getDcMain() (*http.Response, error) {
	res, err := http.Get("https://www.dcinside.com")
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *DcClient) getCookies(res *http.Response) (cic string) {
	for _, value := range res.Header.Values("Set-Cookie") {
		if !strings.Contains(value, "ci_c") {
			continue
		}

		for _, str := range strings.Split(value, "; ") {
			if sp := strings.Split(str, "="); len(sp) == 2 && sp[0] == "ci_c" {
				cic = sp[1]
			}
		}
	}
	return
}

func (c *DcClient) getHiddenInput(res *http.Response) (key, value string) {
	html, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", ""
	}

	hidden := html.Find("form > input").Get(2)

	for _, randKey := range hidden.Attr {
		switch randKey.Key {
		case "name":
			key = randKey.Val
		case "value":
			value = randKey.Val
		}
	}

	return
}
