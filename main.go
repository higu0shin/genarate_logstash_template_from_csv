package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var path string
	var sc bufio.Scanner
	//ファイル名の取得。
	// 引数にパス指定が無ければ入力させる。
	if len(os.Args) < 2 {
		//fmt.Println("CSVファイルのパスを入力してください。(ex. /home/hogehoge/sample.csv)")
		sc = *bufio.NewScanner(os.Stdin)

	} else {
		// 引数にパス指定があればそれを利用。
		path = os.Args[1]
		path, _ = filepath.Abs(path)
		file, err := os.Open(path)
		if err != nil {
			// Openエラー
			fmt.Println("ファイルを開けませんでした。パス・権限を確認してください。")
			os.Exit(-1)
		}
		defer file.Close()
		sc = *bufio.NewScanner(file)
	}

	columns := getFirstLine(&sc)

	logstashInput := "input{\n  file{\n    path => \"" + path + "\"\n    start_position => \"beginning\"\n    sincedb_path => \"/dev/null\"\n    #sincedb_path => \"nul\"\n  }\n}\n"
	logstashFilter := "filter{\n  csv{\n    columns=> [" + printHeader(columns) + "]\n  }\n}\n"
	logstashOutput := "output{\n  #stdout{\n  #  codec => rubydebug\n  #}\n  elasticsearch{\n    hosts => [\"localhost:9200\"]\n    index => \"my_index\"\n    user => \"elastic\"\n    password => \"changeme\"\n  }\n}\n"

	fmt.Print(logstashInput)
	fmt.Print(logstashFilter)
	fmt.Print(logstashOutput)
}

func printHeader(columns string) string {
	columnSlice := strings.Split(columns, ",")
	var result string
	//headerの内容を出力する
	for i := range columnSlice {
		result += "\""
		result += columnSlice[i]
		result += "\""
		if i == len(columnSlice)-1 {
			break
		}
		result += ","

	}
	return result
}

func getFirstLine(sc *bufio.Scanner) string {
	/*for i := 1; sc.Scan(); i++ {
		if err := sc.Err(); err != nil {
			// エラー処理
			break
		}
		fmt.Printf("%4d行目: %s\n", i, sc.Text())
	}*/
	sc.Scan()
	if err := sc.Err(); err != nil {
		fmt.Print("不正な入力です。")
		os.Exit(-1)
	}
	columns := sc.Text()
	return columns
}
