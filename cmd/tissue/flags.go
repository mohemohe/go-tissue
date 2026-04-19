package main

import (
	"flag"
	"fmt"
	"os"
)

// setUsage は flag.FlagSet の --help 出力に位置引数を含む usage 行を付け足す。
func setUsage(fs *flag.FlagSet, usage string) {
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: "+usage)
		fs.PrintDefaults()
	}
}

// parseMixed は、Go 標準 flag が最初の位置引数で解析を打ち切る挙動を回避し、
// フラグと位置引数が任意の順序で混在していてもすべて処理する。
// 戻り値は検出された位置引数 (順序保存)。
func parseMixed(fs *flag.FlagSet, args []string) []string {
	var positional []string
	for {
		_ = fs.Parse(args)
		if fs.NArg() == 0 {
			break
		}
		positional = append(positional, fs.Arg(0))
		args = fs.Args()[1:]
	}
	return positional
}
