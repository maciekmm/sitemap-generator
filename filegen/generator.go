package filegen

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/maciekmm/sitemap-generator/config"
)

type Generator struct {
	config    config.Config
	Input     chan string
	output    []*output
	waitGroup *sync.WaitGroup
}

func NewGenerator(cfg config.Config, waitGroup *sync.WaitGroup) (*Generator, error) {
	fileGenerator := &Generator{config: cfg, Input: make(chan string, 1024), waitGroup: waitGroup}
	//TODO: Make this customizable
	os.MkdirAll("./output", 0755)
	for _, filter := range cfg.Output {
		log.Println(filter.Regex)
		out := &output{*filter, nil, 0, 0}
		err := out.nextFile()
		if err != nil {
			return nil, err
		}
		fileGenerator.output = append(fileGenerator.output, out)
	}
	return fileGenerator, nil
}

type url struct {
	Location         string     `xml:"loc"`
	LastModification *time.Time `xml:"lastmod,omitempty"`
	config.Changeable
}

//Start starts a filegenerator loop
func (g *Generator) Start() {
	for {
		select {
		case job, ok := <-g.Input:
			if !ok {
				log.Println("FileGenerator: Stopping")
				for _, out := range g.output {
					out.clean()
				}
				g.waitGroup.Done()
				return
			}
			for _, out := range g.output {
				if out.Regex.MatchString(job) {
					out.put(job)
					break
				}
			}
			g.waitGroup.Done()
		}
	}
}
