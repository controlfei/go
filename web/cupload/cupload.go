package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func main()  {
	taeget_url := "http://localhost:9090/upload"
	filename := "./text.txt"
	postFile(filename,taeget_url)
}

func postFile(filename ,target_url string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter,err  := bodyWriter.CreateFormFile("uploadfile",filename)
	if err != nil {
		fmt.Println("error writing to buffer ")
		return err
	}

	//打开文件具柄
	fh,err := os.Open(filename)
	if err != nil {
		fmt.Println("openging file failed")
		return err
	}

	defer fh.Close()

	_,err = io.Copy(fileWriter,fh)
	if err != nil {
		return err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp,err := http.Post(target_url,contentType,bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body , err  := 	ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}
