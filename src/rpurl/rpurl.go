package rpurl

import (
	"errors"
	"net/url"
)

func IsRelative(urlIn string) (isRelative bool, err error) {
	u, err := url.Parse(urlIn)
	if err != nil {
		return false, err
	}

	return !u.IsAbs(), nil
}

func IsHttpScheme(urlIn string) (isHttpScheme bool, err error) {
	u, err := url.Parse(urlIn)
	if err != nil {
		return false, err
	}

	return u.Scheme == "http" || u.Scheme == "https", nil
}

func IsSameDomain(url1 string, url2 string) (isSameDomain bool) {
	u1, err := url.Parse(url1)
	if err != nil {
		return false
	}

	u2, err := url.Parse(url2)
	if err != nil {
		return false
	}

	return u1.Host == u2.Host
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

	u = b.ResolveReference(u)

	return u.String(), nil
}

func AbsoluteHttpUrl(baseUrl string, href string) (url string, err error) {
	isRelative, err := IsRelative(href)
	if err != nil {
		return "", err
	}

	isHttpScheme, err := IsHttpScheme(href)
	if err != nil {
		return "", err
	}

	if !isRelative && !isHttpScheme {
		return "", errors.New("Invalid URL " + href + ".")
	}

	if isRelative {
		href, err = AddBaseHost(baseUrl, href)
		if err != nil {
			return "", err
		}
	}

	return href, nil
}

func AbsoluteHttpUrls(baseUrl string, hrefs []string) (urls []string, errs []error) {
	for _, href := range hrefs {
		url, err := AbsoluteHttpUrl(baseUrl, href)
		if err != nil {
			errs = append(errs, err)
		} else {
			urls = append(urls, url)
		}
	}

	return urls, errs
}
