package u_test

import (
	"github.com/ssgo/u"
	"reflect"
	"testing"
)

func TestFixUpperCase1(t *testing.T) {
	buf := []byte("{\"app\":\"\",\"authLevel\":0,\"clientId\":\"\",\"clientIp\":\"[\",\"FromApp\":\"\",\"FromNode\":\"\",\"Host\":\"localhost:8080\",\"LogTime\":1557033007.432081,\"LogType\":\"request\",\"Method\":\"GET\",\"Node\":\"10.59.5.226:8080\",\"Path\":\"/\",\"Priority\":0,\"Proto\":\"1.1\",\"RequestData\":{},\"RequestHeaders\":{\"Cookie\":\"Phpstorm-52af694c=1c7e6019-5e5e-4085-a376-118790b33e5c\",\"If-Modified-Since\":\"Thu, 18 Apr 2019 10:54:43 GMT\",\"User-Agent\":\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36\",\"X-Host\":\"localhost:8080\",\"X-Real-Ip\":\"[\",\"X-Request-Id\":\"o6Obm1EtsaISY6mth45n7\",\"X-Scheme\":\"http\"},\"requestId\":\"o6Obm1EtsaISY6mth45n7\",\"responseCode\":304,\"responseData\":null,\"responseDataLength\":1749,\"responseHeaders\":{\"Last-Modified\":\"Thu, 18 Apr 2019 10:54:43 GMT\"},\"scheme\":\"http\",\"serverId\":\"nfWX0sr1ewF\",\"sessionId\":\"\",\"traceId\":\"o6Obm1EtsaISY6mth45n7\",\"usedTime\":0.366,\"type\":\"STATIC\"}")
	u.FixUpperCase(buf, []string{"Header"})
	//to := bytes.Buffer{}
	//_ = json.Indent(&to, buf, "", "  ")
	//fmt.Println(to.String())

	//fmt.Println(string(buf))
	rights := "{\"app\":\"\",\"authLevel\":0,\"clientId\":\"\",\"clientIp\":\"[\",\"fromApp\":\"\",\"fromNode\":\"\",\"host\":\"localhost:8080\",\"logTime\":1557033007.432081,\"logType\":\"request\",\"method\":\"GET\",\"node\":\"10.59.5.226:8080\",\"path\":\"/\",\"priority\":0,\"proto\":\"1.1\",\"requestData\":{},\"requestHeaders\":{\"Cookie\":\"Phpstorm-52af694c=1c7e6019-5e5e-4085-a376-118790b33e5c\",\"If-Modified-Since\":\"Thu, 18 Apr 2019 10:54:43 GMT\",\"User-Agent\":\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36\",\"X-Host\":\"localhost:8080\",\"X-Real-Ip\":\"[\",\"X-Request-Id\":\"o6Obm1EtsaISY6mth45n7\",\"X-Scheme\":\"http\"},\"requestId\":\"o6Obm1EtsaISY6mth45n7\",\"responseCode\":304,\"responseData\":null,\"responseDataLength\":1749,\"responseHeaders\":{\"Last-Modified\":\"Thu, 18 Apr 2019 10:54:43 GMT\"},\"scheme\":\"http\",\"serverId\":\"nfWX0sr1ewF\",\"sessionId\":\"\",\"traceId\":\"o6Obm1EtsaISY6mth45n7\",\"usedTime\":0.366,\"type\":\"STATIC\"}"
	if string(buf) != rights {
		t.Error("FixUpperCase1 failed")
	}
}

func TestFixUpperCase2(t *testing.T) {
	buf := []byte("{\"App\":\"\",\"Info\":\"stopping router\",\"LogTime\":1557057463.721667,\"LogType\":\"server\",\"Node\":\"10.59.5.226:18811\",\"Proto\":\"http\",\"StartTime\":1557057463.6995,\"TraceId\":\"g6aXzbekTGZ\",\"Weight\":1}")
	u.FixUpperCase(buf, nil)
	//to := bytes.Buffer{}
	//_ = json.Indent(&to, buf, "", "  ")
	//fmt.Println(to.String())

	//fmt.Println(string(buf))
	rights := "{\"app\":\"\",\"info\":\"stopping router\",\"logTime\":1557057463.721667,\"logType\":\"server\",\"node\":\"10.59.5.226:18811\",\"proto\":\"http\",\"startTime\":1557057463.6995,\"traceId\":\"g6aXzbekTGZ\",\"weight\":1}"
	if string(buf) != rights {
		t.Error("FixUpperCase2 failed")
	}
}

func TestBaseConvert(t *testing.T) {
	froms := []interface{}{
		int(99),
		int8(99),
		int16(99),
		int32(99),
		int64(99),
		uint(99),
		uint8(99),
		uint16(99),
		uint32(99),
		uint64(99),
		float32(99.99),
		float64(99.99),
		true,
		"99",
	}

	tos := []interface{}{
		int(88),
		int8(88),
		int16(88),
		int32(88),
		int64(88),
		uint(88),
		uint8(88),
		uint16(88),
		uint32(88),
		uint64(88),
		float32(88.88),
		float64(88.88),
		false,
		"88",
	}

	for j := 0; j < len(tos); j++ {
		for i := 0; i < len(froms); i++ {
			u.Convert(froms[i], &tos[j])
			if froms[i] != tos[j] {
				t.Error("convert not match", reflect.TypeOf(froms[i]), reflect.TypeOf(tos[j]), froms[j], tos[j])
			}
		}
	}

	toString := ""
	for i := 0; i < len(froms); i++ {
		u.Convert(froms[i], &toString)
		if u.String(froms[i]) != toString {
			t.Error("convert to string not match", reflect.TypeOf(froms[i]), froms[i], toString)
		}
	}

	var toInt int16
	for i := 0; i < len(froms); i++ {
		u.Convert(froms[i], &toInt)
		if int16(u.Int(froms[i])) != toInt {
			t.Error("convert to int not match", reflect.TypeOf(froms[i]), froms[i], toInt)
		}
	}
}

func TestConvertMapToStruct(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	from := map[string]interface{}{
		"aaa":  "111",
		"bbb":  222,
		"bbb2": 222,
		"bbb3": 222,
		"ccc":  333.333,
		"ccc2": []string{"1", "2", "3.0"},
		"ccc3": []int{1, 2, 3},
		"ccc4": []interface{}{"1", 2, 3.0},
		"ccc5": []map[string]interface{}{
			{
				"name": 13,
			},
			{
				"name": "ALin",
				"age":  22,
			},
		},
	}

	type CCC struct {
		Ccc  string
		Ccc2 []int
		Ccc3 []interface{}
		Ccc4 []*float32
		Ccc5 []*User
	}
	type BBB struct {
		Bbb  *string
		BBB2 []byte
		Bbb3 ****string
		CCC
	}
	to := struct {
		Aaa int
		BBB
	}{}

	u.Convert(from, &to)

	if to.Aaa != 111 || string(to.BBB2) != "222" || to.Ccc != "333.333" {
		t.Error("test convert Map to Struct 1", to)
	}

	if (*to.Bbb) != "222" || (****to.Bbb3) != "222" {
		t.Error("test convert Map to Struct 2", to)
	}

	if len(to.Ccc2) < 3 || to.Ccc2[0] != 1 || to.Ccc2[1] != 2 || to.Ccc2[2] != 3 {
		t.Error("test convert Slice to Slice 1", to.Ccc2)
	}

	if len(to.Ccc3) < 3 || to.Ccc3[0] != 1 || to.Ccc3[1] != 2 || to.Ccc3[2] != 3 {
		t.Error("test convert Slice to Slice 2", to.Ccc3)
	}

	if len(to.Ccc4) < 3 || *to.Ccc4[0] != 1 || *to.Ccc4[1] != 2 || *to.Ccc4[2] != 3 {
		t.Error("test convert Slice to Slice 3", to.Ccc4)
	}

	if len(to.Ccc5) < 2 || to.Ccc5[0].Name != "13" || to.Ccc5[1].Name != "ALin" || to.Ccc5[1].Age != 22 {
		t.Error("test convert Slice to Slice 3", to.Ccc5)
	}
}

func TestConvertMapToMap(t *testing.T) {
	from := map[string]interface{}{
		"aaa":  "111",
		"bbb":  222,
		"bbb2": 222,
		"bbb3": 222,
		"ccc":  333.333,
		"ccc2": []string{"1", "2", "3.0"},
		"ccc3": []int{1, 2, 3},
		"ccc4": []interface{}{"1", 2, 3.0},
	}

	to := map[string]interface{}{}
	u.Convert(from, &to)
	if to["aaa"] != "111" || to["bbb2"] != 222 || u.String(to["ccc"]) != "333.333" {
		t.Error("test convert Map to Map 1", to)
	}

	ccc2 := to["ccc2"].([]string)
	if len(ccc2) < 3 || ccc2[0] != "1" || ccc2[1] != "2" || ccc2[2] != "3.0" {
		t.Error("test convert Map to Map 2", ccc2)
	}

	ccc3 := to["ccc3"].([]int)
	if len(ccc3) < 3 || ccc3[0] != 1 || ccc3[1] != 2 || ccc3[2] != 3 {
		t.Error("test convert Map to Map 3", ccc3)
	}

	ccc4 := to["ccc4"].([]interface{})
	if len(ccc4) < 3 || ccc4[0] != "1" || ccc4[1] != 2 || ccc4[2] != 3.0 {
		t.Error("test convert Map to Map 3", ccc4)
	}
}

func TestConvertStructToMap(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	from := []User{
		{
			Name: "Andy",
			Age:  11,
		},
		{
			Name: "Kitty",
			Age:  22,
		},
	}

	to := make([]map[string]interface{}, 0)

	r := u.Convert(from, to)
	if v, ok := r.([]map[string]interface{}); ok {
		if len(v) < 2 || v[0]["Name"] != "Andy" || v[1]["Name"] != "Kitty" || v[1]["Age"] != 22 {
			t.Error("test convert Struct to Map 2", v)
		}
	} else {
		t.Error("test convert Struct to Map 1", r)
	}
}
