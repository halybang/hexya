// Copyright 2019 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/hexya-erp/hexya/src/tools/generate"
	"github.com/hexya-erp/hexya/src/tools/strutils"
)

func main() {
	if len(os.Args) <= 1 {
		panic("You must provide a csv file to load")
	}
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	cReader := csv.NewReader(file)
	headers, err := cReader.Read()
	if err != nil {
		panic(err)
	}
	var res []map[string]string
	for {
		record, err := cReader.Read()
		if err == io.EOF {
			break
		}
		recMap := make(map[string]string)
		for i := 0; i < len(headers); i++ {
			recMap[headers[i]] = record[i]
		}
		dtSubst := map[string]string{
			"%d": "02",
			"%m": "01",
			"%Y": "2006",
			"%y": "06",
			"%H": "15",
			"%I": "03",
			"%M": "04",
			"%S": "05",
			"%p": "PM",
			"%b": "Jan",
			"%B": "January",
			"%A": "Monday",
			"%a": "Mon",
		}
		recMap["date_format"] = recMap["date_format"]
		recMap["time_format"] = recMap["time_format"]
		recMap["date_format_go"] = strutils.Substitute(recMap["date_format"], dtSubst)
		recMap["time_format_go"] = strutils.Substitute(recMap["time_format"], dtSubst)
		dir := recMap["direction"]
		recMap["direction"] = "LanguageDirectionLTR"
		if dir == "Right-to-Left" {
			recMap["direction"] = "LanguageDirectionRTL"
		}
		trans, err := strconv.ParseBool(recMap["translatable"])
		if err != nil {
			panic(err)
		}
		recMap["translatable"] = fmt.Sprintf("%t", trans)
		grp := strings.TrimSuffix(strings.TrimPrefix(recMap["grouping"], "["), "]")
		grps := strings.Split(grp, ",")
		recMap["grouping"] = fmt.Sprintf("NumberGrouping{%s}", strings.Join(grps, ", "))
		res = append(res, recMap)
	}
	generate.CreateFileFromTemplate("locales.go", tmpl, res)
}

var tmpl = template.Must(template.New("").Parse(`
// This file has been generated by generate_locales

package i18n

// locales lists all available locales by ISO code.
var locales = map[string]*Locale{
{{- range . }}
	"{{ .iso_code }}": {
		Name: "{{ .name }}",
		Code: "{{ .code }}",
		ISOCode: "{{ .iso_code }}",
		DateFormat: "{{ .date_format }}",
		TimeFormat: "{{ .time_format }}",
		DateFormatGo: "{{ .date_format_go }}",
		TimeFormatGo: "{{ .time_format_go }}",
		DecimalPoint: "{{ .decimal_point }}",
		ThousandsSep: "{{ .thousands_sep }}",
		Grouping: {{ .grouping }},
	},
{{- end }}
}
`))
