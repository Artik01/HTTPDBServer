
package main

import (
	"net/http"
	"io"
	"bytes"
	"fmt"
	"strings"
	"os"
	"bufio"
)

type (
	Person struct {
		Name         string `json:"name" xml:"name"`
		Surname      string `json:"surname" xml:"surname"`
		PersonalCode string `json:"personalCode" xml:"personalCode"`
	}

	Teacher struct {
		ID        string   `json:"id" xml:"id"`
		Subject   string   `json:"subject" xml:"subject"`
		Salary    float64  `json:"salary" xml:"salary"`
		Classroom []string `json:"classroom" xml:"classroom>value"`
		Person    `json:"person"`
	}

	Student struct {
		ID     string `json:"id" xml:"id"`
		Class  string `json:"class" xml:"class"`
		Person `json:"person"`
	}

	Staff struct {
		ID        string  `json:"id" xml:"id"`
		Salary    float64 `json:"salary" xml:"salary"`
		Classroom string  `json:"classroom" xml:"classroom"`
		Phone     string  `json:"phone" xml:"phone"`
		Person    `json:"person"`
	}
)

func main() {
	for {
		fmt.Print("Request:")
		reader := bufio.NewReader(os.Stdin)
		reqStr, err := reader.ReadString('\n')
		if err != nil {panic(err)}
		
		client := &http.Client{}
		var body bytes.Buffer
		switch {
			case strings.HasPrefix(reqStr, "POST"):
				fileName := strings.Split(reqStr, " ")[1]
				fileName = fileName[:len(fileName)-2]
				f, err := os.Open(fileName)
				if err != nil {panic(err)}
				data, err := io.ReadAll(f)
				if err != nil {panic(err)}
				f.Close()
				body.Write(data)
			case strings.HasPrefix(reqStr, "DELETE"):
				id := []byte(strings.Split(reqStr, " ")[1])
				body.Write(id[:len(id)-2])
			default:
				reqStr = reqStr[:len(reqStr)-2]
		}
		
		cmd := strings.Split(reqStr, " ")[0]
		
		req, err := http.NewRequest(cmd, "http://localhost:8080/", &body)
		if err != nil {panic(err)}
		
		resp, err := client.Do(req)
		if err != nil {panic(err)}
		
		data, err := io.ReadAll(resp.Body)
		if err != nil {panic(err)}
		resp.Body.Close()
		
		fmt.Printf("%s\n", data)
	}
}
