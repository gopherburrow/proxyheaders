// This file is part of Riot Emergence Proxy Utillities.
//
// Riot Emergence Proxy Headers Utillities is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Riot Emergence Proxy Headers Utillities is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Riot Emergence Proxy Utillities.  If not, see <http://www.gnu.org/licenses/>.

//Package proxyheaders contains Proxy Headers Utilites. Most of it functions handle Proxy generated headers.
package proxyheaders

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
)

//Some Constants used in namespaces and error strings.
const (
	ctxErrorValue = "gitlab.com/gopherburrow/proxyheaders Error"
)

//Errors caused by misconfiguration or implementation fault when serving a request in ServeHTTP() method.
//
//These errors can be captured using the Error() method when inside an ErrorHandler.
//
//If not custom-handled using ErrorHandler these errors will return a "400 - Bad Request" page.
var (
	//ErrMustHaveXForwardedFor is returned when the X-Forwarded-For header is not present.
	ErrMustHaveXForwardedFor = errors.New("proxyheaders: must have X-Forwarded-For in headers")
	//ErrMustHaveXForwardedHost is returned when the X-Forwarded-Host header is not present.
	ErrMustHaveXForwardedHost = errors.New("proxyheaders: must have X-Forwarded-Host in headers")
	//ErrMustHaveXForwardedProto is returned when the X-Forwarded-Proto header is not present.
	ErrMustHaveXForwardedProto = errors.New("proxyheaders: must have X-Forwarded-Proto in headers")
	//ErrMustHaveXForwardedProto is returned when the X-Forwarded-Proto header is present, but has an invalid certificate value.
	ErrXForwardedClientCertMustBeValid = errors.New("proxyheaders: cannot parse the PEM encoded X.509 certificates in X-Forwarded-Client-Cert header")
)

//Used in request contexts. Go suggests using a specific type different from string for context keys.
type ctxType string

//The key used to store the error used in proxy header parsing.
//So it is possible to retrieve it inside a ErrorHandler.
var ctxError = ctxType(ctxErrorValue)

//XForwardedRequest process the headers X-Forwarded-*, embed their values in a new request and return it, handling the errors.
func xForwardedRequest(r *http.Request) (*http.Request, error) {
	//Some Constants used in namespaces and error strings.
	//Extract and test the expected X-Forwarded-* headers, returning errors if any of them are missed.
	xfh := r.Header.Get("X-Forwarded-Host")
	if xfh == "" {
		return nil, ErrMustHaveXForwardedHost
	}
	xff := r.Header.Get("X-Forwarded-For")
	if xff == "" {
		return nil, ErrMustHaveXForwardedFor
	}
	xfp := r.Header.Get("X-Forwarded-Proto")
	if xfp == "" {
		return nil, ErrMustHaveXForwardedProto
	}

	//Create a copy of the request...
	rCopy := new(http.Request)
	*rCopy = *r

	//..and remove the headers so there is no confusion if the request came from a
	//handler that already embed the headers.
	rCopy.Header.Del("X-Forwarded-Host")
	rCopy.Header.Del("X-Forwarded-For")
	rCopy.Header.Del("X-Forwarded-Proto")

	//Embed the headers...
	rCopy.Host = xfh
	rCopy.RemoteAddr = xff

	//If it is not https there is nothing else to do. Skip what remmains.
	if xfp != "https" {
		return rCopy, nil
	}

	//In case there is https processing create a dummy TLS field.
	rCopy.TLS = &tls.ConnectionState{}

	//Extract (and remove from request) possible client certificates, and if there is none, skip certificate processing.
	xfcc := r.Header.Get("X-Forwarded-Client-Cert")
	rCopy.Header.Del("X-Forwarded-Client-Cert")
	if xfcc == "" {
		return rCopy, nil
	}

	//In case there are client certificates, decode the certificates in PEM blocks and return them in the request.
	var block *pem.Block
	pemRemainder := []byte(xfcc)
	certs := make([]*x509.Certificate, 0)
	for {
		block, pemRemainder = pem.Decode(pemRemainder)
		if block == nil {
			break
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, ErrXForwardedClientCertMustBeValid
		}
		certs = append(certs, cert)
	}
	rCopy.TLS.PeerCertificates = certs
	return rCopy, nil
}

//XForwardedHandler is a handler that process the, widely used in reverse proxies, headers X-Forwarded-*,
//embed their values in a new request, remove the headers like it were generated without the proxy and,
//if there is no errors, call the Handler.
//
//It is possible to create a special error handler for each case.
//
//It is highly recomended from a security standpoint that the internet inbound proxy does not accept these headers to
//avoid injection by a malicious agent.
//
//The following headers are processed:
//
//• X-Forwarded-Host: translates to http.Request.Host [required];
//
//• X-Forwarded-For: translates to http.Request.RemoteAddr [required];
//
//• X-Forwarded-Proto: translates to a default http.Request.TLS if the value is "https", or simply do nothing if "http" [required];
//
//• X-Forwarded-Client-Cert: translates to parsed PEM X.509 Certificates in http.Request.TLS.PeerCertificates if the value of proto was "https" [optional].
type XForwardedHandler struct {
	//Handler that will be called in case of all required proxy headers are present.
	//If nil a vanilla "404 - Not Found" will be served.
	Handler http.Handler
	//ErrorHandler that will be called in case of any required proxy headers are absent or malformed or
	//any error during the parsing of certificates.
	//It is possible to retrieve the error in the request with the request context value: .
	//If nil a vanilla "400 - Bad Request" will be served.
	ErrorHandler http.Handler
}

//ServeHTTP is the method that dispatches requests that came from proxies, transform the headers in the according http.Request fields,
//and dispatches for the Handler or, in case of errors the error Handlers will be called.
func (xfhh *XForwardedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Tranlate the headers in request fields.
	tr, err := xForwardedRequest(r)

	//In case of no Handler defined a 404 Not Implemented will be served. It is not possible override this
	//behavior because it is not the objective os this handler.
	if err == nil && xfhh.Handler == nil {
		http.Error(w, fmt.Sprintf("%d - %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
		return
	}

	//If there is no error simply serve the handler.
	if err == nil {
		xfhh.Handler.ServeHTTP(w, tr)
		return
	}

	//In case of an error handler is defined, call ErrorHandler with the error in the context.
	if xfhh.ErrorHandler != nil {
		xfhh.ErrorHandler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxError, err)))
		return
	}

	//If there is no handlers defined for each specific error, serve a vanilla "400 - Bad Request".
	http.Error(w, fmt.Sprintf("%d - %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
	return
}

//Error retrieves the proxy parsing error, when inside XForwardedHandler.ErrorHandler.
//If called outside an XForwardedHandler.ErrorHandler it will retun nil.
func Error(r *http.Request) error {
	m, ok := r.Context().Value(ctxError).(error)
	if !ok {
		return nil
	}
	return m
}
