package main

import (
	"fmt"
	"time"
	"os"
	"path/filepath"
	"bufio"
	"regexp"
	"unicode/utf8"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main(){
	var tmp string
	const aa_dir_path = "C:\\Live5ch\\aa"

	backup(aa_dir_path)

	convAa(aa_dir_path)

	fmt.Print("終わったので何かキー押してください")
	fmt.Scanln(&tmp)
}

// バックアップ
func backup(aa_dir_path string){
	today := time.Now()
	const layout = "20060102_150405";  // YYYYMMDD

	backup_dir_name := aa_dir_path + today.Format(layout) + "_bk"

	os.Mkdir(backup_dir_name, 0777)

	fmt.Println("mkdir " + backup_dir_name)

	// ファイルをコピーする
	files , _ := filepath.Glob(aa_dir_path + "\\*.txt")
	for _ , f := range files {
		_ = os.Link(f, backup_dir_name + "\\" + filepath.Base(f))
		fmt.Println("copy " + f)
	}

	fmt.Println("---backup ok ")
	fmt.Println("")

}

// 変換
func convAa(aa_dir_path string){
	var fp *os.File
	var fpw *os.File
	// var err error

	// ファイルの中身をコンバートする
	files , _ := filepath.Glob(aa_dir_path + "\\*.txt")
	for _ , f := range files {
		// aa文字列初期化
		aa_text := ""

		// 読み込みファイルオープン
		fp, _ = os.Open(f)

		//
		scan := bufio.NewScanner(transform.NewReader(fp, japanese.ShiftJIS.NewDecoder()))
		//scan := bufio.NewScanner(fp)

		for scan.Scan() {
			// 1行読み込み
			buf := scan.Text()

			// 変換処理
			// buf = utf8mb4_encode_numericentity(buf, `[　 \n][　 ]{3}[　 \n]`)
			buf = utf8mb4_encode_numericentity(buf, `[　 ]+?`)
			buf = utf8mb4_encode_numericentity(
				buf, "[${ .,'`\"\\:;\\-=~_|＼／<>}\n][${ .,'`\"\\:;\\-=~_|＼／<>}\n}]{3}[${ .,'`\"\\:;\\-=~_|＼／<>}\n]")

			// 連結
			aa_text += buf + "\r\n"
		}

		// ファイルクローズ
		fp.Close()

		// 書き出しでファイルオープン
		fpw, _ = os.Create(f)
		w := bufio.NewWriter(transform.NewWriter(fpw, japanese.ShiftJIS.NewEncoder()))

		// 書き込み
		fmt.Fprint(w, aa_text)
		w.Flush()

		// ファイルクローズ
		fpw.Close()
	}

	fmt.Println("---convert ok ")
	fmt.Println("")
}

func utf8mb4_encode_numericentity(str string, reg string) string {
	re := regexp.MustCompile(reg);
	return re.ReplaceAllStringFunc(str, func(match string) string {
        r, _ := utf8.DecodeLastRuneInString(match)
        return fmt.Sprintf("&#%d;", r)
    });
}


