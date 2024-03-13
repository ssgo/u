package u_test

import (
	"github.com/ssgo/u"
	"strings"
	"testing"
	"time"
)

func TestFixUpperCase1(t *testing.T) {
	buf := []byte("{\"FLAG\":\"321\",\"app\":\"\",\"authLevel\":0,\"clientId\":\"\",\"clientIp\":\"[\",\"FromApp\":\"\",\"FromNode\":\"\",\"Host\":\"localhost:8080\",\"LogTime\":1557033007.432081,\"LogType\":\"request\",\"Method\":\"GET\",\"Node\":\"10.59.5.226:8080\",\"Path\":\"/\",\"Priority\":0,\"Proto\":\"1.1\",\"RequestData\":{},\"RequestHeaders\":{\"Cookie\":\"Phpstorm-52af694c=1c7e6019-5e5e-4085-a376-118790b33e5c\",\"If-Modified-Since\":\"Thu, 18 Apr 2019 10:54:43 GMT\",\"User-Agent\":\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36\",\"X-Host\":\"localhost:8080\",\"X-Real-Ip\":\"[\",\"X-Request-Id\":\"o6Obm1EtsaISY6mth45n7\",\"X-Scheme\":\"http\"},\"requestId\":\"o6Obm1EtsaISY6mth45n7\",\"responseCode\":304,\"responseData\":null,\"responseDataLength\":1749,\"responseHeaders\":{\"Last-Modified\":\"Thu, 18 Apr 2019 10:54:43 GMT\"},\"scheme\":\"http\",\"serverId\":\"nfWX0sr1ewF\",\"sessionId\":\"\",\"traceId\":\"o6Obm1EtsaISY6mth45n7\",\"usedTime\":0.366,\"type\":\"STATIC\"}")
	u.FixUpperCase(buf, []string{"RequestHeaders.", "responseHeaders."})
	//to := bytes.Buffer{}
	//_ = json.Indent(&to, buf, "", "  ")
	//fmt.Println(to.String())

	//fmt.Println(string(buf))
	rights := "{\"FLAG\":\"321\",\"app\":\"\",\"authLevel\":0,\"clientId\":\"\",\"clientIp\":\"[\",\"fromApp\":\"\",\"fromNode\":\"\",\"host\":\"localhost:8080\",\"logTime\":1557033007.432081,\"logType\":\"request\",\"method\":\"GET\",\"node\":\"10.59.5.226:8080\",\"path\":\"/\",\"priority\":0,\"proto\":\"1.1\",\"requestData\":{},\"requestHeaders\":{\"Cookie\":\"Phpstorm-52af694c=1c7e6019-5e5e-4085-a376-118790b33e5c\",\"If-Modified-Since\":\"Thu, 18 Apr 2019 10:54:43 GMT\",\"User-Agent\":\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36\",\"X-Host\":\"localhost:8080\",\"X-Real-Ip\":\"[\",\"X-Request-Id\":\"o6Obm1EtsaISY6mth45n7\",\"X-Scheme\":\"http\"},\"requestId\":\"o6Obm1EtsaISY6mth45n7\",\"responseCode\":304,\"responseData\":null,\"responseDataLength\":1749,\"responseHeaders\":{\"Last-Modified\":\"Thu, 18 Apr 2019 10:54:43 GMT\"},\"scheme\":\"http\",\"serverId\":\"nfWX0sr1ewF\",\"sessionId\":\"\",\"traceId\":\"o6Obm1EtsaISY6mth45n7\",\"usedTime\":0.366,\"type\":\"STATIC\"}"
	if string(buf) != rights {
		t.Error("FixUpperCase1 failed", rights, string(buf), ".")
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

func TestFixUpperCase3(t *testing.T) {
	type Item struct {
		Name string
		Set1 map[string]string
		Set2 map[string]string `keepKey`
		Set3 map[string]string `keepSubKey`
		Set4 map[string]string `keepKey keepSubKey`
	}
	type List struct {
		DefaultSet1 map[string]string
		DefaultSet2 map[string]string `keepKey`
		DefaultSet3 map[string]string `keepSubKey`
		DefaultSet4 map[string]string `keepKey,keepSubKey`
		List1       []Item
		List2       []Item `keepKey`
		List3       []Item `keepSubKey`
		List4       []Item `keepKey keepSubKey`
	}
	set := map[string]string{
		"abc1": "1",
		"Abc2": "2",
	}
	item := Item{
		Name: "测试",
		Set1: set,
		Set2: set,
		Set3: set,
		Set4: set,
	}
	list := List{
		DefaultSet1: set,
		DefaultSet2: set,
		DefaultSet3: set,
		DefaultSet4: set,
		List1:       []Item{item},
		List2:       []Item{item},
		List3:       []Item{item},
		List4:       []Item{item},
	}

	excludeKeys := u.MakeExcludeUpperKeys(list, "")
	//fmt.Println(u.JsonP(excludeKeys), ".")

	buf := u.JsonBytes(list)
	u.FixUpperCase(buf, excludeKeys)
	//fmt.Println(string(buf))
	str := string(buf)

	if !strings.Contains(str, `"defaultSet1":{"abc2":"2","abc1":"1"}`) {
		t.Error("FixUpperCase3 failed on defaultSet1")
	}

	if !strings.Contains(str, `"DefaultSet2":{"abc2":"2","abc1":"1"}`) {
		t.Error("FixUpperCase3 failed on defaultSet2")
	}

	if !strings.Contains(str, `"defaultSet3":{"Abc2":"2","abc1":"1"}`) {
		t.Error("FixUpperCase3 failed on defaultSet3")
	}

	if !strings.Contains(str, `"DefaultSet4":{"Abc2":"2","abc1":"1"}`) {
		t.Error("FixUpperCase3 failed on defaultSet4")
	}

	if !strings.Contains(str, `"list1":[{"name":"测试","set1":{"abc2":"2","abc1":"1"},"Set2":{"abc2":"2","abc1":"1"},"set3":{"Abc2":"2","abc1":"1"},"Set4":{"Abc2":"2","abc1":"1"}}]`) {
		t.Error("FixUpperCase3 failed on list1")
	}

	if !strings.Contains(str, `"List2":[{"name":"测试","set1":{"abc2":"2","abc1":"1"},"Set2":{"abc2":"2","abc1":"1"},"set3":{"Abc2":"2","abc1":"1"},"Set4":{"Abc2":"2","abc1":"1"}}]`) {
		t.Error("FixUpperCase3 failed on list2")
	}

	if !strings.Contains(str, `"list3":[{"Name":"测试","Set1":{"Abc2":"2","abc1":"1"},"Set2":{"Abc2":"2","abc1":"1"},"Set3":{"Abc2":"2","abc1":"1"},"Set4":{"Abc2":"2","abc1":"1"}}]`) {
		t.Error("FixUpperCase3 failed on list3")
	}

	if !strings.Contains(str, `"List4":[{"Name":"测试","Set1":{"Abc2":"2","abc1":"1"},"Set2":{"Abc2":"2","abc1":"1"},"Set3":{"Abc2":"2","abc1":"1"},"Set4":{"Abc2":"2","abc1":"1"}}]`) {
		t.Error("FixUpperCase3 failed on list4")
	}
}

func TestDuration(t *testing.T) {
	sets := map[string]time.Duration{
		"100m":   100 * time.Minute,
		"100s":   100 * time.Second,
		"100ms":  100 * time.Millisecond,
		"100":    100 * time.Millisecond,
		"100us":  100 * time.Microsecond,
		"100ns":  100 * time.Nanosecond,
		"100xxx": 0,
		"hdasds": 0,
	}

	for set, check := range sets {
		if u.Duration(set) != check {
			t.Error("time", set, "is", int64(u.Duration(set)), "!=", int64(check))
		}
	}
}
