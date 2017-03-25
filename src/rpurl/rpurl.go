package rpurl

import (
	"net/url"
)

func IsRelative(urlIn string) (isRelative bool, err error) {
	u, err := url.Parse(urlIn)
	if err != nil {
		return false, err
	}

	return u.Host == "", nil
}

func IsHttpScheme(urlIn string) (isHttpScheme bool, err error) {
	u, err := url.Parse(urlIn)
	if err != nil {
		return false, err
	}

	// Assume lack of a URL scheme implies some form of HTTP
	return u.Scheme == "" || u.Scheme == "http" || u.Scheme == "https", nil
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
