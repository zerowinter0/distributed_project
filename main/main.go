package main

import (
	"example_pkg"
	"flag"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "./cfg.toml", "a cfg file path")
	flag.Parse()

	example_pkg.NewNode(path)
}
