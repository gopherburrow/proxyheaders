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

//Package proxiedhandler contains Proxy Headers Utilites Handlers.
package proxiedhandler

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/gopherburrow/proxyheaders"
)

//Some Constants used in namespaces and error strings.
const (
	ctxErrorValue = "gitlab.com/gopherburrow/proxyheaders/proxiedhandler Error"
)

//Used in request contexts. Go suggests using a specific type different from string for context keys.
type ctxType string

//The key used to store the error used in proxy header parsing.
//So it is possible to retrieve it inside a ErrorHandler.
var ctxError = ctxType(ctxErrorValue)

//ProxyHandler is a handler that process the, widely used in reverse proxies, headers X-Forwarded-*,
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
type ProxiedHandler struct {
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
//and dispatches for the Handler or, in case of errors the ErrorHandler will be called.
func (ph *ProxiedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//In case of no Handler defined a 404 Not Implemented will be served. It is not possible override this
	//behavior because this is symply NOT THE OBJECTIVE of this handler.
	if ph.Handler == nil {
		http.Error(w, fmt.Sprintf("%d - %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
		return
	}

	//Tranlate the headers in request fields.
	tr, err := proxyheaders.InjectIntoNewRequest(r)

	//If there is no error simply serve the handler.
	if err == nil {
		ph.Handler.ServeHTTP(w, tr)
		return
	}

	//In case of an error handler is defined, call ErrorHandler with the error in the context.
	if ph.ErrorHandler != nil {
		ph.ErrorHandler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxError, err)))
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
