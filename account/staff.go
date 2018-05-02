package account

import (
	"net/url"

	"github.com/pkg/errors"
)

func (u *staff) Login() error {
	values := url.Values{}
	values.Add("client_id", u.clientID)
	values.Add("email", u.loginID)
	values.Add("password", u.password)
	values.Add("login_type", "1")
	values.Add("url", "/employee")
	res, err := u.httpClient.PostForm("https://ssl.jobcan.jp/login/pc-employee", values)
	if err != nil {
		return errors.Wrap(err, "failed to post to login form")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.Wrap(err, "Login error StatusCode="+string(res.StatusCode))
	}
	return nil
}
