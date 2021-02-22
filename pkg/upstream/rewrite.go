package upstream

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/justinas/alice"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/logger"
)

// newRewritePath creates a new middleware that will rewrite the request URI
// path before handing the request to the next server.
func newRewritePath(rewriteRegExp *regexp.Regexp, rewriteTarget string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return rewritePath(rewriteRegExp, rewriteTarget, next)
	}
}

// rewritePath uses the regexp to rewrite the request URI based on the provided
// rewriteTarget.
func rewritePath(rewriteRegExp *regexp.Regexp, rewriteTarget string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		reqURL, err := url.Parse(req.RequestURI)
		if err != nil {
			logger.Errorf("could not parse request URI: %v", err)
			// Since the requestURI is created by decoding the request,
			// this should never happen.
			// It's ok in this case to just return a 500.
			rw.WriteHeader(500)
			return
		}

		// Use the regex to rewrite the request path before proxying to the upstream.
		reqURL.Path = rewriteRegExp.ReplaceAllString(reqURL.Path, rewriteTarget)
		req.RequestURI = reqURL.String()

		next.ServeHTTP(rw, req)
	})
}
