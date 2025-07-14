package qq

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/linktomarkdown/htxp"
	"github.com/zeromicro/go-zero/core/logx"
)

// 配置信息
const (
	AuthURL     = "https://graph.qq.com/oauth2.0/authorize"
	TokenURL    = "https://graph.qq.com/oauth2.0/token"
	OpenIDURL   = "https://graph.qq.com/oauth2.0/me"
	UserInfoURL = "https://graph.qq.com/user/get_user_info"
)

// Provider QQProvider QQ登录提供者
type Provider struct {
	AppID       string
	AppKey      string
	RedirectURL string
}

type PrivateInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
}

// NewProvider 实例化QQProvider
func NewProvider(AppID string, AppKey string, RedirectURL string) *Provider {
	return &Provider{
		AppID:       AppID,
		AppKey:      AppKey,
		RedirectURL: RedirectURL,
	}
}

// GetAuthCodeURL 获取授权码URL
func (p *Provider) GetAuthCodeURL(state string) string {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", p.AppID)
	params.Add("redirect_uri", p.RedirectURL)
	//params.Add("scope", "get_user_info")
	params.Add("state", state)
	return fmt.Sprintf("%s?%s", AuthURL, params.Encode())
}

// GetAccessToken 获取AccessToken
func (p *Provider) GetAccessToken(code string) (*PrivateInfo, error) {
	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", p.AppID)
	params.Add("client_secret", p.AppKey)
	params.Add("code", code)
	params.Add("redirect_uri", p.RedirectURL)
	tokenURL := fmt.Sprintf("%s?%s", TokenURL, params.Encode())
	resp, err := http.Get(tokenURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	body := string(bs)
	r := htxp.ConvertToMap(body)
	info := &PrivateInfo{}
	info.AccessToken = r["access_token"]
	info.RefreshToken = r["refresh_token"]
	info.ExpiresIn = r["expires_in"]
	return info, nil
}

// GetOpenID 获取OpenID
func (p *Provider) GetOpenID(accessToken string) (string, error) {
	params := url.Values{}
	params.Add("access_token", accessToken)
	openIDURL := fmt.Sprintf("%s?%s", OpenIDURL, params.Encode())
	resp, err := http.Get(openIDURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	body := string(bs)
	return body[45:77], nil
}

// UserInfo 获取用户信息
type UserInfo struct {
	Ret        int    `json:"ret"`
	Msg        string `json:"msg"`
	Nickname   string `json:"nickname"`
	Figureurl  string `json:"figureurl"`
	Figureurl1 string `json:"figureurl_1"`
	Figureurl2 string `json:"figureurl_2"`
	Gender     string `json:"gender"`
	Vip        string `json:"vip"`
	Level      string `json:"level"`
}

// GetUserInfo 获取用户信息
func (p *Provider) GetUserInfo(accessToken, openID string) (*UserInfo, error) {
	params := url.Values{}
	params.Add("access_token", accessToken)
	params.Add("oauth_consumer_key", p.AppID)
	params.Add("openid", openID)
	userInfoURL := fmt.Sprintf("%s?%s", UserInfoURL, params.Encode())
	resp, err := http.Get(userInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	var userInfo UserInfo
	err = json.Unmarshal(bs, &userInfo)
	if err != nil {
		logx.Errorf("Failed to parse user info JSON: %v", err)
	}
	return &userInfo, nil
}
