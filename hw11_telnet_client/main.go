package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func createClient() (TelnetClient, error) {
	flag.Parse()
	if len(os.Args) != flag.NArg()+2 {
		return nil, errors.New("wrong arguments number")
	}

	address := os.Args[flag.NArg()]
	portStr := os.Args[flag.NArg()+1]
	_, err := strconv.Atoi(portStr)

	if err != nil {
		return nil, errors.New("wrong port number")
	}

	client := NewTelnetClient(address+":"+portStr, timeout, os.Stdin, os.Stdout)
	return client, nil
}

func main() {
	client, err := createClient()
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	wg := sync.WaitGroup{}

	// цикл выполнения действий до первой ошибки
	// если эта ошибка EOF - выводим сообщение и закрываем соединение
	actionLoop := func(action func() error, EOFMessage string) {
		defer wg.Done()
		for {
			err := action()
			if err == nil {
				continue
			}

			if errors.Is(err, io.EOF) {
				fmt.Fprintln(os.Stderr, EOFMessage)
				client.Close()
			}
			return
		}
	}

	wg.Add(2)

	// Цикл получения сообщений
	go actionLoop(client.Receive, "...Connection closed by peer. Press <Enter> to exit")
	// Цикл отправки сообщений
	go actionLoop(client.Send, "...EOF")

	// Ждём завершения обоих циклов
	wg.Wait()
}
