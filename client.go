package gocinside

import "github.com/go-resty/resty/v2"

type DcClient struct {
	session *resty.Client
	id, pw  string
	ci_c    string
	guest   bool
}

func newClient(guest bool) *DcClient {
	client := new(DcClient)
	client.session = resty.New().
		SetHeader(
			"User-Agent",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36").
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(15)).
		SetContentLength(true)
	client.guest = guest
	return client
}

// 고정닉, 비고정닉 계정이 있는 경우 이 메소드를 통해 클라이언트를 생성합니다.
func NewMemberClient() *DcClient { return newClient(false) }

// 유동으로 이용할 경우 이 메소드를 통해 클라이언트를 생성합니다.
func NewGuestClient() *DcClient { return newClient(true) }

func (c *DcClient) SetID(id string) *DcClient {
	c.id = id
	return c
}

func (c *DcClient) SetPassword(pw string) *DcClient {
	c.pw = pw
	return c
}
