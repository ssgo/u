package u

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type MemFile struct {
	Name       string
	absName    string
	ModTime    time.Time
	IsDir      bool
	Compressed bool
	Size       int64
	Data       []byte
	SafeData   *SafeBuf
}

type MemFileB64 struct {
	Name       string
	ModTime    time.Time
	IsDir      bool
	DataB64    []byte
	Compressed bool
	Size       int64
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

func (mf *MemFile) GetSafeData() *SecretPlaintext {
	return mf.SafeData.Open()
}

func GetAbsFilename(filename string) string {
	if !filepath.IsAbs(filename) {
		if absName, err := filepath.Abs(filename); err == nil {
			filename = absName
		}
	}
	return filename
}

func ReadFileFromMemory(name string) *MemFile {
	name = GetAbsFilename(name)
	memFilesLock.RLock()
	defer memFilesLock.RUnlock()
	return memFiles[name]
}

func ReadDirFromMemory(name string) []MemFile {
	name = GetAbsFilename(name)
	var dirFiles []MemFile
	memFilesLock.RLock()
	mfList := memFilesByDir[name]
	if mfList != nil {
		dirFiles = make([]MemFile, len(mfList))
		copy(dirFiles, mfList)
	}
	memFilesLock.RUnlock()
	return dirFiles
}

func AddFileToMemory(memFile MemFile) {
	memFile.Name = GetAbsFilename(memFile.Name)
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
	loadFileToMemory(filename, false, false)
}

func SafeLoadFileToMemory(filename string) {
	loadFileToMemory(filename, false, true)
}

func LoadFileToMemoryWithCompress(filename string) {
	loadFileToMemory(filename, true, false)
}

func SafeLoadFileToMemoryWithCompress(filename string) {
	loadFileToMemory(filename, true, true)
}

func loadFileToMemory(filename string, compress bool, isSafe bool) {
	if info, err := os.Stat(filename); err == nil {
		if info.IsDir() {
			AddFileToMemory(MemFile{
				Name:    filename,
				ModTime: info.ModTime(),
				IsDir:   true,
				Size:    info.Size(),
			})
			if files, err := os.ReadDir(filename); err == nil {
				for _, file := range files {
					loadFileToMemory(filepath.Join(filename, file.Name()), compress, isSafe)
				}
			}
		} else {
			if data, err := ReadFileBytes(filename); err == nil {
				compressed := false
				var dataBuf *SafeBuf
				if compress {
					if data2, err := Gzip(data); err == nil {
						if isSafe {
							ZeroMemory(data) // 安全模式下，清空原始数据
						}
						data = data2
						compressed = true
					}
				}
				if isSafe {
					dataBuf = NewSafeBuf(data)
					ZeroMemory(data) // 安全模式下，清空原始数据（或压缩后的明文数据）
					data = nil
				}
				AddFileToMemory(MemFile{
					Name:       filename,
					ModTime:    info.ModTime(),
					IsDir:      false,
					Size:       info.Size(),
					Data:       data,
					SafeData:   dataBuf,
					Compressed: compressed,
				})
			}
		}
	}
}

// func Gzip(data []byte) ([]byte, error) {
// 	var buf bytes.Buffer
// 	var err error
// 	gzWriter := gzip.NewWriter(&buf)
// 	if _, err = gzWriter.Write(data); err == nil {
// 		err = gzWriter.Close()
// 		return buf.Bytes(), err
// 	}
// 	return data, err
// }

// func Gunzip(data []byte) ([]byte, error) {
// 	bufR := bytes.NewReader(data)
// 	var err error
// 	var gzReader *gzip.Reader
// 	if gzReader, err = gzip.NewReader(bufR); err == nil {
// 		defer gzReader.Close()
// 		var buf []byte
// 		if buf, err = io.ReadAll(gzReader); err == nil {
// 			return buf, nil
// 		}
// 	}
// 	return data, err
// }

// func GzipN(data []byte) []byte {
// 	buf, _ := Gzip(data)
// 	return buf
// }

// func GunzipN(data []byte) []byte {
// 	buf, _ := Gunzip(data)
// 	return buf
// }

func Compress(data []byte, cType string) ([]byte, error) {
	var buf bytes.Buffer
	var w io.WriteCloser

	lcType := strings.ToLower(cType)
	switch lcType {
	case "gzip", "gz":
		w = gzip.NewWriter(&buf)
	default:
		w = zlib.NewWriter(&buf)
	}

	if _, err := w.Write(data); err != nil {
		return data, err
	}
	w.Close()
	return buf.Bytes(), nil
}

func Decompress(data []byte, cType string) ([]byte, error) {
	bufR := bytes.NewReader(data)
	var r io.ReadCloser
	var err error

	lcType := strings.ToLower(cType)
	switch lcType {
	case "gzip", "gz":
		r, err = gzip.NewReader(bufR)
	default:
		r, err = zlib.NewReader(bufR)
	}

	if err != nil {
		return data, err
	}
	defer r.Close()
	if out, err := io.ReadAll(r); err == nil {
		return out, nil
	}
	return data, err
}

// 你要求的原有接口保持不变
func Gzip(data []byte) ([]byte, error)   { return Compress(data, "gzip") }
func Gunzip(data []byte) ([]byte, error) { return Decompress(data, "gzip") }
func GzipN(data []byte) []byte           { b, _ := Gzip(data); return b }
func GunzipN(data []byte) []byte         { b, _ := Gunzip(data); return b }
func Zip(data []byte) ([]byte, error)    { return Compress(data, "zlib") }
func Unzip(data []byte) ([]byte, error)  { return Decompress(data, "zlib") }
func ZipN(data []byte) []byte            { b, _ := Zip(data); return b }
func UnzipN(data []byte) []byte          { b, _ := Unzip(data); return b }

// Extract 自动识别格式并解压 (支持 .zip, .tar.gz, .tgz, .tar, .gz, .bz2)
func Extract(srcFile, destDir string, stripRoot bool) error {
	f, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer f.Close()

	lowerSrc := strings.ToLower(srcFile)

	switch {
	case strings.HasSuffix(lowerSrc, ".zip"):
		return extractZip(srcFile, destDir, stripRoot)

	case strings.HasSuffix(lowerSrc, ".tar.gz") || strings.HasSuffix(lowerSrc, ".tgz"):
		gzr, _ := gzip.NewReader(f)
		defer gzr.Close()
		return extractTar(gzr, destDir, stripRoot)

	case strings.HasSuffix(lowerSrc, ".tar"):
		return extractTar(f, destDir, stripRoot)

	case strings.HasSuffix(lowerSrc, ".gz"):
		// 纯 gzip 文件，解压为单个文件
		gzr, _ := gzip.NewReader(f)
		defer gzr.Close()
		return extractSingleFile(gzr, destDir, strings.TrimSuffix(filepath.Base(srcFile), ".gz"))

	case strings.HasSuffix(lowerSrc, ".bz2"):
		// bzip2 仅支持解压
		bzr := bzip2.NewReader(f)
		return extractSingleFile(bzr, destDir, strings.TrimSuffix(filepath.Base(srcFile), ".bz2"))

	default:
		return extractZip(srcFile, destDir, stripRoot)
	}
}

// extractTar 处理所有的 Tar 变体 (tar, tar.gz)
func extractTar(r io.Reader, dest string, strip bool) error {
	tr := tar.NewReader(r)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := writeEntry(dest, h.Name, h.FileInfo().Mode(), h.Typeflag == tar.TypeDir, h.Linkname, tr, strip); err != nil {
			return err
		}
	}
}

// extractZip 处理 Zip
func extractZip(src, dest string, strip bool) error {
	rz, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer rz.Close()
	for _, f := range rz.File {
		rc, _ := f.Open()
		var linkTarget string
		if f.Mode()&os.ModeSymlink != 0 {
			b, _ := io.ReadAll(rc)
			linkTarget = string(b)
		}
		err := writeEntry(dest, f.Name, f.Mode(), f.FileInfo().IsDir(), linkTarget, rc, strip)
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// extractSingleFile 用于处理 .gz 或 .bz2 这种单文件压缩
func extractSingleFile(r io.Reader, destDir, fileName string) error {
	os.MkdirAll(destDir, 0755)
	target := filepath.Join(destDir, fileName)
	out, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, r)
	return err
}

// writeEntry 核心磁盘写入逻辑
func writeEntry(destDir, name string, mode os.FileMode, isDir bool, linkPath string, r io.Reader, strip bool) error {
	if strip {
		parts := strings.SplitN(name, "/", 2)
		if len(parts) < 2 || parts[1] == "" {
			return nil
		}
		name = parts[1]
	}
	target := filepath.Join(destDir, name)

	if isDir {
		return os.MkdirAll(target, 0755)
	}
	os.MkdirAll(filepath.Dir(target), 0755)

	if mode&os.ModeSymlink != 0 {
		os.Remove(target)
		return os.Symlink(linkPath, target)
	}

	out, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, r)
	return err
}

// Archive 将目录或文件打包成 zip 或 tar.gz
func Archive(srcPath, destFile string) error {
	f, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer f.Close()

	lower := strings.ToLower(destFile)
	if strings.HasSuffix(lower, ".zip") {
		zw := zip.NewWriter(f)
		defer zw.Close()
		return walkAndAdd(srcPath, func(relPath string, info os.FileInfo, fileReader io.Reader) error {
			header, _ := zip.FileInfoHeader(info)
			header.Name = relPath
			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}
			w, err := zw.CreateHeader(header)
			if err != nil || info.IsDir() {
				return err
			}
			_, err = io.Copy(w, fileReader)
			return err
		})
	}

	// 默认处理为 .tar.gz
	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	return walkAndAdd(srcPath, func(relPath string, info os.FileInfo, fileReader io.Reader) error {
		header, _ := tar.FileInfoHeader(info, "")
		header.Name = relPath
		if info.Mode()&os.ModeSymlink != 0 {
			link, _ := os.Readlink(filepath.Join(srcPath, "..", relPath)) // 简化逻辑
			header.Linkname = link
		}
		if err := tw.WriteHeader(header); err != nil || info.IsDir() || header.Typeflag == tar.TypeSymlink {
			return err
		}
		_, err = io.Copy(tw, fileReader)
		return err
	})
}

// walkAndAdd 内部通用的文件树遍历辅助函数
func walkAndAdd(srcPath string, addFn func(string, os.FileInfo, io.Reader) error) error {
	baseDir := filepath.Dir(srcPath)
	return filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(baseDir, path)
		if relPath == "." {
			return nil
		}

		var r io.Reader
		if !info.IsDir() && info.Mode()&os.ModeSymlink == 0 {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			r = f
		}
		return addFn(filepath.ToSlash(relPath), info, r)
	})
}

func LoadFileToB64(filename string) *MemFileB64 {
	if info, err := os.Stat(filename); err == nil {
		if info.IsDir() {
			out := MemFileB64{
				Name:     filename,
				ModTime:  info.ModTime(),
				IsDir:    true,
				Size:     info.Size(),
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
					Size:       info.Size(),
					DataB64:    Base64(data),
					Compressed: compressed,
				}
			}
		}
	}
	return nil
}

func LoadFilesToMemoryFromB64(b64File *MemFileB64) {
	if data, err := UnBase64(b64File.DataB64); err == nil {
		if b64File.Compressed {
			data = GunzipN(data)
		}
		memFile := MemFile{
			Name:    b64File.Name,
			ModTime: b64File.ModTime,
			IsDir:   b64File.IsDir,
			Size:    b64File.Size,
			Data:    data,
		}
		AddFileToMemory(memFile)
		if memFile.IsDir && len(b64File.Children) > 0 {
			for _, child := range b64File.Children {
				LoadFilesToMemoryFromB64(&child)
			}
		}
	}
}

func LoadFilesToMemoryFromB64KeepGzip(b64File *MemFileB64) {
	if data, err := UnBase64(b64File.DataB64); err == nil {
		memFile := MemFile{
			Name:    b64File.Name,
			ModTime: b64File.ModTime,
			IsDir:   b64File.IsDir,
			Size:    b64File.Size,
			Data:    data,
		}
		AddFileToMemory(memFile)
		if memFile.IsDir && len(b64File.Children) > 0 {
			for _, child := range b64File.Children {
				LoadFilesToMemoryFromB64(&child)
			}
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
			if info, _ := f.Info(); info != nil {
				out = append(out, FileInfo{
					Name:     f.Name(),
					FullName: filepath.Join(filename, f.Name()),
					IsDir:    info.IsDir(),
					Size:     info.Size(),
					ModTime:  info.ModTime(),
				})
			}
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
	absFilename := GetAbsFilename(filename)
	memFilesLock.RLock()
	mf := memFiles[absFilename]
	memFilesLock.RUnlock()
	if mf != nil {
		memFilesLock.Lock()
		mf.Data = content
		memFilesLock.Unlock()
	}

	CheckPath(filename)
	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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

func GetFileInfo(filename string) *FileInfo {
	if mf := ReadFileFromMemory(filename); mf != nil {
		return &FileInfo{
			Name:     mf.Name,
			FullName: mf.absName,
			IsDir:    mf.IsDir,
			Size:     mf.Size,
			ModTime:  mf.ModTime,
		}
	}
	if fi, err := os.Stat(filename); err == nil {
		fullName := filename
		if !filepath.IsAbs(filename) {
			fullName, _ = filepath.Abs(filename)
		}
		return &FileInfo{
			Name:     filename,
			FullName: fullName,
			IsDir:    fi.IsDir(),
			Size:     fi.Size(),
			ModTime:  fi.ModTime(),
		}
	}
	return nil
}

func CheckPath(filename string) {
	pos := strings.LastIndexByte(filename, os.PathSeparator)
	if pos < 0 {
		return
	}
	path := filename[0:pos]
	if _, err := os.Stat(path); err != nil {
		_ = os.MkdirAll(path, 0755)
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

func LoadX(filename string, to any) error {
	var in = map[string]any{}
	if err := Load(filename, &in); err == nil {
		Convert(in, to)
		return nil
	} else {
		return err
	}
}

func Load(filename string, to any) error {
	if strings.HasSuffix(filename, "yml") || strings.HasSuffix(filename, "yaml") {
		return load(filename, true, to)
	} else {
		return load(filename, false, to)
	}
}

func LoadYaml(filename string, to any) error {
	return load(filename, true, to)
}

func LoadJson(filename string, to any) error {
	return load(filename, false, to)
}

func load(filename string, isYaml bool, to any) error {
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
	if fp, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err == nil {
		defer fp.Close()
		io.Copy(fp, from)
		return nil
	} else {
		return err
	}
}

func CopyFile(from, to string) error {
	fromStat, err := os.Stat(from)
	if err != nil {
		return err
	}
	if fromStat.IsDir() {
		// copy dir
		for _, f := range ReadDirN(from) {
			err := CopyFile(filepath.Join(from, f.Name), filepath.Join(to, f.Name))
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		// copy file
		toStat, err := os.Stat(to)
		if err == nil && toStat.IsDir() {
			to = filepath.Join(to, filepath.Base(from))
		}
		CheckPath(to)
		if writer, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err == nil {
			defer writer.Close()
			if reader, err := os.OpenFile(from, os.O_RDONLY, 0644); err == nil {
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
}

func Save(filename string, data any) error {
	if strings.HasSuffix(filename, "yml") || strings.HasSuffix(filename, "yaml") {
		return save(filename, true, data, true)
	} else {
		return save(filename, false, data, false)
	}
}

func SaveYaml(filename string, data any) error {
	return save(filename, true, data, true)
}

func SaveJson(filename string, data any) error {
	return save(filename, false, data, false)
}

func SaveJsonP(filename string, data any) error {
	return save(filename, false, data, true)
}

func save(filename string, isYaml bool, data any, indent bool) error {
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

	absFilename := GetAbsFilename(filename)
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

	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err == nil {
		_, err = fp.Write(buf)
		_ = fp.Close()
	}
	return err
}
