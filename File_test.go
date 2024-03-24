package u_test

import (
	"bytes"
	"github.com/ssgo/u"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGzip(t *testing.T) {
	str := "hello 123!"
	buf := u.GzipN([]byte(str))
	str2 := string(u.GunzipN(buf))
	if str2 != str {
		t.Fatal("gzip error", str2, buf)
	}
}

func TestFileAndMemFile(t *testing.T) {
	testFile1 := filepath.Join("testDir", "testDir2", "testFile1.txt")
	testFile2 := filepath.Join("testDir", "testFile2.hex")
	u.WriteFile(testFile1, "abc123")
	u.WriteFileBytes(testFile2, []byte{0, 1, 2, 3, 255})
	str1 := u.ReadFileN(testFile1)
	buf2 := u.ReadFileBytesN(testFile2)
	if str1 != "abc123" || bytes.Compare(buf2, []byte{0, 1, 2, 3, 255}) != 0 {
		t.Fatal("read file error", str1, buf2)
	}

	type TestData struct {
		AAA string
		BBB int
	}
	testFile3 := filepath.Join("testDirX", "testDir2", "testFile1.json")
	testFile4 := filepath.Join("testDirX", "testFile2.yml")
	testData := TestData{"111", 222}
	u.SaveJson(testFile3, testData)
	u.SaveYaml(testFile4, testData)
	testData3 := TestData{}
	testData4 := TestData{}
	u.LoadX(testFile3, &testData3)
	u.LoadYaml(testFile4, &testData4)
	if testData3.AAA != "111" || testData3.BBB != 222 || testData4.AAA != "111" || testData4.BBB != 222 {
		t.Fatal("read json&yaml file error", testData3, testData4)
	}

	files := u.ReadDirN("testDir")
	if len(files) != 2 || !files[0].IsDir || files[1].Name != "testFile2.hex" || files[1].FullName != "testDir/testFile2.hex" {
		t.Fatal("read dir error", files)
	}

	u.LoadFileToMemory("testDir")
	u.LoadFileToMemoryWithCompress(testFile1)
	memFiles := u.LoadFileToB64("testDirX")
	//fmt.Println(u.JsonP(memFiles), 111)
	u.LoadFilesToMemoryFromB64(memFiles)
	os.RemoveAll("testDir")
	os.RemoveAll("testDirX")

	buf1 := u.ReadFileBytesN(testFile1)
	str1 = string(u.GunzipN(buf1))
	buf2 = u.ReadFileBytesN(testFile2)
	if str1 != "abc123" || bytes.Compare(buf2, []byte{0, 1, 2, 3, 255}) != 0 {
		t.Fatal("read mem file error", str1, buf2)
	}

	testData3 = TestData{}
	testData4 = TestData{}
	u.LoadX(testFile3, &testData3)
	u.LoadYaml(testFile4, &testData4)
	if testData3.AAA != "111" || testData3.BBB != 222 || testData4.AAA != "111" || testData4.BBB != 222 {
		t.Fatal("read mem json&yaml file error", testData3, testData4)
	}

	files = u.ReadDirN("testDir")
	if len(files) != 2 || !files[0].IsDir || filepath.Base(files[1].Name) != "testFile2.hex" || !strings.HasSuffix(files[1].FullName, "testDir/testFile2.hex") {
		t.Fatal("read dir error", files)
	}
}
