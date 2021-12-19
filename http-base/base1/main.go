package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/",indexHandle)
	http.HandleFunc("/hello",helloHandler)
	log.Fatal(http.ListenAndServe(":9999",nil))
}
//得到访问的url地址
func indexHandle(w http.ResponseWriter,req *http.Request)  {
	//输出url地址
	//URL.Path = "/"
	fmt.Fprintf(w,"URL.Path = %q\n",req.URL.Path)
}
//得到请求头文件
func helloHandler(w http.ResponseWriter,req *http.Request)  {
	//循环输出header文件
	/*
	Header ["Accept-Encoding"] = ["gzip, deflate, br"]
	Header ["Connection"] = ["keep-alive"]
	Header ["Content-Length"] = ["0"]
	Header ["Accept"] = ["application/json"]
	Header ["User-Agent"] = ["PostmanRuntime/7.28.4"]
	 */
	for k,v := range req.Header {
		fmt.Fprintf(w,"Header [%q] = %q\n",k,v)
	}
}
