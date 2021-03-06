package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kalmhq/kalm/api/resources"
	"github.com/kalmhq/kalm/controller/api/v1alpha1"
	"github.com/kalmhq/kalm/controller/controllers"
	"github.com/stretchr/testify/suite"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
	"testing"
)

type HttpsCertTestSuite struct {
	WithControllerTestSuite
}

func TestHttpsCertTestSuite(t *testing.T) {
	suite.Run(t, new(HttpsCertTestSuite))
}

func (suite *HttpsCertTestSuite) SetupSuite() {
	suite.WithControllerTestSuite.SetupSuite()

	suite.ensureNamespaceExist("istio-system")
	suite.ensureNamespaceExist(controllers.CertManagerNamespace)
}

func (suite *HttpsCertTestSuite) TearDownTest() {
	suite.ensureObjectDeleted(&v1alpha1.HttpsCert{ObjectMeta: metaV1.ObjectMeta{Name: "foobar-cert"}})
	suite.ensureObjectDeleted(&coreV1.Secret{ObjectMeta: metaV1.ObjectMeta{Name: "kalm-self-managed-foobar-cert", Namespace: "istio-system"}})
}

func (suite *HttpsCertTestSuite) TestGetEmptyHttpsCertList() {
	rec := suite.NewRequest(http.MethodGet, "/v1alpha1/httpscerts", nil)

	var res []resources.HttpsCert
	rec.BodyAsJSON(&res)

	suite.Equal(200, rec.Code)
}

func (suite *HttpsCertTestSuite) TestCreateHttpsCert() {
	body := `{
  "name":    "foobar-cert",
  "httpsCertIssuer":  "foobar-issuer",
  "domains": ["example.com"]
}`

	rec := suite.NewRequest(http.MethodPost, "/v1alpha1/httpscerts", body)

	var httpsCert resources.HttpsCert
	rec.BodyAsJSON(&httpsCert)

	suite.Equal(201, rec.Code)
	suite.Equal("foobar-cert", httpsCert.Name)

	var res v1alpha1.HttpsCertList
	err := suite.List(&res)
	suite.Nil(err)

	suite.Equal(1, len(res.Items))
	suite.Equal("foobar-cert", res.Items[0].Name)
	//fmt.Println("item:", res.Items[0])
	suite.Equal("foobar-issuer", res.Items[0].Spec.HttpsCertIssuer)
	suite.Equal("example.com", strings.Join(res.Items[0].Spec.Domains, ""))

	// check size & content of cert list
	rec = suite.NewRequest(http.MethodGet, "/v1alpha1/httpscerts", nil)

	var resList []resources.HttpsCertResp
	rec.BodyAsJSON(&resList)

	suite.Equal(200, rec.Code)
	suite.Equal(1, len(resList))
	suite.Equal(string(coreV1.ConditionUnknown), resList[0].Ready)
	suite.Equal(resources.ReasonForNoReadyConditions, resList[0].Reason)
}

func (suite *HttpsCertTestSuite) TestCreateHttpsCertWithoutName() {
	body := `{
  "httpsCertIssuer":  "foobar-issuer",
  "domains": ["example.com"]
}`

	rec := suite.NewRequest(http.MethodPost, "/v1alpha1/httpscerts", body)

	var httpsCert resources.HttpsCert
	rec.BodyAsJSON(&httpsCert)

	suite.Equal(201, rec.Code)
	suite.NotEqual("", httpsCert.Name)

	var res v1alpha1.HttpsCertList
	err := suite.List(&res)
	suite.Nil(err)

	suite.Equal(1, len(res.Items))
	fmt.Println(res.Items[0].Name)
	suite.True(strings.HasPrefix(res.Items[0].Name, "example-com-"))
	suite.Equal("foobar-issuer", res.Items[0].Spec.HttpsCertIssuer)
	suite.Equal("example.com", strings.Join(res.Items[0].Spec.Domains, ""))

	rec = suite.NewRequest(http.MethodDelete, "/v1alpha1/httpscerts/"+httpsCert.Name, nil)
	suite.Equal(200, rec.Code)
}

func (suite *HttpsCertTestSuite) TestCreateHttpsCertWithoutSetIssuer() {
	body := `{
  "name":    "foobar-cert",
  "domains": ["example.com"]
}`

	rec := suite.NewRequest(http.MethodPost, "/v1alpha1/httpscerts", body)

	var httpsCert resources.HttpsCert
	rec.BodyAsJSON(&httpsCert)

	suite.Equal(201, rec.Code)
	suite.Equal("foobar-cert", httpsCert.Name)

	var res v1alpha1.HttpsCertList
	err := suite.List(&res)
	suite.Nil(err)

	suite.Equal(1, len(res.Items))
	suite.Equal("foobar-cert", res.Items[0].Name)
	suite.Equal(controllers.DefaultHTTP01IssuerName, res.Items[0].Spec.HttpsCertIssuer)
	suite.Equal("example.com", strings.Join(res.Items[0].Spec.Domains, ""))
}

const tlsCert = `-----BEGIN CERTIFICATE-----
MIIFVjCCBD6gAwIBAgISBPNCxpUJsb9iD+AX7DviehGrMA0GCSqGSIb3DQEBCwUA
MEoxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MSMwIQYDVQQD
ExpMZXQncyBFbmNyeXB0IEF1dGhvcml0eSBYMzAeFw0yMDA1MjAwMzI5MzNaFw0y
MDA4MTgwMzI5MzNaMBoxGDAWBgNVBAMTD2hlbGxvLmthcHAubGl2ZTCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBAJ48RtSGIUl66BXE5H7TF81dm2JHWxS9
WaPLB9fw+7aE7Q80MqNemxC9919eiLsY43/5vE+oGyqCxy5OA+NjNhqkRyfRtLOy
C+qT30s+bSGVc7iwRqyBSA/1tVjvlio+bD3jmiKP8G4fX0MswJmUIhUqDjrgcz73
7WCn0SxfJrxRVihgQ0BYdwn7rhXd9owQK5KIYG80a/oLwnsplJCzI3MeDzhLz/oK
pcaPI8qoLH4Bxb7Od/tKODpp80ub6c4x+M+qI62Goo50+Sm6vwVzc8CsSlz2lGDN
608tBWDZ6HJrn0ogBa/qUFdSFXrQcpeVNVi+/suT/+wGJ1KtiDInM0ECAwEAAaOC
AmQwggJgMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYB
BQUHAwIwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQU88bxblZUQdX7RYMsUnKxUTwK
Z04wHwYDVR0jBBgwFoAUqEpqYwR93brm0Tm3pkVl7/Oo7KEwbwYIKwYBBQUHAQEE
YzBhMC4GCCsGAQUFBzABhiJodHRwOi8vb2NzcC5pbnQteDMubGV0c2VuY3J5cHQu
b3JnMC8GCCsGAQUFBzAChiNodHRwOi8vY2VydC5pbnQteDMubGV0c2VuY3J5cHQu
b3JnLzAaBgNVHREEEzARgg9oZWxsby5rYXBwLmxpdmUwTAYDVR0gBEUwQzAIBgZn
gQwBAgEwNwYLKwYBBAGC3xMBAQEwKDAmBggrBgEFBQcCARYaaHR0cDovL2Nwcy5s
ZXRzZW5jcnlwdC5vcmcwggEEBgorBgEEAdZ5AgQCBIH1BIHyAPAAdgDwlaRZ8gDR
gkAQLS+TiI6tS/4dR+OZ4dA0prCoqo6ycwAAAXIwWAAmAAAEAwBHMEUCIQD/weuk
7dWqw7iswofV7vt4ANxIvVFKfynHik1tDWGX2QIgUZwvdLxNjduE15kPB5G3zpbp
7I8Y2VIWIgxyZIGR3BEAdgCyHgXMi6LNiiBOh2b5K7mKJSBna9r6cOeySVMt74uQ
XgAAAXIwWAAVAAAEAwBHMEUCIGJwq3BvFcWn8CwRwXsMkOR2FUKAV1XcDwcJsbJa
jFsgAiEA5dqDJ0oL2V2ItThyNQRGTD7WvVKjghCL+EIO5yaYZaswDQYJKoZIhvcN
AQELBQADggEBAJJ8mKQ+IyuFlOMijD5RhU3U7l8rR5R/f9ITRUK5Q3NgkmhvNG8t
wBCcnVr3nNnKFYloLJ0rBSPyqs/nE9KljHzhZoootkP8PfXHe7A6FOR8xohzqHR0
u54xd1p+jTluVOE+Ofwa32VZ4VkIUIezyoSpLz1xk00tVtFlIrBn1Bk2vskTA5XK
znkm2KnVBuj75tteXjjMthi+yKW1Bfd1I2mCuSz8sylKyXx+2sX6YVR5o1+NamBi
7mK92Uhdb4Zq21RpDNYWrITAjVIunofNSgGjFu165ZGvPCMG0DwhvFnzWb97dsB5
fZLKGiQmUq6JTawO7JIZDdVDK7zZQTsEQG8=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIEkjCCA3qgAwIBAgIQCgFBQgAAAVOFc2oLheynCDANBgkqhkiG9w0BAQsFADA/
MSQwIgYDVQQKExtEaWdpdGFsIFNpZ25hdHVyZSBUcnVzdCBDby4xFzAVBgNVBAMT
DkRTVCBSb290IENBIFgzMB4XDTE2MDMxNzE2NDA0NloXDTIxMDMxNzE2NDA0Nlow
SjELMAkGA1UEBhMCVVMxFjAUBgNVBAoTDUxldCdzIEVuY3J5cHQxIzAhBgNVBAMT
GkxldCdzIEVuY3J5cHQgQXV0aG9yaXR5IFgzMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAnNMM8FrlLke3cl03g7NoYzDq1zUmGSXhvb418XCSL7e4S0EF
q6meNQhY7LEqxGiHC6PjdeTm86dicbp5gWAf15Gan/PQeGdxyGkOlZHP/uaZ6WA8
SMx+yk13EiSdRxta67nsHjcAHJyse6cF6s5K671B5TaYucv9bTyWaN8jKkKQDIZ0
Z8h/pZq4UmEUEz9l6YKHy9v6Dlb2honzhT+Xhq+w3Brvaw2VFn3EK6BlspkENnWA
a6xK8xuQSXgvopZPKiAlKQTGdMDQMc2PMTiVFrqoM7hD8bEfwzB/onkxEz0tNvjj
/PIzark5McWvxI0NHWQWM6r6hCm21AvA2H3DkwIDAQABo4IBfTCCAXkwEgYDVR0T
AQH/BAgwBgEB/wIBADAOBgNVHQ8BAf8EBAMCAYYwfwYIKwYBBQUHAQEEczBxMDIG
CCsGAQUFBzABhiZodHRwOi8vaXNyZy50cnVzdGlkLm9jc3AuaWRlbnRydXN0LmNv
bTA7BggrBgEFBQcwAoYvaHR0cDovL2FwcHMuaWRlbnRydXN0LmNvbS9yb290cy9k
c3Ryb290Y2F4My5wN2MwHwYDVR0jBBgwFoAUxKexpHsscfrb4UuQdf/EFWCFiRAw
VAYDVR0gBE0wSzAIBgZngQwBAgEwPwYLKwYBBAGC3xMBAQEwMDAuBggrBgEFBQcC
ARYiaHR0cDovL2Nwcy5yb290LXgxLmxldHNlbmNyeXB0Lm9yZzA8BgNVHR8ENTAz
MDGgL6AthitodHRwOi8vY3JsLmlkZW50cnVzdC5jb20vRFNUUk9PVENBWDNDUkwu
Y3JsMB0GA1UdDgQWBBSoSmpjBH3duubRObemRWXv86jsoTANBgkqhkiG9w0BAQsF
AAOCAQEA3TPXEfNjWDjdGBX7CVW+dla5cEilaUcne8IkCJLxWh9KEik3JHRRHGJo
uM2VcGfl96S8TihRzZvoroed6ti6WqEBmtzw3Wodatg+VyOeph4EYpr/1wXKtx8/
wApIvJSwtmVi4MFU5aMqrSDE6ea73Mj2tcMyo5jMd6jmeWUHK8so/joWUoHOUgwu
X4Po1QYz+3dszkDqMp4fklxBwXRsW10KXzPMTZ+sOPAveyxindmjkW8lGy+QsRlG
PfZ+G6Z6h7mjem0Y+iWlkYcV4PIWL1iwBi8saCbGS5jN2p8M+X+Q7UNKEkROb3N6
KOqkqm57TH2H3eDJAkSnh6/DNFu0Qg==
-----END CERTIFICATE-----`

func (suite *HttpsCertTestSuite) TestUploadHttpsCert() {

	bodyMap := map[string]interface{}{
		"name":                      "foobar-cert",
		"isSelfManaged":             true,
		"selfManagedCertContent":    tlsCert,
		"selfManagedCertPrivateKey": "",
	}

	bodyBytes, _ := json.Marshal(bodyMap)
	body := string(bodyBytes)

	//body = `{"name":"hydro","managedType":"selfManaged","selfManagedCertContent":"-----BEGIN CERTIFICATE-----\nMIIBnTCCAQYCCQCenBNa6J7xKDANBgkqhkiG9w0BAQUFADATMREwDwYDVQQDDAho\neWRyby5pbzAeFw0yMDA1MzAwNjU5NDdaFw0zMDA1MjgwNjU5NDdaMBMxETAPBgNV\nBAMMCGh5ZHJvLmlvMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC3Zli6e0Tc\ndKgPr5L72Kx4NbathnS5qoN6yR5Couql24l3R7ex3l1rxUt9XOrvtf6AxuZ6Uvkl\nutG6XdMF0ER/X95S/ExCif9Qcj7jephuWikXEYKr81ROyd66fiVmwZEtxsRcn+Bk\nuaOm1u8z4z5n9f0bHCyBeMS2c8GK7dBruwIDAQABMA0GCSqGSIb3DQEBBQUAA4GB\nAB1M37ckPSVzS4YrLss5G2ua/yYP8U/d504Oma3+AQ5HIUJIa0SVaeb07X2AqxRW\n9DCMPBKgBH6H5Dk9asJQKNXoV3597OexPqgf71YnlI5Fe65dWzFwgka6vzrVj96Z\nFLhV1tAWV3qTgAQtPuRq5WCXR0gpzXGsAzjoq8v91nM1\n-----END CERTIFICATE-----\n","selfManagedCertPrivateKey":"-----BEGIN RSA PRIVATE KEY-----\nMIICXgIBAAKBgQCmBL0loX5LGPUt4lLwKsRy0RP1/PsIek7cL2ypWiySxpmb/J5P\nuQh7ioD9GhpDz3P0KmP/70IrxsGI2NuFkrLvEEu/as/rGQdePY5MSvE5IAmw2+q3\nKxPGqslSMyUu/9BcA673wsHGmFx/WQN1pcH8Uvn5OpVtBe0FYGfmbx4NdQIDAQAB\nAoGBAJ2Jj6ce2MYQV9ADnOoz2xiM1+Er/1ZrIVwFUJpfbD6XYHpJvDwwfCEnLAGA\n/H1uZTzD4QKjGmqJ/2qJRb45l7qFQAJQus/QjlMSiYpmyB/K4WTt+xQPguDBqrJa\nwgtaqEz3pajHzis4g2HlzvgfNNDdDzPb//SgWiHNFal9o9K9AkEA2LxLMUiSIqbr\nwWdQ7y86e6UJ7bvs9sef2s6+owucntt0pQR7kqzdEOx0YnAM7a+QJ6aAp764cxC/\nqwgP87Z9zwJBAMQYTd5PntR5CmVpCQP+AJ9GjqdpguZ2l2mXfhkBFzG2k+ElzhIs\n9vh+kXXUR/BR4cgkqrpdrZyG5VMCU04QdXsCQQDHbfuCjdRqLk4g2ICQllSqEWLh\nblrNeUKOhE7GmJi1VBow/b73aDjCqdEuAwziaPmsgnk/4S64MmErb2++Qd3hAkA3\nWGi54xz/+P2vK5sIzmrciWx/4a65XyVS/xfu9LjYJiDMvf/Qb3JxJj/l0rlPIb0o\n4PhyyzfR4tzctd5PPSP7AkEAnubKj9ZR9qBE6YKm3A5zy5vUAEYgRGE/vMcS6qa+\n0Hj2yaFlbFHQ/QZqzOSYe7uUPQ+daOa4YF9ecyB+lvgM5Q==\n-----END RSA PRIVATE KEY-----\n","isSelfManaged":true}`
	rec := suite.NewRequest(http.MethodPost, "/v1alpha1/httpscerts/upload", body)

	var httpsCert resources.HttpsCert
	rec.BodyAsJSON(&httpsCert)

	suite.Equal(201, rec.Code)
	suite.Equal("foobar-cert", httpsCert.Name)

	var res v1alpha1.HttpsCertList
	err := suite.List(&res)
	suite.Nil(err)

	suite.Equal(1, len(res.Items))
	suite.Equal("foobar-cert", res.Items[0].Name)
	//fmt.Println("item:", res.Items[0])
	suite.Equal(true, res.Items[0].Spec.IsSelfManaged)
	suite.Equal("", res.Items[0].Spec.HttpsCertIssuer)
	suite.Equal("hello.kapp.live", strings.Join(res.Items[0].Spec.Domains, ""))
	suite.Equal("kalm-self-managed-foobar-cert", res.Items[0].Spec.SelfManagedCertSecretName)

	// sec
	var sec coreV1.Secret
	err = suite.Get("istio-system", "kalm-self-managed-foobar-cert", &sec)
	suite.Nil(err)
	suite.Equal(sec.Data["tls.key"], []byte(""))
	suite.Equal(sec.Data["tls.crt"], []byte(tlsCert))

	rec = suite.NewRequest(http.MethodGet, "/v1alpha1/httpscerts", nil)
	var resList []resources.HttpsCert
	rec.BodyAsJSON(&resList)
	suite.Equal(200, rec.Code)
	suite.Equal(1, len(resList))
	suite.Equal("hello.kapp.live", strings.Join(resList[0].Domains, ""))
}

func (suite *HttpsCertTestSuite) TestUpdateSelfManagedHttpsCert() {

	// upload
	bodyBytes, _ := json.Marshal(map[string]interface{}{
		"name":                      "foobar-cert",
		"isSelfManaged":             true,
		"selfManagedCertContent":    tlsCert,
		"selfManagedCertPrivateKey": "",
	})
	rec := suite.NewRequest(http.MethodPost, "/v1alpha1/httpscerts/upload", string(bodyBytes))
	suite.Equal(201, rec.Code)

	// update
	updateBodyBytes, _ := json.Marshal(map[string]interface{}{
		"name":                      "foobar-cert",
		"isSelfManaged":             true,
		"selfManagedCertContent":    tlsCert,
		"selfManagedCertPrivateKey": "updatedPrvKey",
	})
	rec = suite.NewRequest(http.MethodPut, "/v1alpha1/httpscerts/foobar-cert", string(updateBodyBytes))
	suite.Equal(200, rec.Code)

	var res v1alpha1.HttpsCertList
	err := suite.List(&res)
	suite.Nil(err)

	suite.Equal(1, len(res.Items))
	suite.Equal("foobar-cert", res.Items[0].Name)
	suite.Equal(true, res.Items[0].Spec.IsSelfManaged)
	suite.Equal("", res.Items[0].Spec.HttpsCertIssuer)
	suite.Equal("hello.kapp.live", strings.Join(res.Items[0].Spec.Domains, ""))
	suite.Equal("kalm-self-managed-foobar-cert", res.Items[0].Spec.SelfManagedCertSecretName)

	// sec
	var sec coreV1.Secret
	err = suite.Get("istio-system", "kalm-self-managed-foobar-cert", &sec)
	suite.Nil(err)
	suite.Equal(sec.Data["tls.key"], []byte("updatedPrvKey"))
}

func (suite *HttpsCertTestSuite) TestUpdateAutoManagedHttpsCert() {
	body := `{
  "name":    "foobar-cert",
  "httpsCertIssuer":  "foobar-issuer",
  "domains": ["example.com"]
}`

	rec := suite.NewRequest(http.MethodPost, "/v1alpha1/httpscerts", body)

	var httpsCert resources.HttpsCert
	rec.BodyAsJSON(&httpsCert)

	suite.Equal(201, rec.Code)
	suite.Equal("foobar-cert", httpsCert.Name)

	var res v1alpha1.HttpsCertList
	err := suite.List(&res)
	suite.Nil(err)

	suite.Equal(1, len(res.Items))
	suite.Equal("foobar-cert", res.Items[0].Name)
	suite.Equal("foobar-issuer", res.Items[0].Spec.HttpsCertIssuer)
	suite.Equal("example.com", strings.Join(res.Items[0].Spec.Domains, ""))

	body = `{
  "name":    "foobar-cert",
  "httpsCertIssuer":  "foobar-issuer2",
  "domains": ["example.com2"]
}`
	rec = suite.NewRequest(http.MethodPut, "/v1alpha1/httpscerts/foobar-cert", body)

	rec.BodyAsJSON(&httpsCert)

	suite.Equal(200, rec.Code)
	suite.Equal("foobar-cert", httpsCert.Name)

	err = suite.List(&res)
	suite.Nil(err)

	suite.Equal(1, len(res.Items))
	suite.Equal("foobar-cert", res.Items[0].Name)
	suite.Equal("foobar-issuer2", res.Items[0].Spec.HttpsCertIssuer)
	suite.Equal("example.com2", strings.Join(res.Items[0].Spec.Domains, ""))
}

func (suite *HttpsCertTestSuite) TestDeleteHttpsCert() {
	body := `{
  "name":    "foobar-cert",
  "httpsCertIssuer":  "foobar-issuer",
  "domains": ["example.com"]
}`

	rec := suite.NewRequest(http.MethodPost, "/v1alpha1/httpscerts", body)

	var httpsCert resources.HttpsCert
	rec.BodyAsJSON(&httpsCert)

	suite.Equal(201, rec.Code)
	suite.Equal("foobar-cert", httpsCert.Name)

	var res v1alpha1.HttpsCertList
	err := suite.List(&res)
	suite.Nil(err)

	suite.Equal(1, len(res.Items))
	suite.Equal("foobar-cert", res.Items[0].Name)
	//fmt.Println("item:", res.Items[0])
	suite.Equal("foobar-issuer", res.Items[0].Spec.HttpsCertIssuer)
	suite.Equal("example.com", strings.Join(res.Items[0].Spec.Domains, ""))

	rec = suite.NewRequest(http.MethodDelete, "/v1alpha1/httpscerts/foobar-cert", nil)
	suite.Equal(200, rec.Code)

	err = suite.List(&res)
	suite.Nil(err)
	suite.Equal(0, len(res.Items))
}
