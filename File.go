package u

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type MemFile struct {
	Name       string
	absName    string
	ModTime    time.Time
	IsDir      bool
	Compressed bool
	Data       []byte
}

type MemFileB64 struct {
	Name       string
	ModTime    time.Time
	IsDir      bool
	DataB64    string
	Compressed bool
	Children   []MemFileB64
}

type FileInfo struct {
	Name     string
	FullName string
	IsDir    bool
	Size     int64
	ModTime  time.Time
}

var memFiles = make(map[string]*MemFile)
var memFilesByDir = make(map[string][]MemFile)
var memFilesLock sync.RWMutex

func (mf *MemFile) GetData() []byte {
	if mf.Compressed {
		return GunzipN(mf.Data)
	}
	return mf.Data
}

func GetAbsName(filename string) string {
	if !filepath.IsAbs(filename) {
		if absName, err := filepath.Abs(filename); err == nil {
			filename = absName
		}
	}
	return filename
}

func ReadFileFromMemory(name string) *MemFile {
	name = GetAbsName(name)
	memFilesLock.RLock()
	defer memFilesLock.RUnlock()
	return memFiles[name]
}

func ReadDirFromMemory(name string) []MemFile {
	name = GetAbsName(name)
	var dirFiles []MemFile
	memFilesLock.RLock()
	mfList := memFilesByDir[name]
	if mfList != nil {
		dirFiles = make([]MemFile, len(mfList))
		for i, mf := range mfList {
			dirFiles[i] = mf
		}
	}
	memFilesLock.RUnlock()
	return dirFiles
}

func AddFileToMemory(memFile MemFile) {
	memFile.Name = GetAbsName(memFile.Name)
	dirName := filepath.Dir(memFile.Name)
	memFilesLock.Lock()
	memFiles[memFile.Name] = &memFile
	if dirName != "" && dirName != "." {
		if memFilesByDir[dirName] == nil {
			memFilesByDir[dirName] = make([]MemFile, 0)
		}
		memFilesByDir[dirName] = append(memFilesByDir[dirName], memFile)
	}
	memFilesLock.Unlock()
}

func LoadFileToMemory(filename string) {
	loadFileToMemory(filename, false)
}

func LoadFileToMemoryWithCompress(filename string) {
	loadFileToMemory(filename, true)
}

func loadFileToMemory(filename string, compress bool) {
	if info, err := os.Stat(filename); err == nil {
		if info.IsDir() {
			AddFileToMemory(MemFile{
				Name:    filename,
				ModTime: info.ModTime(),
				IsDir:   true,
				Data:    nil,
			})
			if files, err := os.ReadDir(filename); err == nil {
				for _, file := range files {
					LoadFileToMemory(filepath.Join(filename, file.Name()))
				}
			}
		} else {
			if data, err := ReadFileBytes(filename); err == nil {
				compressed := false
				if compress {
					if data2, err := Gzip(data); err == nil {
						data = data2
						compressed = true
					}
				}
				AddFileToMemory(MemFile{
					Name:       filename,
					ModTime:    info.ModTime(),
					IsDir:      false,
					Data:       data,
					Compressed: compressed,
				})
			}
		}
	}
}

func Gzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	gzWriter := gzip.NewWriter(&buf)
	if _, err = gzWriter.Write(data); err == nil {
		err = gzWriter.Close()
		return buf.Bytes(), err
	}
	return data, err
}

func Gunzip(data []byte) ([]byte, error) {
	bufR := bytes.NewReader(data)
	var err error
	var gzReader *gzip.Reader
	if gzReader, err = gzip.NewReader(bufR); err == nil {
		defer gzReader.Close()
		var buf []byte
		if buf, err = io.ReadAll(gzReader); err == nil {
			return buf, nil
		}
	}
	return data, err
}

func GzipN(data []byte) []byte {
	buf, _ := Gzip(data)
	return buf
}

func GunzipN(data []byte) []byte {
	buf, _ := Gunzip(data)
	return buf
}

func LoadFileToB64(filename string) *MemFileB64 {
	if info, err := os.Stat(filename); err == nil {
		if info.IsDir() {
			out := MemFileB64{
				Name:     filename,
				ModTime:  info.ModTime(),
				IsDir:    true,
				Children: make([]MemFileB64, 0),
			}
			if files, err := os.ReadDir(filename); err == nil {
				for _, file := range files {
					if mfB64 := LoadFileToB64(filepath.Join(filename, file.Name())); mfB64 != nil {
						out.Children = append(out.Children, *mfB64)
					}
				}
			}
			return &out
		} else {
			if data, err := ReadFileBytes(filename); err == nil {
				compressed := false
				if buf, err := Gzip(data); err == nil {
					data = buf
					compressed = true
				}
				return &MemFileB64{
					Name:       filename,
					ModTime:    info.ModTime(),
					IsDir:      false,
					DataB64:    Base64(data),
					Compressed: compressed,
				}
			}
		}
	}
	return nil
}

func LoadFilesToMemoryFromB64(b64File *MemFileB64) {
	data := UnBase64(b64File.DataB64)
	if b64File.Compressed {
		data = GunzipN(data)
	}
	memFile := MemFile{
		Name:    b64File.Name,
		ModTime: b64File.ModTime,
		IsDir:   b64File.IsDir,
		Data:    data,
	}
	AddFileToMemory(memFile)
	if memFile.IsDir && len(b64File.Children) > 0 {
		for _, child := range b64File.Children {
			LoadFilesToMemoryFromB64(&child)
		}
	}
}

func LoadFilesToMemoryFromB64KeepGzip(b64File *MemFileB64) {
	data := UnBase64(b64File.DataB64)
	memFile := MemFile{
		Name:    b64File.Name,
		ModTime: b64File.ModTime,
		IsDir:   b64File.IsDir,
		Data:    data,
	}
	AddFileToMemory(memFile)
	if memFile.IsDir && len(b64File.Children) > 0 {
		for _, child := range b64File.Children {
			LoadFilesToMemoryFromB64(&child)
		}
	}
}

func LoadFilesToMemoryFromJson(jsonFiles string) {
	dirData := MemFileB64{}
	UnJson(jsonFiles, &dirData)
	LoadFilesToMemoryFromB64(&dirData)
}

func RunCommand(command string, args ...string) ([]string, error) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	outs := make([]string, 0)
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(io.MultiReader(stdout, stderr))
	for {
		lineBuf, _, err2 := reader.ReadLine()

		if err2 != nil || io.EOF == err2 {
			break
		}
		line := strings.TrimRight(string(lineBuf), "\r\n")
		outs = append(outs, line)
	}

	_ = cmd.Wait()
	return outs, nil
}

func ReadDir(filename string) ([]FileInfo, error) {
	out := make([]FileInfo, 0)
	if mfList := ReadDirFromMemory(filename); mfList != nil {
		for _, f := range mfList {
			out = append(out, FileInfo{
				Name:     filepath.Base(f.Name),
				FullName: f.Name,
				IsDir:    f.IsDir,
				Size:     int64(len(f.Data)),
				ModTime:  f.ModTime,
			})
		}
		return out, nil
	}

	files, err := os.ReadDir(filename)
	if err == nil {
		for _, f := range files {
			info, _ := f.Info()
			out = append(out, FileInfo{
				Name:     f.Name(),
				FullName: filepath.Join(filename, f.Name()),
				IsDir:    f.IsDir(),
				Size:     info.Size(),
				ModTime:  info.ModTime(),
			})
		}
	}
	return out, err
}

func ReadDirN(filename string) []FileInfo {
	files, _ := ReadDir(filename)
	return files
}

func ReadFileLines(filename string) ([]string, error) {
	outs := make([]string, 0)
	if mf := ReadFileFromMemory(filename); mf != nil {
		return strings.Split(string(mf.Data), "\n"), nil
	}
	fd, err := os.OpenFile(filename, os.O_RDONLY, 0400)
	if err != nil {
		return outs, err
	}

	inputReader := bufio.NewReader(fd)
	for {
		line, err := inputReader.ReadString('\n')
		line = strings.TrimRight(string(line), "\r\n")
		outs = append(outs, line)
		if err != nil {
			break
		}
	}
	_ = fd.Close()
	return outs, nil
}

func ReadFileLinesN(filename string) []string {
	lines, _ := ReadFileLines(filename)
	return lines
}

func ReadFile(filename string) (string, error) {
	if mf := ReadFileFromMemory(filename); mf != nil {
		return string(mf.Data), nil
	}
	buf, err := ReadFileBytes(filename)
	return string(buf), err
}

func ReadFileN(filename string) string {
	buf, _ := ReadFile(filename)
	return buf
}

func ReadFileBytes(filename string) ([]byte, error) {
	if mf := ReadFileFromMemory(filename); mf != nil {
		return mf.Data, nil
	}
	var maxLen uint
	if fi, _ := os.Stat(filename); fi != nil {
		maxLen = uint(fi.Size())
	} else {
		maxLen = 1024000
	}

	fd, err := os.OpenFile(filename, os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, maxLen)
	n, err := fd.Read(buf)
	_ = fd.Close()
	if err != nil {
		return nil, err
	}

	return buf[0:n], nil
}

func ReadFileBytesN(filename string) []byte {
	buf, _ := ReadFileBytes(filename)
	return buf
}

func WriteFile(filename string, content string) error {
	return WriteFileBytes(filename, []byte(content))
}

func WriteFileBytes(filename string, content []byte) error {
	absFilename := GetAbsName(filename)
	memFilesLock.RLock()
	mf := memFiles[absFilename]
	memFilesLock.RUnlock()
	if mf != nil {
		memFilesLock.Lock()
		mf.Data = content
		memFilesLock.Unlock()
	}

	CheckPath(filename)
	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	_, err = fd.Write(content)
	_ = fd.Close()
	if err != nil {
		return err
	}

	return nil
}

func FileExists(filename string) bool {
	if mf := ReadFileFromMemory(filename); mf != nil {
		return true
	}
	fi, err := os.Stat(filename)
	return err == nil && fi != nil
}

func CheckPath(filename string) {
	pos := strings.LastIndexByte(filename, os.PathSeparator)
	if pos < 0 {
		return
	}
	path := filename[0:pos]
	if _, err := os.Stat(path); err != nil {
		_ = os.MkdirAll(path, 0700)
	}
}

func FixPath(path string) string {
	const spe = string(os.PathSeparator)
	if !strings.HasSuffix(path, spe) {
		return path + spe
	}
	return path
}

var fileLocksLock = sync.Mutex{}
var fileLocks = map[string]*sync.Mutex{}

func LoadX(filename string, to interface{}) error {
	var in = map[string]interface{}{}
	if err := Load(filename, &in); err == nil {
		Convert(in, to)
		return nil
	} else {
		return err
	}
}

func Load(filename string, to interface{}) error {
	if strings.HasSuffix(filename, "yml") || strings.HasSuffix(filename, "yaml") {
		return load(filename, true, to)
	} else {
		return load(filename, false, to)
	}
}

func LoadYaml(filename string, to interface{}) error {
	return load(filename, true, to)
}

func LoadJson(filename string, to interface{}) error {
	return load(filename, false, to)
}

func load(filename string, isYaml bool, to interface{}) error {
	if mf := ReadFileFromMemory(filename); mf != nil {
		if isYaml {
			return yaml.Unmarshal(mf.Data, to)
		} else {
			return json.Unmarshal(mf.Data, to)
		}
	}

	fileLocksLock.Lock()
	if fileLocks[filename] == nil {
		fileLocks[filename] = new(sync.Mutex)
	}
	lock := fileLocks[filename]
	fileLocksLock.Unlock()

	lock.Lock()
	defer lock.Unlock()

	fp, err := os.Open(filename)
	if err == nil {
		if isYaml {
			decoder := yaml.NewDecoder(fp)
			err = decoder.Decode(to)
		} else {
			decoder := json.NewDecoder(fp)
			err = decoder.Decode(to)
		}
		_ = fp.Close()
	}
	return err
}

func CopyToFile(from io.Reader, to string) error {
	if fp, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600); err == nil {
		defer fp.Close()
		io.Copy(fp, from)
		return nil
	} else {
		return err
	}
}

func CopyFile(from, to string) error {
	if writer, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600); err == nil {
		defer writer.Close()
		if reader, err := os.OpenFile(from, os.O_RDONLY, 0600); err == nil {
			defer reader.Close()
			_, err = io.Copy(writer, reader)
			return err
		} else {
			return err
		}
	} else {
		return err
	}
}

func Save(filename string, data interface{}) error {
	if strings.HasSuffix(filename, "yml") || strings.HasSuffix(filename, "yaml") {
		return save(filename, true, data, true)
	} else {
		return save(filename, false, data, false)
	}
}

func SaveYaml(filename string, data interface{}) error {
	return save(filename, true, data, true)
}

func SaveJson(filename string, data interface{}) error {
	return save(filename, false, data, false)
}

func SaveJsonP(filename string, data interface{}) error {
	return save(filename, false, data, true)
}

func save(filename string, isYaml bool, data interface{}, indent bool) error {
	var buf []byte
	var err error
	if isYaml {
		buf, err = yaml.Marshal(data)
	} else {
		buffer := bytes.Buffer{}
		enc := json.NewEncoder(&buffer)
		enc.SetEscapeHTML(false)
		if indent {
			enc.SetIndent("", "  ")
		}
		err := enc.Encode(data)

		//buf, err = json.Marshal(data)
		if err == nil {
			buf = buffer.Bytes()
			FixUpperCase(buf, nil)
			//if indent {
			//	buf2 := bytes.Buffer{}
			//	err2 := json.Indent(&buf2, buf, "", "  ")
			//	if err2 == nil {
			//		buf = buf2.Bytes()
			//	}
			//}
		}
	}
	if err != nil {
		return err
	}

	absFilename := GetAbsName(filename)
	memFilesLock.RLock()
	mf := memFiles[absFilename]
	memFilesLock.RUnlock()
	if mf != nil {
		memFilesLock.Lock()
		mf.Data = buf
		memFilesLock.Unlock()
	}

	CheckPath(filename)
	fileLocksLock.Lock()
	if fileLocks[filename] == nil {
		fileLocks[filename] = new(sync.Mutex)
	}
	lock := fileLocks[filename]
	fileLocksLock.Unlock()

	lock.Lock()
	defer lock.Unlock()

	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err == nil {
		_, err = fp.Write(buf)
		_ = fp.Close()
	}
	return err
}
