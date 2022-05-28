package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"

	"h12.io/socks"
)

type IP struct {
	IP string
}

var (
	wg    sync.WaitGroup
	valid []string
)

func main() {
	runtime.GOMAXPROCS(50)
	var timeout string
	fmt.Printf("Enter timeout in s")
	fmt.Scan(&timeout)

	file := "proxies.txt"
	reader, _ := os.Open(file)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		wg.Add(1)
		go CheckSocks5(scanner.Text(), &wg, timeout)
	}
	wg.Wait()
	for _, v := range valid {
		fmt.Println(v)
	}
}

func CheckSocks5(proxyy string, wg *sync.WaitGroup, timeout string) (err error) {
	defer wg.Done()
	dialSocksProxy := socks.Dial("socks5://" + proxyy + "?timeout=" + string(timeout) + "s")
	tr := &http.Transport{Dial: dialSocksProxy}
	httpClient := &http.Client{Transport: tr}
	resp, err := httpClient.Get("http://checkip.amazonaws.com")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("err")
		return err
	}
	valid = append(valid, proxyy)
	return nil
}
