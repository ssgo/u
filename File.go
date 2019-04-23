package u

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(command string, args ...string) ([]string, error) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
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
	fd.Close()
	return outs, nil
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
		os.MkdirAll(path, 0700)
	}
}

func Load(fileName string, to interface{}) error {
	CheckPath(fileName)

	fp, err := os.Open(fileName)
	if err == nil {
		defer fp.Close()
		decoder := json.NewDecoder(fp)
		err = decoder.Decode(&to)
	}
	return err
}

func Save(fileName string, data interface{}) error {
	CheckPath(fileName)

	fp, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err == nil {
		defer fp.Close()
		var buf []byte
		buf, err = json.MarshalIndent(data, "", "  ")
		if err == nil {
			fp.Write(buf)
		}
	}
	return err
}
