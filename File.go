package u

import (
	"bufio"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

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
	cmd.Start()
	reader := bufio.NewReader(io.MultiReader(stdout, stderr))
	for {
		lineBuf, _, err2 := reader.ReadLine()

		if err2 != nil || io.EOF == err2 {
			break
		}
		line := strings.TrimRight(string(lineBuf), "\r\n")
		outs = append(outs, line)
	}

	cmd.Wait()
	return outs, nil
}

func ReadFile(fileName string) ([]string, error) {
	outs := make([]string, 0)
	fd, err := os.OpenFile(fileName, os.O_RDONLY, 0400)
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

func ReadFullFile(fileName string, maxLen uint) (string, error) {
	fd, err := os.OpenFile(fileName, os.O_RDONLY, 0400)
	if err != nil {
		return "", err
	}

	buf := make([]byte, maxLen)
	n, err := fd.Read(buf)
	_ = fd.Close()
	if err != nil {
		return "", err
	}

	return string(buf[0:n]), nil
}

func FileExists(fileName string) bool {
	fi, err := os.Stat(fileName)
	return err == nil && fi != nil
}

func CheckPath(fileName string) {
	pos := strings.LastIndexByte(fileName, '/')
	if pos < 0 {
		return
	}
	path := fileName[0:pos]
	if _, err := os.Stat(path); err != nil {
		_ = os.MkdirAll(path, 0700)
	}
}

var fileLocks = map[string]*sync.Mutex{}

func Load(fileName string, to interface{}) error {
	if strings.HasSuffix(fileName, "yml") || strings.HasSuffix(fileName, "yaml") {
		return load(fileName, true, to)
	} else {
		return load(fileName, false, to)
	}
}

func LoadYaml(fileName string, to interface{}) error {
	return load(fileName, true, to)
}

func LoadJson(fileName string, to interface{}) error {
	return load(fileName, false, to)
}

func load(fileName string, isYaml bool, to interface{}) error {
	if fileLocks[fileName] == nil {
		fileLocks[fileName] = new(sync.Mutex)
	}
	fileLocks[fileName].Lock()
	defer fileLocks[fileName].Unlock()

	fp, err := os.Open(fileName)
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

func Save(fileName string, data interface{}) error {
	if strings.HasSuffix(fileName, "yml") || strings.HasSuffix(fileName, "yaml") {
		return save(fileName, true, data)
	} else {
		return save(fileName, false, data)
	}
}

func SaveYaml(fileName string, data interface{}) error {
	return save(fileName, true, data)
}

func SaveJson(fileName string, data interface{}) error {
	return save(fileName, false, data)
}

func save(fileName string, isYaml bool, data interface{}) error {
	CheckPath(fileName)

	if fileLocks[fileName] == nil {
		fileLocks[fileName] = new(sync.Mutex)
	}
	fileLocks[fileName].Lock()
	defer fileLocks[fileName].Unlock()

	fp, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err == nil {
		var buf []byte
		if isYaml {
			buf, err = yaml.Marshal(data)
		} else {
			buf, err = json.MarshalIndent(data, "", "  ")
		}
		if err == nil {
			_, err = fp.Write(buf)
			_ = fp.Close()
		}
	}
	return err
}
