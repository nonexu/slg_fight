package gd_config

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	//"fmt"
)

var  path = "../../gd_config/"

func init() {
	LoadTables()
	loadData()
}

var (
	_tables map[string]*Table
)

type Record struct {
	Fields map[string]string
}

type Table struct {
	Records []*Record
}

type Query struct {
	Key   string
	Value string
}

func Sto8(v string) int8 {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return int8(n)
}

func Sto16(v string) int16 {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return int16(n)
}

func Sto32(v string) int32 {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return int32(n)
}

func Sto64(v string) int64 {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return int64(n)
}

func Stoi(v string) int {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return n
}

func Stou(v string) uint {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return uint(n)
}

func Stou64(v string) uint64 {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return uint64(n)
}

func Stobool(v string) bool {
	if v == "0" {
		return false
	} else {
		return true
	}
}

func Stobool2(v string) bool {
	if v == "1" {
		return true
	} else {
		return false
	}
}

func Stof32(v string) float32 {
	if val, err := strconv.ParseFloat(v, 32); err == nil {
		return float32(val)
	}
	return 0
}

func Stof64(v string) float64 {
	if val, err := strconv.ParseFloat(v, 64); err == nil {
		return val
	}
	return 0
}

func GetLines(table string, querys []*Query) []*Record {
	return getLines(table, querys)
}

func GetLine(table string, querys []*Query) *Record {
	ret := getLines(table, querys)
	if len(ret) > 0 {
		return ret[0]
	}
	return nil
}

func GetAll(table string) []*Record {
	return getAll(table)
}

// 加载目录下的所有csv文件
func LoadTables() {
	_tables = make(map[string]*Table)

	path := filepath.Join(path, "*.csv")

	files, err := filepath.Glob(path)
	if err != nil {
		//loG.println("get file list err: " + err.Error())
		panic(err)
	}

	total, ok := 0, 0
	for _, fn := range files {

		total += 1
		f, err := os.Open(fn)
		if err != nil {
			//loG.println("Open file error", err.Error(), fn)
			continue
		}
		defer f.Close()

		// 通过gd_config/xxx.csv得到xxx
		name := strings.TrimSuffix(filepath.Base(fn), filepath.Ext(fn))
		initTable(name, f)
		//fmt.Println(name)
		ok += 1
	}
}

func ReloadByName(table string) bool {

	fn := filepath.Join(path, table+".csv")
	f, err := os.Open(fn)
	if err != nil {
		//loG.println("Open file error", err.Error(), fn)
		return false
	}
	defer f.Close()
	initTable(table, f)

	return true
}

func prevFormat(s string) string {
	if len(s) < 4 {
		return s
	}

	v := s[:4]
	if v == "INT_" || v == "STR_" || v == "int_" || v == "str_" {
		return s[4:]
	}
	return s
}

func initTable(table string, f *os.File) {
	reader := csv.NewReader(f)
	title, err := reader.Read()
	if err != nil {
		//loG.println("Error:" + err.Error())
		panic(err)
	}

	// 去除 titile 格式前缀(INT_ STR_ 等)
	for idx, val := range title {
		title[idx] = prevFormat(val)
	}

	t := Table{Records: make([]*Record, 0)}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if nil != err {
			//fileInfo, _ := f.Stat()
			//log.Printf("file %v err: %v\n", fileInfo.Name(), err)
			return
		}

		rec := Record{Fields: make(map[string]string)}
		for idx, val := range line {
			rec.Fields[title[idx]] = val
		}
		t.Records = append(t.Records, &rec)
	}

	_tables[table] = &t
}

func getOne(tableName string, queryField string, val string) *Record {
	table, ok := _tables[tableName]
	if !ok {
		return nil
	}

	for _, rec := range table.Records {
		v, ok := rec.Fields[queryField]
		if !ok {
			return nil
		}
		if v == val {
			return rec
		}
	}

	return nil
}

func getLines(tableName string, querys []*Query) []*Record {
	table, ok := _tables[tableName]
	if !ok {
		return nil
	}

	var match bool
	lines := make([]*Record, 0)
	for _, rec := range table.Records {
		match = true
		for _, query := range querys {
			v, ok := rec.Fields[query.Key]
			if !ok {
				match = false
				break
			}
			if v != query.Value {
				match = false
				break
			}
		}
		if match {
			lines = append(lines, rec)
		}
	}

	return lines
}

func getAll(tableName string) []*Record {
	table, ok := _tables[tableName]
	if !ok {
		return nil
	}

	lines := make([]*Record, len(table.Records))
	idx := 0
	for _, rec := range table.Records {
		lines[idx] = rec
		idx++
	}

	return lines
}
