package gocinside

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var errWriteGuestBook = errors.New("error writing guest book")

type GuestBook struct {
	client *DcClient
	userID string
}

func (c *DcClient) NewGuestBook(userID string) *GuestBook {
	return &GuestBook{
		client: c,
		userID: userID}
}

// 에러 발생 시 빈 문자열을 반환합니다.
func (g *GuestBook) Username() string {
	res, err := http.Get("https://gallog.dcinside.com/" + g.userID)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	html, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return ""
	}

	return html.Find("strong").Text()
}

// 방명록을 작성합니다.
// 비밀글로 작성하길 원하는 경우 WriteSecret() 메소드를 사용하세요.
func (g *GuestBook) Write(memo string) error { return g.write(memo, false) }

// 방명록을 작성합니다.
// 공개글로 작성하길 원하는 경우 Write() 메소드를 사용하세요.
func (g *GuestBook) WriteSecret(memo string) error { return g.write(memo, true) }

func (g *GuestBook) write(memo string, secret bool) error {
	formData := map[string]string{}
	formData["memo"] = memo

	if g.client.guest {
		formData["name"] = g.client.id
		formData["password"] = g.client.pw
	}

	isSecret := map[bool]int{
		true: 1, false: 0,
	}[secret]

	res, err := g.client.session.R().
		SetHeader("Referer", "https://gallog.dcinside.com/"+g.userID+"/guestbook").
		SetHeader("X-Requested-With", "XMLHttpRequest").
		SetFormData(formData).
		SetBody(map[string]any{
			"ci_t":      g.client.ci_c,
			"is_secret": isSecret}).
		Post("https://gallog.dcinside.com/" + g.userID + "/ajax/guestbook_ajax/write")

	if res.IsError() {
		return fmt.Errorf(
			"%v\nResponse: \n%s", errWriteGuestBook, string(res.Body()))
	}

	return err
}
