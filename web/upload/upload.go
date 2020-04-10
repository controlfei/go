package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main()  {
	http.HandleFunc("/upload",uploads)

	err := http.ListenAndServe(":9080",nil)
	if err != nil {
		log.Fatal("listenandserve",err)
	}
}

func uploads(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("method:",r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		//formatint  返回crutime十进制的字符串
		io.WriteString(w,strconv.FormatInt(crutime,10))
		token := fmt.Sprintf("%v",h.Sum(nil))

		t,_ := template.ParseFiles("upload.gtpl")
		t.Execute(w,token)

	}else{
		r.ParseMultipartForm(32 << 20)
		file , handler ,err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w,"%v",handler.Header)
		f,err  := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE,0666)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer f.Close()

		io.Copy(f,file)
	}
}
