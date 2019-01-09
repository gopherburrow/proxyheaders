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

//Package proxiedhandler_test contains Proxy Headers Utility Handlers tests.
package proxiedhandler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/gopherburrow/proxyheaders"
	"gitlab.com/gopherburrow/proxyheaders/proxiedhandler"
)

const validCert = `
-----BEGIN CERTIFICATE-----
MIIF3DCCA8SgAwIBAgICEAMwDQYJKoZIhvcNAQELBQAwgcIxCzAJBgNVBAYTAkJS
MRkwFwYDVQQIDBBEaXN0cml0byBGZWRlcmFsMSQwIgYDVQQKDBtSaW90IEVtZXJn
ZW5jZSBPcmdhbml6YXRpb24xKTAnBgNVBAsMIENlcnRpZmljYXRlcyBJc3N1aW5n
IERlcGFydGFtZW50MRgwFgYDVQQDDA9JbnRlcm1lZGlhdGUgQ0ExLTArBgkqhkiG
9w0BCQEWHmNlcnRpZmljYXRlc0ByaW90ZW1lcmdlbmNlLm9yZzAeFw0xODA1Mjkx
MzE4MjRaFw0xOTA2MDgxMzE4MjRaMIHZMQswCQYDVQQGEwJCUjEZMBcGA1UECAwQ
RGlzdHJpdG8gRmVkZXJhbDEeMBwGA1UEBwwVUmlvdCBFbWVyZ2VuY2UgU3RyZWV0
MSQwIgYDVQQKDBtSaW90IEVtZXJnZW5jZSBPcmdhbml6YXRpb24xKzApBgNVBAsM
IkNlcnRpZmljYXRpb24gQXV0aG9yaXR5IERlcGFydG1lbnQxETAPBgNVBAMMCEpv
aG4gRG9lMSkwJwYJKoZIhvcNAQkBFhpqb2huLmRvZUByaW90ZW1lcmdlbmNlLm9y
ZzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMdtfkEUVPNCVAkIPP3F
gxXj7o82aagGd5+nRLwiLdAXXyMaCg9tmUbg59Fcupece/wwFvFobe4Ro+Ob9HIV
n8mXLOzDAIR6bHNjf7xNzT7rSDDb1mXuPGUv8L8NRINZ7UyMlKtsgixK7VM+p6Tx
fgbQA+3IBrY5nkKP/BmZukmf5RHGNyTavXd+V04aWatJ2BU0X9zxmyPjxsGYPqr3
m/YkIjZnzXal7QHTU6IV5v+1zOOOWST2NkYbMWgk9OMYT4tkfWDSxokw2TlSrNpN
lOTBBLHuzmDsVPESH2lQKD72dUvFZgAjKfOthxEkiNUWULKJSGDepzkWzSiEDhou
hO8CAwEAAaOBwjCBvzAJBgNVHRMEAjAAMBEGCWCGSAGG+EIBAQQEAwIFoDAwBglg
hkgBhvhCAQ0EIxYhUmlvdCBFbWVyZ2VuY2UgQ2xpZW50IENlcnRpZmljYXRlMB0G
A1UdDgQWBBRxhliVUvKetuIFYzTuMb9EvKnqRDAfBgNVHSMEGDAWgBRWmOHShoPM
4TTp0k4SpBPrEZ9RLTAOBgNVHQ8BAf8EBAMCBeAwHQYDVR0lBBYwFAYIKwYBBQUH
AwIGCCsGAQUFBwMEMA0GCSqGSIb3DQEBCwUAA4ICAQDfHEIJ0tZehQFqkNEjeFIL
/QZ483lDRRO9yF39acXXU/FKRsOuQCSuCX4rpM/y7Kn34jXsaIu45UOl22/t7NDm
VKrZ1+v/SNEMZOUEyaakJ0UCv9AsF7FUj3TB7ZpbDDlhQPDKig+UT6Zi7MRnuAF+
MXo9IP19HW0wOPwOJ05HVATqNaPRf8/YaH9NW+5A/dbo5YJ2AvJyNaVuvH8YZSf0
TLE7pmf7JjaBxQB3FrTltLFCaiHQp1P4ql6OnXe0zzuhtZw+MOiZnCMk/uUcHBYw
vLli78F5ubawcDeuzu3exc6XL9XSkE1JMrqY60qR6sWEeyuYT0tIiWWwGC4QqMnL
o3eagSyaN8V/+noOLUlObRClD293Q0+TH9IQw1IsPouCTsdaWrCvQqEq9tVca/VW
FovgEzuqQCesHgguuQW+qOQDCTVoKChyqnYYQNl33tceDfTYlVRkf127KAsRUgF5
RWPC/TLH39KV/z2NuGIVPHeA5KP9z/mrRr0HdJI7IgidhO0rPbxl1tt8ZMDN3xjZ
l57WeXdZnUA1dY+5ns+Yphc3EAtTksenLBlc1ANPDurLfwOrsoSrNMmpUV6AB/AY
0u6Bkw37QD5536tXZi65rokheVWDZ9YG/01IHr+z7Y+IwO7XERFyG4fvq57ReEa/
PgNj8cl6ndHEvt4IlMpLPQ==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIGfDCCBGSgAwIBAgICEAAwDQYJKoZIhvcNAQELBQAwgdYxCzAJBgNVBAYTAkJS
MRkwFwYDVQQIDBBEaXN0cml0byBGZWRlcmFsMR4wHAYDVQQHDBVSaW90IEVtZXJn
ZW5jZSBTdHJlZXQxJDAiBgNVBAoMG1Jpb3QgRW1lcmdlbmNlIE9yZ2FuaXphdGlv
bjErMCkGA1UECwwiQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkgRGVwYXJ0bWVudDEQ
MA4GA1UEAwwHUm9vdCBDQTEnMCUGCSqGSIb3DQEJARYYcm9vdGNhQHJpb3RlbWVy
Z2VuY2Uub3JnMB4XDTE4MDUwNzE0MDEzNFoXDTI4MDUwNDE0MDEzNFowgcIxCzAJ
BgNVBAYTAkJSMRkwFwYDVQQIDBBEaXN0cml0byBGZWRlcmFsMSQwIgYDVQQKDBtS
aW90IEVtZXJnZW5jZSBPcmdhbml6YXRpb24xKTAnBgNVBAsMIENlcnRpZmljYXRl
cyBJc3N1aW5nIERlcGFydGFtZW50MRgwFgYDVQQDDA9JbnRlcm1lZGlhdGUgQ0Ex
LTArBgkqhkiG9w0BCQEWHmNlcnRpZmljYXRlc0ByaW90ZW1lcmdlbmNlLm9yZzCC
AiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAN/invYVfdA2Vs09RIYa7317
R51sDLCJsazQQA0wxnvhlHcjjkAeiCbsg7jzSSk15nQRfXMr9oY/l7NoPm4n92Pa
BO98DD17UMTgVYl5019D31xaxQiet6hEww2KAnaMXzcw4//o4pCk9dz8cVAEBM1R
Ar3HhhHe9ojMfE37Gfe0C34RksFfSFehF8p47/g+GMFZ/t1awszh8qIjfTFS0Hrz
f9G/yn19DmesNTM22U83Nnu73MYG4gVuzrreaBbc16RWXdwWyIuG7ESVQklE9fO3
2efJR/thtzlS7D3ElBBqxg6AsSJ7fmQnUsLwUwaTYlca/TdYvzcPaZLOHig/3m+R
GvnP0FwA/NYWIUVdEm4wqudDevSS0+OcdGAM68zjabzkJQUXCjB+3njdhVBZHij3
mdh17Dk+BeIcrgiXJUIjEOklEzkTz6+lnTwaTPYptyrOgKkjBbeXLaoUsZD5mgw4
xMqJnSXlk6zECFl0gFI3ZuaGyL1sS/Y+/TDS+I1BB6+VNgoojXWFRiE2TPjP3/Uq
T3UugahkeK6Ga0J+kNJmI5w0j+dHUraeLRJ+gAR5Rk/0C7/jrtOTCIMG6H7w115f
QztL7t0U6re756RoU+86ENlmBFZFBdlz6LleA1IV2P1QkW2nTiTlIGZJDwcPDi41
2iiBe36Ai08/zR0q7K17AgMBAAGjZjBkMB0GA1UdDgQWBBRWmOHShoPM4TTp0k4S
pBPrEZ9RLTAfBgNVHSMEGDAWgBR181MuS/zd6wQwk6AW1UR6aRtB1DASBgNVHRMB
Af8ECDAGAQH/AgEAMA4GA1UdDwEB/wQEAwIBhjANBgkqhkiG9w0BAQsFAAOCAgEA
dq/ZXgoi3lMd/vPgMsTkkOm2xAgHqImaZlEiMUoYxMCpw/FTqZ+XKfOuhOYuTmJY
sQymhcTf40rclmAB8adCct9TrBY6RGRTZn3Ol4mts/fbJcEyeU/tP6kSbAmp1bPg
7jBBjSYUgcd+0wBUbcSKZUrssY3F4iLUtm+h3NNs5BRk2QGz2ZZEvfcIVfoNpnBT
cxfjhwDY+yfEj6O+MHFjLd7E3IQ2fbw/C0to8pGRE+0czG5JCjDN0xT+LIvwBEgr
NfRzYWZPavDHvCnT6/E28Y9MK4UXwjbLJmBdY7Jdvuph0rZEGLYiGbDOGjApQ+Go
zfafRPCqB7GHHvXcq6pfVAle/gG12uiRWQOP1MOqHuEOLWXbR4pa3akTbnBEg60j
FFolFuFVYUJDq1zzHXeuS4UhCDECWZgVACKr6ClIUH7T5N81+11OsGBuSP/bEKqF
MYhKH1+Ql/u8l3TpHQoe0JWrO3mloucREbtKXC5EN8/aYzV/S/I668HVbuWVhv30
HiH3hb3kxkvlwCinh//2xCD3di4nfNFzFp3jy4YAmDkCw734L2WX97UW5ka69U5G
+bx0pS3zyVYJ75eUt2sst/DuR+99nAELkYJasuIbZUSjiR7AVCiVbMepLVM6NRn8
P/EP/vuOmxtTbgEqUKvE/6Uqb3UX13Yx1B/v3RVFDrk=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIGlDCCBHygAwIBAgIJAPHBS3E/Ojn2MA0GCSqGSIb3DQEBCwUAMIHWMQswCQYD
VQQGEwJCUjEZMBcGA1UECAwQRGlzdHJpdG8gRmVkZXJhbDEeMBwGA1UEBwwVUmlv
dCBFbWVyZ2VuY2UgU3RyZWV0MSQwIgYDVQQKDBtSaW90IEVtZXJnZW5jZSBPcmdh
bml6YXRpb24xKzApBgNVBAsMIkNlcnRpZmljYXRpb24gQXV0aG9yaXR5IERlcGFy
dG1lbnQxEDAOBgNVBAMMB1Jvb3QgQ0ExJzAlBgkqhkiG9w0BCQEWGHJvb3RjYUBy
aW90ZW1lcmdlbmNlLm9yZzAeFw0xODA1MDcxNDAwNDhaFw0zODA1MDIxNDAwNDha
MIHWMQswCQYDVQQGEwJCUjEZMBcGA1UECAwQRGlzdHJpdG8gRmVkZXJhbDEeMBwG
A1UEBwwVUmlvdCBFbWVyZ2VuY2UgU3RyZWV0MSQwIgYDVQQKDBtSaW90IEVtZXJn
ZW5jZSBPcmdhbml6YXRpb24xKzApBgNVBAsMIkNlcnRpZmljYXRpb24gQXV0aG9y
aXR5IERlcGFydG1lbnQxEDAOBgNVBAMMB1Jvb3QgQ0ExJzAlBgkqhkiG9w0BCQEW
GHJvb3RjYUByaW90ZW1lcmdlbmNlLm9yZzCCAiIwDQYJKoZIhvcNAQEBBQADggIP
ADCCAgoCggIBAJ3yGITjiMHhFDyUbhBe4giEF9FlUjr1lBPo8hOYfLlRW2qx0TEy
LOFCJSPc4NtN8gt3qKjqBIQt+qahUJH5v6n+YGvS/q3eX9AA2iVjJ6VyAkSwf8YY
MIfnevkCuYH6tVATT90i4M0qKDFrzRZXOKNM9ndXt+0mpEDKU9RB/9dru3ge/iI/
Vi23W1Ix+Pdxa0coIVnEQk4sNVRVhM0omhhfWkUV5hQhpJOOMheH4fYQXdQUP0DC
l9sOHUOtwfTJDVoOp0qc08y0+2jmWBeNa0v0cPSKdRFgxTwvtIxmkXOVCbHSbfWC
QacArNu9vLfJzxt9kp1E2q/iW4FA30iiLCDngfgS5xiJwusOmwEroiAZvZht1rz7
Df52bcnpmEVHb6abyPsR/159yRzgD7+GTKjt/PHwNu1lBr7C7HH0owm63/MUNT0m
v8nK12t/Cez6TXmPYb4Xa3ZTG4ndCBklyyzSv0xHuNJGJz9wBNgTXoNhuTirSBxo
8Y+EyljoAa+g/wZgMTdgjvVFGjzZqAW+l1TdYrr/SNAJ3isikctHDO2Jpz3rEaE0
AD2L8SM0NQ6ct9HXYrGXeAKqT6dlB5dYBlUQ4Tvc8hqTjShXi5i6dG1s1tiC0di3
Tw8Ll9ppN/hB7VNAQK1xDFH1tC8pEIsgItDx0J9mBb2rYlpgV1T2T3lPAgMBAAGj
YzBhMB0GA1UdDgQWBBR181MuS/zd6wQwk6AW1UR6aRtB1DAfBgNVHSMEGDAWgBR1
81MuS/zd6wQwk6AW1UR6aRtB1DAPBgNVHRMBAf8EBTADAQH/MA4GA1UdDwEB/wQE
AwIBhjANBgkqhkiG9w0BAQsFAAOCAgEAlb0HJu8CSi8al6sTYS2sKvww/GPNIjOs
2yutdEJV+gJcQxfpx9r9pwLvdDHFHStp7DPt+tm4d9ph3CaBZcPGBasWN0xgkqxA
/NDwN37sr0p91E1LhUbY9qxmp/QTeadm5Ej4n2LpwckMFfJLFhHk1XiHb/Xvnpny
nYH+4WWN5hsV2ZQNd3ack7yG+B15UKNAxJ3aVdc44s+d9fOXPBogsc1HF+W4Rtuy
oEIrV57y3oeGx5sBauDk2AhnAvg/DCF9saJnl66H5PikcP7QDXGE208ImOH6AtEE
MabBQRIgmNeZI1lfT5Jk0rH2T97iT1oTbCUoSrJ5F4qjoITNmyUL0QUyRfXSWShO
O0UcXdOKAPQAOOv7C6hZsQcRwrZfd5KEkFmsQoXHaUHo2ztcIgra2UttYu87GpcW
CvZ7Qs4y9Ae8pEaExa8giYC4/z9FVBDIxsoqa7013yPIIXXKAr9xRE+IWaUQUxBu
yg/aHk3FLviVaoShfh7uEn7sOsM80PvzfSDCGXGnLe8pK+wT5sAoRbIQ/i1sdA28
cxqbRUNj8Qv770wPD9Kcq+WBw1gXCvpd9M7kmrW1qUAtc3EZYwD7w2HPbz9el54D
VxO4VAdujXqk5ps3oooZOeFx6Y5z+ddz66NAUlduvxdfKcq8SHQzTgjTFHUHimvQ
rYikYp2kMGA=
-----END CERTIFICATE-----
`

const invalidCert = `
-----BEGIN CERTIFICATE-----
MIIF3DCCA8SgAwIBAgICEAMwDQYJKoZIhvcNAQELBQAwgcIxCzAJBgNVBAYTAkJS
MRkwFwYDVQQIDBBEaXN0cml0byBGZWRlcmFsMSQwIgYDVQQKDBtSaW90IEVtZXJn
ZW5jZSBPcmdhbml6YXRpb24xKTAnBgNVBAsMIENlcnRpZmljYXRlcyBJc3N1aW5n
IERlcGFydGFtZW50MRgwFgYDVQQDDA9JbnRlcm1lZGlhdGUgQ0ExLTArBgkqhkiG
9w0BCQEWHmNlcnRpZmljYXRlc0ByaW90ZW1lcmdlbmNlLm9yZzAeFw0xODA1Mjkx
MzE4MjRaFw0xOTA2MDgxMzE4MjRaMIHZMQswCQYDVQQGEwJCUjEZMBcGA1UECAwQ
RGlzdHJpdG8gRmVkZXJhbDEeMBwGA1UEBwwVUmlvdCBFbWVyZ2VuY2UgU3RyZWV0
MSQwIgYDVQQKDBtSaW90IEVtZXJnZW5jZSBPcmdhbml6YXRpb24xKzApBgNVBAsM
IkNlcnRpZmljYXRpb24gQXV0aG9yaXR5IERlcGFydG1lbnQxETAPBgNVBAMMCEpv
aG4gRG9lMSkwJwYJKoZIhvcNAQkBFhpqb2huLmRvZUByaW90ZW1lcmdlbmNlLm9y
ZzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMdtfkEUVPNCVAkIPP3F
gxXj7o82aagGd5+nRLwiLdAXXyMaCg9tmUbg59Fcupece/wwFvFobe4Ro+Ob9HIV
n8mXLOzDAIR6bHNjf7xNzT7rSDDb1mXuPGUv8L8NRINZ7UyMlKtsgixK7VM+p6Tx
fgbQA+3IBrY5nkKP/BmZukmf5RHGNyTavXd+V04aWatJ2BU0X9zxmyPjxsGYPqr3
m/YkIjZnzXal7QHTU6IV5v+1zOOOWST2NkYbMWgk9OMYT4tkfWDSxokw2TlSrNpN
lOTBBLHuzmDsVPESH2lQKD72dUvFZgAjKfOthxEkiNUWULKJSGDepzkWzSiEDhou
hO8CAwEAAaOBwjCBvzAJBgNVHRMEAjAAMBEGCWCGSAGG+EIBAQQEAwIFoDAwBglg
hkgBhvhCAQ0EIxYhUmlvdCBFbWVyZ2VuY2UgQ2xpZW50IENlcnRpZmljYXRlMB0G
A1UdDgQWBBRxhliVUvKetuIFYzTuMb9EvKnqRDAfBgNVHSMEGDAWgBRWmOHShoPM
4TTp0k4SpBPrEZ9RLTAOBgNVHQ8BAf8EBAMCBeAwHQYDVR0lBBYwFAYIKwYBBQUH
AwIGCCsGAQUFBwMEMA0GCSqGSIb3DQEBCwUAA4ICAQDfHEIJ0tZehQFqkNEjeFIL
/QZ483lDRRO9yF39acXXU/FKRsOuQCSuCX4rpM/y7Kn34jXsaIu45UOl22/t7NDm
VKrZ1+v/SNEMZOUEyaakJ0UCv9AsF7FUj3TB7ZpbDDlhQPDKig+UT6Zi7MRnuAF+
TLE7pmf7JjaBxQB3FrTltLFCaiHQp1P4ql6OnXe0zzuhtZw+MOiZnCMk/uUcHBYw
vLli78F5ubawcDeuzu3exc6XL9XSkE1JMrqY60qR6sWEeyuYT0tIiWWwGC4QqMnL
o3eagSyaN8V/+noOLUlObRClD293Q0+TH9IQw1IsPouCTsdaWrCvQqEq9tVca/VW
FovgEzuqQCesHgguuQW+qOQDCTVoKChyqnYYQNl33tceDfTYlVRkf127KAsRUgF5
RWPC/TLH39KV/z2NuGIVPHeA5KP9z/mrRr0HdJI7IgidhO0rPbxl1tt8ZMDN3xjZ
l57WeXdZnUA1dY+5ns+Yphc3EAtTksenLBlc1ANPDurLfwOrsoSrNMmpUV6AB/AY
0u6Bkw37QD5536tXZi65rokheVWDZ9YG/01IHr+z7Y+IwO7XERFyG4fvq57ReEa/
PgNj8cl6ndHEvt4IlMpLPQ==
-----END CERTIFICATE-----
`

func ErrorHandlerFunc(w http.ResponseWriter, r *http.Request) {
	err := proxiedhandler.Error(r)
	if err == nil {
		http.Error(w, "500 - Error Not Found", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, err.Error())
}

func DumpServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain;charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
}

func TestProxiedHandler_ServeHTTP_Success(t *testing.T) {
	xfh := &proxiedhandler.ProxiedHandler{
		Handler: http.HandlerFunc(DumpServeHTTP),
	}

	{
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
		req.Header.Add("X-Forwarded-For", "1.2.3.4")
		req.Header.Add("X-Forwarded-Host", "www.example.com")
		req.Header.Add("X-Forwarded-Proto", "http")

		rr := httptest.NewRecorder()
		xfh.ServeHTTP(rr, req)
		if want, got := http.StatusOK, rr.Code; want != got {
			t.Fatalf("want=%d, got=%d", want, got)
		}
	}

	{
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
		req.Header.Add("X-Forwarded-For", "1.2.3.4")
		req.Header.Add("X-Forwarded-Host", "www.example.com")
		req.Header.Add("X-Forwarded-Proto", "https")

		rr := httptest.NewRecorder()
		xfh.ServeHTTP(rr, req)
		if want, got := http.StatusOK, rr.Code; want != got {
			t.Fatalf("want=%d, got=%d", want, got)
		}
	}

	{
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
		req.Header.Add("X-Forwarded-For", "1.2.3.4")
		req.Header.Add("X-Forwarded-Host", "www.example.com")
		req.Header.Add("X-Forwarded-Proto", "https")
		req.Header.Add("X-Forwarded-Client-Cert", validCert)

		rr := httptest.NewRecorder()
		xfh.ServeHTTP(rr, req)
		if want, got := http.StatusOK, rr.Code; want != got {
			t.Fatalf("want=%d, got=%d", want, got)
		}
	}
}

func TestProxiedHandler_ServeHTTP_failMustHaveXForwardedHostWithUndefinedErrorHandler(t *testing.T) {
	xfh := &proxiedhandler.ProxiedHandler{
		Handler: http.HandlerFunc(DumpServeHTTP),
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("X-Forwarded-Proto", "https")
	req.Header.Add("X-Forwarded-Client-Cert", validCert)

	rr := httptest.NewRecorder()
	xfh.ServeHTTP(rr, req)
	if want, got := http.StatusBadRequest, rr.Code; want != got {
		t.Fatalf("want=%d, got=%d", want, got)
	}
}

func TestProxiedHandler_ServeHTTP_failMustHaveXForwardedHost(t *testing.T) {
	xfh := &proxiedhandler.ProxiedHandler{
		Handler:      http.HandlerFunc(DumpServeHTTP),
		ErrorHandler: http.HandlerFunc(ErrorHandlerFunc),
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("X-Forwarded-Proto", "https")
	req.Header.Add("X-Forwarded-Client-Cert", validCert)
	rr := httptest.NewRecorder()
	xfh.ServeHTTP(rr, req)
	if want, got := http.StatusBadRequest, rr.Code; want != got {
		t.Fatalf("want=%d, got=%d", want, got)
	}

	if want, got := proxyheaders.ErrMustHaveXForwardedHost.Error(), rr.Body.String(); want != got {
		t.Fatalf("want='%q', got='%q'", want, got)
	}
}

func TestProxiedHandler_ServeHTTP_failMustHaveXForwardedFor(t *testing.T) {
	xfh := &proxiedhandler.ProxiedHandler{
		Handler:      http.HandlerFunc(DumpServeHTTP),
		ErrorHandler: http.HandlerFunc(ErrorHandlerFunc),
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	req.Header.Add("X-Forwarded-Host", "www.example.com")
	req.Header.Add("X-Forwarded-Proto", "https")
	req.Header.Add("X-Forwarded-Client-Cert", validCert)
	rr := httptest.NewRecorder()
	xfh.ServeHTTP(rr, req)
	if want, got := http.StatusBadRequest, rr.Code; want != got {
		t.Fatalf("want=%d, got=%d", want, got)
	}
	if want, got := proxyheaders.ErrMustHaveXForwardedFor.Error(), rr.Body.String(); want != got {
		t.Fatalf("want=%q, got=%q", want, got)
	}
}

func TestProxiedHandler_ServeHTTP_failMustHaveXForwardedProto(t *testing.T) {
	xfh := &proxiedhandler.ProxiedHandler{
		Handler:      http.HandlerFunc(DumpServeHTTP),
		ErrorHandler: http.HandlerFunc(ErrorHandlerFunc),
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("X-Forwarded-Host", "www.example.com")
	req.Header.Add("X-Forwarded-Client-Cert", validCert)
	rr := httptest.NewRecorder()
	xfh.ServeHTTP(rr, req)
	if want, got := http.StatusBadRequest, rr.Code; want != got {
		t.Fatalf("want=%d, got=%d", want, got)
	}
	if want, got := proxyheaders.ErrMustHaveXForwardedProto.Error(), rr.Body.String(); want != got {
		t.Fatalf("want=%q, got=%q", want, got)
	}
}

func TestProxiedHandler_ServeHTTP_failInvalidXForwardedClientCert(t *testing.T) {
	xfh := &proxiedhandler.ProxiedHandler{
		Handler:      http.HandlerFunc(DumpServeHTTP),
		ErrorHandler: http.HandlerFunc(ErrorHandlerFunc),
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("X-Forwarded-Host", "www.example.com")
	req.Header.Add("X-Forwarded-Proto", "https")
	req.Header.Add("X-Forwarded-Client-Cert", invalidCert)

	rr := httptest.NewRecorder()
	xfh.ServeHTTP(rr, req)
	if want, got := http.StatusBadRequest, rr.Code; want != got {
		t.Fatalf("want=%d, got=%d", want, got)
	}
	if want, got := proxyheaders.ErrXForwardedClientCertMustBeValid.Error(), rr.Body.String(); want != got {
		t.Fatalf("want=%q, got=%q", want, got)
	}
}

func TestProxiedHandler_ServeHTTP_failUndefinedHandler(t *testing.T) {
	xfh := &proxiedhandler.ProxiedHandler{}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("X-Forwarded-Host", "www.example.com")
	req.Header.Add("X-Forwarded-Proto", "https")
	req.Header.Add("X-Forwarded-Client-Cert", validCert)

	rr := httptest.NewRecorder()
	xfh.ServeHTTP(rr, req)
	if want, got := http.StatusNotFound, rr.Code; want != got {
		t.Fatalf("want=%d, got=%d", want, got)
	}
}

func TestProxiedHandler_GetError_failCalledOutsideErrorHandler(t *testing.T) {
	xfh := &proxiedhandler.ProxiedHandler{
		Handler: http.HandlerFunc(ErrorHandlerFunc),
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("X-Forwarded-Host", "www.example.com")
	req.Header.Add("X-Forwarded-Proto", "https")
	req.Header.Add("X-Forwarded-Client-Cert", validCert)

	rr := httptest.NewRecorder()
	xfh.ServeHTTP(rr, req)
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("want=%d, got=%d", want, got)
	}
}
