package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)


func main(){
	var path string

	//ファイル名の取得。
	// 引数にパス指定が無ければ入力させる。
	if len(os.Args)<2 {
		fmt.Println("CSVファイルのパスを入力してください。(ex. /home/hogehoge/sample.csv)")
		fmt.Scan(&path)
	}else {
		// 引数にパス指定があればそれを利用。
		path = os.Args[1]
	}
	file, err := os.Open(path)
	if err != nil {
		// Openエラー
		fmt.Println("ファイルを開けませんでした。パス・権限を確認してください。")
		os.Exit(-1)
	}

	defer file.Close()
	columns := get_firdt_line(file)

	logstash_input := "input{\n  file{\n    path => \"" + path + "\"\n    start_position => \"beginning\"\n    sincedb_path => \"/dev/null\"\n    #sincedb_path => \"nul\"\n  }\n}\n"
	logstash_filter := "filter{\n  csv{\n    columns=> [" + print_header(columns) + "]\n  }\n}\n"
	logstash_output := "output{\n  #stdout{\n  #  codec => rubydebug\n  #}\n  elasticsearch{\n    hosts => [\"localhost:9200\"]\n    index => \"my_index\"\n    user => \"elastic\"\n    password => \"changeme\"\n  }\n}\n"

	fmt.Print(logstash_input)
	fmt.Print(logstash_filter)
	fmt.Print(logstash_output)
}

func print_header(columns string) string{
	column_slice := strings.Split(columns, ",")
	var result string
	//headerの内容を出力する
	for i := range (column_slice) {
		result += "\""
		result += column_slice[i]
		result += "\""
		if i == len(column_slice)-1 {
			break
		}
		result += ","

	}
	return result
}

func get_firdt_line(file *os.File) string {
	sc := bufio.NewScanner(file)
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
