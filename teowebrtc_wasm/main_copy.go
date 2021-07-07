// Example from: https://github.com/golang/go/wiki/WebAssembly#getting-started
//
/* Build:

GOOS=js GOARCH=wasm go build -o main.wasm
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .

# install goexec: go get -u github.com/shurcooL/goexec
goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'

*/

package main

import "fmt"

func Main_copy() {
	fmt.Println("Hello, WebAssembly!")
}
