package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/maciekmm/sitemap-generator"
	"github.com/maciekmm/sitemap-generator/config"
)

func main() {
	configFile := flag.String("config", "", "Path to config file")
	flag.Parse()
	if *configFile == "" {
		fmt.Fprintln(os.Stderr, "Config flag can't be empty")
		os.Exit(1)
	}
	cfg, err := config.FromFile(*configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	gen := sitemapgen.NewGenerator(cfg)
	gen.Start()
}
