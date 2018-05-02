package account

import (
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

func (u *admin) Login() error {
	values := url.Values{}
	values.Add("client_login_id", u.clientID)
	values.Add("client_manager_login_id", u.loginID)
	values.Add("client_login_password", u.password)
	values.Add("login_type", "2")
	values.Add("url", "https://ssl.jobcan.jp/client/")
	res, err := u.httpClient.PostForm("https://ssl.jobcan.jp/login/client", values)
	if err != nil {
		return errors.Wrap(err, "failed to post to login form")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.Wrap(err, "Login error StatusCode="+string(res.StatusCode))
	}
	u.employeeLogin()
	return nil
}

func (u *admin) employeeLogin() error {
	code, err := u.fetchEmployeeCode()
	if err != nil {
		return errors.Wrap(err, "failed to get employee code")
	}
	res, err := u.httpClient.Get("https://ssl.jobcan.jp/login/pc-employee/try?code=" + code)
	if err != nil {
		return errors.Wrap(err, "failed to login as an employee")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.Wrap(err, "Login error StatusCode="+string(res.StatusCode))
	}
	return nil
}

func (u *admin) fetchEmployeeCode() (string, error) {
	res, err := u.httpClient.Get("https://ssl.jobcan.jp/client")
	if err != nil {
		return "", errors.Wrap(err, "failed to get client page")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", errors.Wrap(err, "Login error StatusCode="+string(res.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to create a new document")
	}
	attr, _ := doc.Find("#rollover-menu > li:nth-child(2)").Attr("onclick")
	str := []byte(attr)
	assigned := regexp.MustCompile("code=([0-9a-f]+)")
	group := assigned.FindSubmatch(str)
	return string(group[1]), nil
}
