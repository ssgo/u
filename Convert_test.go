package u_test

import (
	"fmt"
	"github.com/ssgo/u"
	"reflect"
	"testing"
)

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
		sex  string
	}

	from := map[string]interface{}{
		"aaa":   "111",
		"bbb":   222,
		"bbb2":  222,
		"bbb3":  222,
		"ccc":   333.333,
		"ccc-2": []string{"1", "2", "3.0"},
		"ccc3":  []int{1, 2, 3},
		"ccc4":  []interface{}{"1", 2, 3.0},
		"ccc5": []map[string]interface{}{
			{
				"name": 13,
			},
			{
				"name": "ALin",
				"age":  22,
			},
		},
		"sex": "male",
	}

	type CCC struct {
		Ccc    string
		Ccc2__ []int
		Ccc3   []interface{}
		Ccc4   []*float32
		Ccc5   []*User
		sex    string
	}
	type BBB struct {
		Bbb  *string
		BBB2 []byte
		Bbb3 ****string
		CCC
		sex string
	}
	to := struct {
		Aaa int
		BBB
		sex string
	}{}

	u.Convert(from, &to)

	if to.Aaa != 111 || string(to.BBB2) != "222" || to.Ccc != "333.333" {
		t.Error("test convert Map to Struct 1", to)
	}

	if (*to.Bbb) != "222" || (****to.Bbb3) != "222" {
		t.Error("test convert Map to Struct 2", to)
	}

	if len(to.Ccc2__) < 3 || to.Ccc2__[0] != 1 || to.Ccc2__[1] != 2 || to.Ccc2__[2] != 3 {
		t.Error("test convert Slice to Slice 1", to.Ccc2__)
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

	u.Convert(from, &to)
	if len(to) < 2 || to[0]["name"] != "Andy" || to[1]["name"] != "Kitty" || to[1]["age"] != 22 {
		t.Error("test convert Struct to Map 2", u.JsonP(to), len(to) < 2, to[0]["name"] != "Andy", to[1]["name"] != "Kitty", to[1]["age"] != 22)
	}
}

func TestConvertSliceToSlice(t *testing.T) {
	from := []int{1, 2, 3}
	to := make([]string, 0)

	xto := &to
	u.Convert(from, &xto)
	if len(to) < 3 || to[0] != "1" || to[1] != "2" || to[2] != "3" {
		t.Error("test convert Slice to Slice 1", to)
	}

	from2 := 9
	to2 := make([]string, 0)

	u.Convert(from2, &to2)
	if len(to2) < 1 || to2[0] != "9" {
		t.Error("test convert Slice to Slice 2", to2)
	}
}

func TestConvertStructToStruct(t *testing.T) {
	type User1 struct {
		My_Name string
		Age     int
		F1      func(a int, b string, c bool) (d int, e string, f bool)
		F2      func(a ...interface{}) interface{}
	}

	type User2 struct {
		MYNAME_ string
		Level   int
		Class   int
		F1      func(a int, b int, c bool) (d int, e string, f bool)
		F2      func(a int, b int, c bool) (d int, e string, f bool)
	}

	from := User1{My_Name: "Tom", Age: 23, F1: func(a int, b string, c bool) (d int, e string, f bool) {
		return a + 1, b + "1", !c
	}, F2: func(a ...interface{}) interface{} {
		return []interface{}{u.Int(a[0]) + 1, u.String(a[1]) + "1", !u.Bool(a[2])}
	}}
	to := User2{MYNAME_: "Jeff", Level: -1}

	u.Convert(from, &to)
	if to.MYNAME_ != "Tom" {
		t.Error("test convert Struct to Struct", to)
	}

	a, b, c := to.F1(1, 1, false)
	if a != 2 || b != "11" || c != true {
		t.Error("test convert Func to Func", a, b, c)
	}

	a, b, c = to.F1(1, 1, false)
	if a != 2 || b != "11" || c != true {
		t.Error("test convert Func to Func with variadic", a, b, c)
	}
}

type TestStructA struct {
	field1 string
	Field2 string
}

func (item *TestStructA) method1() {
}

func (item *TestStructA) Method2() {
}

type TestStructB struct {
	TestStructA
	field3 string
	Field4 string
}

func (item *TestStructB) method3() {
}

func (item *TestStructB) Method4() {
}

func TestFlatStruct(t *testing.T) {
	data := &TestStructB{
		TestStructA: TestStructA{
			field1: "111",
			Field2: "222",
		},
		field3: "333",
		Field4: "444",
	}

	info := u.FlatStructWithUnexported(data)
	fmt.Println(u.JsonP(info.Values))
	fmt.Println(u.JsonP(info.MethodValues))
}

//func TestToInterfaceArray(t *testing.T) {
//	in := []int{1,2,3}
//	out := u.ToInterfaceArray(in)
//	fmt.Println(u.JsonP(out))
//}

//func TestConvertS(t *testing.T) {
//	s := `{
//  "Apps": null,
//  "Binds": {
//    "xxx": [
//      "aaa",
//      "bbb"
//    ]
//  },
//  "Desc": "",
//  "Name": "",
//  "Vars": null,
//  "name": "c1"
//}`
//	args := map[string]interface{}{}
//	_ = json.Unmarshal([]byte(s), &args)
//
//	in := struct {
//		Name  string
//		Desc  string
//		Vars  map[string]*string
//		Binds map[string][]string
//	}{}
//	u.Convert(args, &in)
//
//	fmt.Println(u.JsonP(in))
//}

//func TestConvertIntArray(t *testing.T) {
//	s := `{
//      "department": [2,3,7],
//      "name": "AAA",
//      "userid": "aaa"
//    }`
//
//	args := map[string]interface{}{}
//	_ = json.Unmarshal([]byte(s), &args)
//	fmt.Println(u.JsonP(args))
//
//	in := struct {
//		Department  []string
//		Name  string
//	}{}
//	u.Convert(args, &in)
//
//	fmt.Println(u.JsonP(in))
//}
