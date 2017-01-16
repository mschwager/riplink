package rpurl

import (
	"net/url"
)

func HasHost(urlIn string) (hasHost bool, err error) {
	u, err := url.Parse(urlIn)
	if err != nil {
		return false, err
	}

	return u.Host != "", nil
}

func AddBaseHost(baseHost string, urlPath string) (urlOut string, err error) {
	b, err := url.Parse(baseHost)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(urlPath)
	if err != nil {
		return "", err
	}

	u.Scheme = b.Scheme
	u.Host = b.Host
	u.User = b.User

	result, err := url.QueryUnescape(u.String())
	if err != nil {
		return "", err
	}

	return result, nil
}
