// This file is part of Gopher Burrow Proxy Headers Utilities.
//
// Gopher Burrow Proxy Headers Utilities is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Gopher Burrow Proxy Headers Utilities is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Gopher Burrow Proxy Headers Utilities.  If not, see <http://www.gnu.org/licenses/>.

//Package proxyheaders contains Proxy Headers Utilites. Most of it functions handle Proxy generated headers.
package proxyheaders

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
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

//NewProxiedRequest process the headers X-Forwarded-*, embed their values in a new request copied from r and return it, handling the errors.
func NewProxiedRequest(r *http.Request) (*http.Request, error) {
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
