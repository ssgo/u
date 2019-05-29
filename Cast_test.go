package u_test

import (
	"github.com/ssgo/u"
	"testing"
)

func TestFixUpperCase1(t *testing.T) {
	buf := []byte("{\"FLAG\":\"321\",\"app\":\"\",\"authLevel\":0,\"clientId\":\"\",\"clientIp\":\"[\",\"FromApp\":\"\",\"FromNode\":\"\",\"Host\":\"localhost:8080\",\"LogTime\":1557033007.432081,\"LogType\":\"request\",\"Method\":\"GET\",\"Node\":\"10.59.5.226:8080\",\"Path\":\"/\",\"Priority\":0,\"Proto\":\"1.1\",\"RequestData\":{},\"RequestHeaders\":{\"Cookie\":\"Phpstorm-52af694c=1c7e6019-5e5e-4085-a376-118790b33e5c\",\"If-Modified-Since\":\"Thu, 18 Apr 2019 10:54:43 GMT\",\"User-Agent\":\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36\",\"X-Host\":\"localhost:8080\",\"X-Real-Ip\":\"[\",\"X-Request-Id\":\"o6Obm1EtsaISY6mth45n7\",\"X-Scheme\":\"http\"},\"requestId\":\"o6Obm1EtsaISY6mth45n7\",\"responseCode\":304,\"responseData\":null,\"responseDataLength\":1749,\"responseHeaders\":{\"Last-Modified\":\"Thu, 18 Apr 2019 10:54:43 GMT\"},\"scheme\":\"http\",\"serverId\":\"nfWX0sr1ewF\",\"sessionId\":\"\",\"traceId\":\"o6Obm1EtsaISY6mth45n7\",\"usedTime\":0.366,\"type\":\"STATIC\"}")
	u.FixUpperCase(buf, []string{"Header"})
	//to := bytes.Buffer{}
	//_ = json.Indent(&to, buf, "", "  ")
	//fmt.Println(to.String())

	//fmt.Println(string(buf))
	rights := "{\"FLAG\":\"321\",\"app\":\"\",\"authLevel\":0,\"clientId\":\"\",\"clientIp\":\"[\",\"fromApp\":\"\",\"fromNode\":\"\",\"host\":\"localhost:8080\",\"logTime\":1557033007.432081,\"logType\":\"request\",\"method\":\"GET\",\"node\":\"10.59.5.226:8080\",\"path\":\"/\",\"priority\":0,\"proto\":\"1.1\",\"requestData\":{},\"requestHeaders\":{\"Cookie\":\"Phpstorm-52af694c=1c7e6019-5e5e-4085-a376-118790b33e5c\",\"If-Modified-Since\":\"Thu, 18 Apr 2019 10:54:43 GMT\",\"User-Agent\":\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36\",\"X-Host\":\"localhost:8080\",\"X-Real-Ip\":\"[\",\"X-Request-Id\":\"o6Obm1EtsaISY6mth45n7\",\"X-Scheme\":\"http\"},\"requestId\":\"o6Obm1EtsaISY6mth45n7\",\"responseCode\":304,\"responseData\":null,\"responseDataLength\":1749,\"responseHeaders\":{\"Last-Modified\":\"Thu, 18 Apr 2019 10:54:43 GMT\"},\"scheme\":\"http\",\"serverId\":\"nfWX0sr1ewF\",\"sessionId\":\"\",\"traceId\":\"o6Obm1EtsaISY6mth45n7\",\"usedTime\":0.366,\"type\":\"STATIC\"}"
	if string(buf) != rights {
		t.Error("FixUpperCase1 failed")
	}
}

func TestFixUpperCase2(t *testing.T) {
	buf := []byte("{\"FLAG\":\"321\",\"App\":\"\",\"Info\":\"stopping router\",\"LogTime\":1557057463.721667,\"LogType\":\"server\",\"Node\":\"10.59.5.226:18811\",\"Proto\":\"http\",\"StartTime\":1557057463.6995,\"TraceId\":\"g6aXzbekTGZ\",\"Weight\":1}")
	u.FixUpperCase(buf, nil)
	//to := bytes.Buffer{}
	//_ = json.Indent(&to, buf, "", "  ")
	//fmt.Println(to.String())

	//fmt.Println(string(buf))
	rights := "{\"FLAG\":\"321\",\"app\":\"\",\"info\":\"stopping router\",\"logTime\":1557057463.721667,\"logType\":\"server\",\"node\":\"10.59.5.226:18811\",\"proto\":\"http\",\"startTime\":1557057463.6995,\"traceId\":\"g6aXzbekTGZ\",\"weight\":1}"
	if string(buf) != rights {
		t.Error("FixUpperCase2 failed")
	}
}
