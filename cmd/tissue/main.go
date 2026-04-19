package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	name := os.Args[1]
	args := os.Args[2:]
	switch name {
	case "configure":
		cmdConfigure(args)
	case "me":
		cmdMe(args)
	case "checkin":
		cmdCheckin(args)
	case "collection":
		cmdCollection(args)
	case "search":
		cmdSearch(args)
	case "tags":
		cmdTags(args)
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", name)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: tissue <command> [options]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  configure   認証情報の設定")
	fmt.Fprintln(os.Stderr, "  me          自分のユーザー情報を表示")
	fmt.Fprintln(os.Stderr, "  checkin     チェックイン操作 (add/list/get/update/delete)")
	fmt.Fprintln(os.Stderr, "  collection  コレクション操作 (list/create/update/delete/item ...)")
	fmt.Fprintln(os.Stderr, "  search      チェックインを検索")
	fmt.Fprintln(os.Stderr, "  tags        最近使用したタグ")
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		die("failed to encode: %v", err)
	}
}
