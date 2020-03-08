package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

/*
. Запущенные воркеры берут задачу из канала, выполняют ее и отправляют в канал выполненных задач, в котором анализирующая горутина считает статистику и формирует выходные данные программы.
*/
func Worker(wg *sync.WaitGroup, site, pattern string, result chan<- *AnalyzerData) {
	defer wg.Done()

	data := &AnalyzerData{
		site: site,
	}

	resp, err := http.Get(site)
	if err != nil {
		data.err = err
		result <- data
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		data.err = err
		result <- data
		return
	}

	err = resp.Body.Close()
	if err != nil {
		data.err = err
		result <- data
		return
	}

	var lines []string
	count := 0

	for _, line := range strings.Split(string(body), "\n") {
		if strings.Contains(line, pattern) {
			count++
			lines = append(lines, line)
		}
	}

	data.count = count
	data.result = lines

	result <- data
}
