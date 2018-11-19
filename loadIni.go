package helper

import (
	"os"
	"bufio"
	"io"
	"fmt"
	"strings"
	"regexp"
	"errors"
	"strconv"
)

const (
	TYPE_SECTION   = iota
	TYPE_PARAMETER
	TYPE_UNKNOWN
)

const (
	sectionReg   = `^\[(\w+)\]$`
	parameterReg = `(.+?)=(.*)`
)

type configStruct = map[string]map[string]string

var configCache map[string]configStruct

type Config struct {
	File    string
	Section string
}

type lineParseResult struct {
	Tpe int
	Key string
	Val string
}

func (s *lineParseResult) String() string {
	switch s.Tpe {
	case TYPE_SECTION:
		return "section: `" + s.Val + "`"
	case TYPE_PARAMETER:
		return "key:`" + s.Key + "`, val:`" + s.Val + "`"
	default:
		return "unknown: `" + s.Val + "`"
	}
}

func (c Config) Get(key string, defaultValue string) (string, error) {
	if strings.Index(key, ".") > 0 {
		sl := strings.Split(key, ".")
		c.File, c.Section, key = sl[0], sl[1], sl[2]
	}
	section, e := c.GetSection()
	if e != nil {
		return defaultValue, e
	}
	if _, ok := section[key]; !ok {
		return defaultValue, errors.New("Key: " + key + " is not exist")
	}
	return configCache[c.File][c.Section][key], nil
}

func (c Config) Int(key string, defaultValue int) (i int, e error) {
	s, e := c.Get(key, "")
	if e != nil {
		return
	}
	i, _ = strconv.Atoi(s)
	return
}

func (c Config) GetSection() (map[string]string, error) {
	if c.File == "" || c.Section == "" {
		return nil, errors.New("file or section is empty")
	}

	if configCache == nil {
		configCache = make(map[string]configStruct)
	}
	if _, ok := configCache[c.File]; !ok {
		filePath := c.File + ".ini"

		configCache[c.File] = make(configStruct)
		configCache[c.File], _ = ReadConf(filePath)
	}
	if _, ok := configCache[c.File][c.Section]; !ok {
		return nil, errors.New("Section: " + c.Section + " is not defined in " + c.File)
	}
	return configCache[c.File][c.Section], nil
}

func ReadConf(filePath string) (conf map[string]map[string]string, err error) {
	fh, err := os.Open(filePath)
	defer fh.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	var Conf = make(configStruct)
	currentSection := ""
	buff := bufio.NewReader(fh)
	eof := false
	for !eof {
		line, err := buff.ReadString('\n')
		if err == io.EOF {
			eof = true
		} else if err != nil {
			fmt.Println(err)
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0:1] == ";" {
			continue
		}

		parseResult := parseLine(line)
		if parseResult.Tpe == TYPE_UNKNOWN && parseResult.Val == "" {
			fmt.Println(parseResult)
			panic("Ini file: " + filePath + " parsed failed.")
		}

		switch parseResult.Tpe {
		case TYPE_SECTION:
			currentSection = parseResult.Val
		case TYPE_PARAMETER:
			if _, ok := Conf[currentSection]; !ok {
				Conf[currentSection] = make(map[string]string)
			}
			Conf[currentSection][parseResult.Key] = parseResult.Val
		case TYPE_UNKNOWN:
			continue
		}
		//fmt.Println(&parseResult)
	}

	return Conf, nil
}

func parseLine(line string) (lineParseResult) {
	if i := strings.Index(line, ";"); i >= 0 {
		line = line[0:i]
	}

	// section
	reg := regexp.MustCompile(sectionReg)
	rs := reg.FindStringSubmatch(line)
	if len(rs) > 0 {
		return lineParseResult{Tpe: TYPE_SECTION, Val: rs[1]}
	}

	// parameter
	reg = regexp.MustCompile(parameterReg)
	rs = reg.FindStringSubmatch(line)
	if len(rs) > 0 {
		name, value := strings.TrimSpace(rs[1]), strings.TrimSpace(rs[2])
		if name != "" {
			return lineParseResult{Tpe: TYPE_PARAMETER, Key: name, Val: value}
		}
	}

	return lineParseResult{Tpe: TYPE_UNKNOWN}
}
