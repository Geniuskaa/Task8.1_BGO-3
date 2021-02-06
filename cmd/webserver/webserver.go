package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	if err := execute(); err != nil {
		os.Exit(1)
	}
}

type transaction struct {
	Id string
	From string
	To string
	Amount int64
}

func execute() (err error) {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			log.Println(cerr)
			if err == nil {
				err = cerr
			}
		}
	}()

	for  {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		conn.SetReadDeadline(time.Now().Add(time.Second * 30))
		conn.SetWriteDeadline(time.Now().Add(time.Minute * 5))
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Println(cerr)
		}
	}()

	reader := bufio.NewReader(conn)
	const delim = '\n'
	line, err := reader.ReadString(delim)
	if err != nil {
		if err != io.EOF {
			log.Println(err)
		}
		log.Printf("received: %s\n", line)
		return
	}
	log.Printf("received: %s\n", line)

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		log.Printf("invalid request line: %s", line)
		return
	}

	path := parts[1]

	switch path {
	case "/":
		err = writeIndex(conn)
	case "/application/json":
		err = writeOperations(conn, "json")
	case "/application/xml":
		err = writeOperations(conn, "xml")
	default:
		err = write404(conn)
	}
	if err != nil {
		log.Println(err)
		return
	}
}

func writeIndex(writer io.Writer) error {
	username := "Василий"
	balance := "1 000.50"

	page, err := ioutil.ReadFile("web/template/index.html")
	if err != nil {
		log.Println(err)
		return err
	}

	page = bytes.ReplaceAll(page, []byte("{username}"), []byte(username))
	page = bytes.ReplaceAll(page, []byte("{balance}"), []byte(balance))

	return writeResponse(writer, 200, []string{
		"Content-Type: text/html;charset=utf-8",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}


func writeOperations(writer io.Writer, typeOfApplication string) error { // generate JSON //////////////////////////////////////////////

	strings.ToLower(typeOfApplication)
	var page []byte
	switch typeOfApplication {
	case "json":
		page = []byte("\"id\":\"1\",\"from\":\"0001\",\"to\":\"0002\",\"amount\":10000,\"created\":1598613478\n")
	case "xml":
		page, _ = xml.Marshal(transaction{
			Id:     "20",
			From:   "0001",
			To:     "0002",
			Amount: 10_203,
		})
	}

	return writeResponse(writer, 200, []string{
		"Content-Type: text/" + typeOfApplication,
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}


func write404(writer io.Writer) error {
	page, err := ioutil.ReadFile("web/template/index.html")
	if err != nil {
		_, err = ioutil.ReadFile("web/template/index.html")
	}

	return writeResponse(writer, 200, []string{
		"Content-Type: text/html;charset=utf-8",
		fmt.Sprintf("Content-Length: #{len(page)}"),
		"Connection: close",
	}, page)
}


func writeResponse(
	writer io.Writer,
	status int,
	headers []string,
	content []byte,
) error {
	const CRLF = "\r\n"
	var err error

	w := bufio.NewWriter(writer)
	_, err = w.WriteString(fmt.Sprintf("HTTP/1.1 %d OK%s", status, CRLF))
	if err != nil {
		log.Println(err)
		return err
	}

	for _, h := range headers {
		_, err = w.WriteString(h + CRLF)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	_, err = w.WriteString(CRLF)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = w.Write(content)
	if err != nil {
		log.Println(err)
		return err
	}

	err = w.Flush()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}



