package parser

import (
	"fmt"
	"github.com/google/pprof/profile"
	snapshot "gold/api"

	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func NewParser() *Parser {
	return &Parser{logger: log.Logger{}, wg: sync.WaitGroup{}}
}

func (p *Parser) Parse(snapshot *snapshot.SnapShot) (*ParsedSnapShot, error) {
	if snapshot == nil {
		p.logger.Println("GoroutineDump:", ErrBufEmpty.Error())
		return nil, ErrBufEmpty
	}
	var (
		goroutiness *[]Goroutine
		blocks      *[]Block
		mutex       *[]Mutex
	)
	p.runAsync("GoroutineDump", func() error {
		var err error
		goroutiness, err = p.parseGoroutineDump(&snapshot.GoroutineDump, &p.wg)
		return err
	})

	p.runAsync("BlockProfile", func() error {
		var err error
		blocks, err = p.parseBlockProfile(&snapshot.BlockProfile, &p.wg)
		return err
	})

	p.runAsync("MutexProfile", func() error {
		var err error
		mutex, err = p.parseMutexProfile(&snapshot.MutexProfile, &p.wg)
		return err
	})

	p.wg.Wait()
	result := &ParsedSnapShot{
		snapshot.ServiceName,
		snapshot.Timestamp,
		*goroutiness,
		*blocks,
		*mutex,
		snapshot.RuntimeMetrics,
		snapshot.Version,
	}
	return result, nil
}

func (p *Parser) parseGoroutineDump(buf *[]byte, wg *sync.WaitGroup) (*[]Goroutine, error) {
	defer wg.Done()
	raw := string((*buf)[:len(*buf)-1])
	stacks := strings.Split(raw, "\n")
	var Goroutines []Goroutine
	var goroutine Goroutine
	for _, v := range stacks {

		if v == "" {
			continue
		}

		if strings.Contains(v, "created by") {
			continue
		}
		v = strings.Replace(v, "+0x", "", -1)
		if strings.Contains(v, "(") {
			v = func(string) string {
				re := regexp.MustCompile(`\([^()]*\)$`)
				return re.ReplaceAllString(v, "")
			}(v)

		}
		if strings.Contains(v, "goroutine") {

			re := regexp.MustCompile(`^goroutine\s+(\d+)`)
			matches := re.FindStringSubmatch(v)
			if len(matches) > 1 {
				id, err := strconv.ParseUint(matches[1], 10, 64)
				if err == nil {
					goroutine.Id = id
				} else {
					log.Printf("failed to parse goroutine id: %v", err)
				}
			}
			v = re.ReplaceAllString(v, "")
			if len(goroutine.Data) != 0 {

				Goroutines = append(Goroutines, goroutine)
				goroutine = Goroutine{}
			}
		}

		goroutine.Data = append(goroutine.Data, v)

	}
	if goroutine.Data != nil {
		Goroutines = append(Goroutines, goroutine)
	}

	return &Goroutines, nil
}

func (p *Parser) parseBlockProfile(buf *[]byte, wg *sync.WaitGroup) (*[]Block, error) {
	defer wg.Done()
	if buf == nil {
		p.logger.Println("BlockProfile:", ErrBufEmpty.Error())
		return nil, ErrBufEmpty
	}
	data, err := profile.ParseData(*buf)
	if err != nil {
		return nil, err
	}
	var out []Block
	for _, sample := range data.Sample {
		b := Block{
			Count:  uint64(sample.Value[0]),
			Cycles: sample.Value[1],
			Frames: []string{},
		}

		// Проходим по стеку
		for _, loc := range sample.Location {
			for _, line := range loc.Line {
				if line.Function != nil {
					frame := fmt.Sprintf("%s:%d", line.Function.Name, line.Line)
					b.Frames = append(b.Frames, frame)
				}
			}
		}

		out = append(out, b)
	}

	return &out, nil
}
func (p *Parser) parseMutexProfile(buf *[]byte, wg *sync.WaitGroup) (*[]Mutex, error) {
	defer wg.Done()
	if buf == nil {
		p.logger.Println("MutexProfile:", ErrBufEmpty.Error())
		return nil, ErrBufEmpty
	}
	data, err := profile.ParseData(*buf)
	if err != nil {
		return nil, err
	}
	var out []Mutex

	for _, sample := range data.Sample {
		m := Mutex{
			Count:  uint64(sample.Value[0]),
			Cycles: sample.Value[1],
			Frames: []string{},
		}
		for _, loc := range sample.Location {
			for _, line := range loc.Line {
				if line.Function != nil {
					frame := fmt.Sprintf("%s:%d", line.Function.Name, line.Line)
					m.Frames = append(m.Frames, frame)
				}
			}
		}
		out = append(out, m)
	}
	return &out, nil
}
func (p *Parser) runAsync(name string, task func() error) {
	go func() {
		if err := task(); err != nil {
			p.logger.Printf("%s: %v\n", name, err)
		}
	}()
}
