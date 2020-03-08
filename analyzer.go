package main

import (
	"fmt"
	"sync"
)

type AnalyzerData struct {
	site   string
	count  int
	result []string
	err    error
}

/*. Запущенные воркеры берут задачу из канала, выполняют ее и отправляют в канал выполненных задач, в котором анализирующая горутина считает статистику и формирует выходные данные программы. */
func Analyzer(wg *sync.WaitGroup, analyzerData chan *AnalyzerData) {
	defer wg.Done()

	var result []string

LOOP:
	for {

		data, ok := <-analyzerData
		if !ok {
			break LOOP
		}

		if data.err != nil {
			fmt.Printf("Error happends with %s: %s\n", data.site, data.err.Error())
		} else {
			fmt.Printf("Count for %s: %d\n", data.site, len(data.result))
			for _, line := range data.result {
				fmt.Printf("\t%s\n", line)
			}
			result = append(result, data.result...)
		}
	}

	fmt.Printf("Total: %d\n", len(result))
}
