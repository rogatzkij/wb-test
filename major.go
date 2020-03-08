package main

import (
	"bufio"
	"os"
	"sync"
)

const GORUTINE_LIMITS = 5
const PATTERN = "go"

/*
Главная горутина создает необходимые каналы для взаимодействия между горутинами и запускает анализирующую горутину. Далее главная горутина последовательно считывает входные данные из stdin, запускает необходимое число горутин-воркеров, формирует задачу и отправляет ее в канал задач.
...
При достижении конца входных данных главная горутина закрывает канал задач, передавая сигнал завершения для горутин-воркеров. Каждая горутина-воркер при завершении отправляет сигнал в специальный канал. По-скольку главная горутина знает число запущенных горутин-воркеров, она ждет до тех пор, пока каждый воркер пришлет сигнал завершения. Далее главная горутина дожидается завершения анализирующей горутины и завершает работу.
*/
func Major(wg *sync.WaitGroup) {
	defer wg.Done()

	chanLimits := make(chan bool, GORUTINE_LIMITS) // Канал ограничивает максималльное кол-во запущенных горутин
	wgWorker := &sync.WaitGroup{}                  // WaitGroup для Worker'ов

	wgAnalyzer := &sync.WaitGroup{} // WaitGroup для Analyzer'ов
	wgAnalyzer.Add(1)

	chanAnalyzerData := make(chan *AnalyzerData, GORUTINE_LIMITS) // Каннал для данных о работе Worker'ов

	go Analyzer(wgAnalyzer, chanAnalyzerData) // Запускаем горутину-анализатор

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() { // Считываем данные с stdin
		site := scanner.Text()

		chanLimits <- true
		wgWorker.Add(1)

		go func(chanLimits chan bool) {
			Worker(wgWorker, site, PATTERN, chanAnalyzerData)
			_ = <-chanLimits
		}(chanLimits)
	}

	wgWorker.Wait()
	close(chanAnalyzerData)

	wgAnalyzer.Wait()
}
