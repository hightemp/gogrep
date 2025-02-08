package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

var (
	filesBufferSize        int = 10
	linesBufferSize        int = 1000
	searchInFileNumWorkers int = 4
)

func isTextFile(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	sample := make([]byte, 512)
	n, err := file.Read(sample)
	if err != nil && err != io.EOF {
		return false
	}
	sample = sample[:n]

	mimeType := http.DetectContentType(sample)
	if strings.HasPrefix(mimeType, "text/") {
		return true
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != "" {
		mimeType := mime.TypeByExtension(ext)
		if strings.HasPrefix(mimeType, "text/") {
			return true
		}
	}

	if bytes.Contains(sample, []byte{0}) {
		return false
	}

	textChars := 0
	for _, b := range sample {
		if (b >= 32 && b < 127) || b == '\n' || b == '\r' || b == '\t' {
			textChars++
		}
	}

	return float64(textChars)/float64(len(sample)) > 0.9
}

type FilesSearchWorkerPool struct {
	FilesCh       chan string
	WorkersNumber int
	regex         *regexp.Regexp
	Wg            *sync.WaitGroup
}

func (fswp *FilesSearchWorkerPool) Start() {
	for i := 0; i < fswp.WorkersNumber; i++ {
		fswp.Wg.Add(1)
		go fswp.Work()
	}
}

func (fswp *FilesSearchWorkerPool) Work() {
	for filePath := range fswp.FilesCh {
		if !isTextFile(filePath) {
			continue
		}
		SearchInFile(fswp.regex, filePath)
	}
	fswp.Wg.Done()
}

type LineInfo struct {
	number int
	text   string
}

func SearchInFile(regex *regexp.Regexp, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	linesCh := make(chan LineInfo, linesBufferSize)
	matchesCh := make(chan LineInfo, linesBufferSize)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(linesCh)

		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			linesCh <- LineInfo{lineNum, scanner.Text()}
		}
	}()

	for i := 0; i < searchInFileNumWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range linesCh {
				if regex.MatchString(line.text) {
					matchesCh <- line
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(matchesCh)
	}()

	for match := range matchesCh {
		fmt.Printf("%s:%d:%s\n", filePath, match.number, match.text)
	}
}

func SearchInFiles(regex *regexp.Regexp, files []string) {
	fswp := &FilesSearchWorkerPool{
		FilesCh:       make(chan string, filesBufferSize),
		WorkersNumber: filesBufferSize,
		regex:         regex,
		Wg:            &sync.WaitGroup{},
	}

	fswp.Start()

	for _, filePath := range files {
		fswp.FilesCh <- filePath
	}
	close(fswp.FilesCh)

	fswp.Wg.Wait()
}

func SearchInDirRecursive(regex *regexp.Regexp, dir string) {
	fswp := &FilesSearchWorkerPool{
		FilesCh:       make(chan string, filesBufferSize),
		WorkersNumber: filesBufferSize,
		regex:         regex,
		Wg:            &sync.WaitGroup{},
	}

	fswp.Start()

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fswp.FilesCh <- path
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory: %v\n", err)
		return
	}
	close(fswp.FilesCh)

	fswp.Wg.Wait()
}

func SearchInStdin(regex *regexp.Regexp) {
	scanner := bufio.NewScanner(os.Stdin)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if regex.MatchString(line) {
			fmt.Printf("%d:%s\n", lineNumber, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
	}
}

func main() {
	var recursiveSearchFlag bool

	flag.BoolVar(&recursiveSearchFlag, "r", false, "Search recursive in directory")

	flag.IntVar(&filesBufferSize, "fbs", 100, "Size of file buffer")
	flag.IntVar(&linesBufferSize, "lbs", 1000, "Size of lines buffer")
	flag.IntVar(&searchInFileNumWorkers, "snw", runtime.NumCPU()*2, "Number of workers for searching in file")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: gogrep [options] pattern [file...]")
		os.Exit(1)
	}

	pattern := flag.Arg(0)

	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Invalid regular expression:", err)
		os.Exit(1)
	}

	if recursiveSearchFlag {
		dir := flag.Arg(1)
		SearchInDirRecursive(regex, dir)
	} else {
		files := flag.Args()[1:]
		if len(files) == 0 {
			SearchInStdin(regex)
		} else {
			SearchInFiles(regex, files)
		}
	}
}
