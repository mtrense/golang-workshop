# Installing Go

Go can be installed using homebrew on Mac OS
```bash
$ brew install go
```

## GOPATH

After that, we should set the global `GOPATH` to a reasonable value
```bash
$ cat <<EOF >> ~/.bash_profile
export GOPATH=~/.go
export PATH=$PATH:$GOPATH/bin
EOF
$ source ~/.bash_profile
```

This will set the standard `GOPATH` to `~/.go`, which allows us to install certain utilities globally:

```
$ go get -u github.com/nsf/gocode
$ go get -u github.com/lukehoban/go-outline
$ go get -u github.com/newhook/go-symbols
$ go get -u github.com/uudashr/gopkgs/cmd/gopkgs
$ go get -u golang.org/x/tools/cmd/gorename
$ go get -u github.com/fatih/gomodifytags
$ go get -u github.com/cweill/gotests/...
$ go get -u github.com/tylerb/gotype-live
$ go get -u github.com/constabulary/gb/...
$ go get -u github.com/onsi/ginkgo/ginkgo
$ go get -u github.com/onsi/gomega/...
```

`GOPATH` is much like the known `PATH` variable in Unix. It contains paths delimited by colons.
*Be aware that `go get` will always install into the first available path from the `GOPATH` environment variable*

## Utilities

[**gocode**](https://github.com/nsf/gocode) Provides Intellisense and autocompletion for Go programs
[**go-outline**](https://github.com/lukehoban/go-outline) Extracts declarations from the source tree
[**go-symbols**](https://github.com/newhook/go-symbols) Extracts symbols and identifiers from the source tree
[**gopkgs**](https://github.com/uudashr/gopkgs) Lists packages that are available in the current `GOPATH`
[**gorename**](https://godoc.org/golang.org/x/tools/cmd/gorename) Provides type-safe renaming of identifiers in Go
[**gomodifytags**](https://github.com/fatih/gomodifytags) Savely modifies tags on struct fields
[**gotests**](https://github.com/cweill/gotests) Generates sophisticated test stubs from your code
[**gotype-live**](https://github.com/tylerb/gotype-live) Parses and type checks code as it is written
[**gb**](https://getgb.io) A sane build and dependency management tool for Go projects


# Setting up the project

First create a directory somewhere on your disk (this will be our workspace in Go terms):
```bash
$ mkdir -p <PROJECTS>/golang-workshop/src
$ cd <PROJECTS>/golang-workshop
```
(A valid `GOPATH` is any directory that has a subdirectory called `src`)

#### Configuring the `GOPATH`

Now prepend your newly created path to the `GOPATH` variable:
```bash
$ export GOPATH=$(pwd)/vendor:$(pwd):$GOPATH
```

This adds two directories (`vendor` in the current directory and the current directory itself) to the `GOPATH`. You can easily automate that using a tool like for example [`direnv`](https://direnv.net). 

#### Persistent configuration using `direnv`

A proper `.envrc` for a golang project would contain something like:
```bash
export PROJECT_ROOT=$(pwd)
export GOPATH=$PROJECT_ROOT/vendor:$PROJECT_ROOT:$GOPATH
```

Adding the `vendor` directory to your `GOPATH` has the additional benefit that any golang tools that adhere to go's path conventions will work out of the box.

#### Editor configuration

Editors like VS Code usually infer these paths from the environment if possible, but to be explicit about it, you can add the following configuration (for VS Code) to your **workspace** settings:
```json
{
  "go.gopath": "${workspaceRoot}/vendor:${workspaceRoot}"
}
```

## Building the project

First add a simple `main.go` file as placeholder for the server that we will create later:
```bash
$ mkdir $PROJECT_ROOT/src/server
$ cat <<EOF > $PROJECT_ROOT/src/server/main.go
package main

import "fmt"

func main() {
  fmt.Println("Hello World!")
}
EOF
```

Now with a little bit of code in place, we can use `gb` to build that into a binary:
```bash
$ gb build
```

This has created two more folders:
```
bin/
└── server
pkg/
└── darwin-amd64
    └── server.a
```
with `pkg` containing the (statically) compiled library from your code and `bin` containing the linked binary. `gb` names the binary after the directory it is found in.

Executing the resulting binary leads to the expected result:
```bash
$ ./bin/server
Hello World!
```

## Adding dependencies

Let's add the first dependency, [echo](https://github.com/labstack/echo), a web library for go that nicely fits the gap between the http package contained in the standard library and a serious web application. It ships lots of helpers and middlewares for common problems in the web world, so we'll use it as a starting point to dive into webservices with Go. Install the library through the `gb vendor` command:
```bash
$ gb vendor fetch github.com/labstack/echo
```
*Please note that the path appendix of `/...` (as written in the documentation of echo) is not necessary when using `gb`, as the vendor plugin automatically takes care of transitive dependencies.*

Now that we have echo installed, we should see a new `vendor` directory within the project's path:
```
vendor/
├── manifest
└── src
    ├── github.com
    │   └── ...
    └── golang.org
        └── ...
```

Now we can use that library by importing `github.com/labstack/echo` anywhere in our codebase:
```go
import "github.com/labstack/echo"
```

## A simple webservice

Let's replace the contents of the `main` function in `main.go` with something more web-aware:
```go
server := echo.New()
server.GET("/", func(c echo.Context) error {
	return c.String(http.StatusOK, "Hello CLΛRK!")
})
server.Start(":1323")
```

When we run `go run $PROJECT_ROOT/src/server/main.go` now, the server should start on port 1323:
```bash
$ curl http://localhost:1323
Hello CLΛRK!
```

## Testing your code

Go already has testing support right in the standard library. While being quite productive, the `testing` package has some shortages when it comes to readability and modern testing techniques. The Ginkgo project has tackled these problems and provides a testing library that's way more like rspec and modern JUnit versions:

```go
var _ = Describe("Main", func() {
	Context(".helloWorld", func() {
		It("Returns a String containing 'Hello World'", func() {
			...
			Expect(helloWorld(c)).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(Equal("Hello CLΛRK!"))
		})
	})
})
```

To get started with ginkgo, run the bootstrap subcommand in a folder with go sourcecode:
```bash
$ ( cd $PROJECT_ROOT/src/server ; ginkgo bootstrap )
```

This will create a file called `server_suite_test.go` which includes some glue between Ginkgo and the standard testing library. That way any library (especially CI integrations) that expects the standard Go testing library in place will just work as expected.

### Extracting the HTTP handler function for better testability

Right now we have inlined the logic of our web endpoint:
```go
server.GET("/", func(c echo.Context) error {
    return c.String(http.StatusOK, "Hello CLΛRK!")
})
```

While this was easy for starting, it is not easy to test, so let's extract the logic into it's own function:
```go
func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello CLΛRK!")
}
```
And reference that function in the routing code:
```go
server.GET("/", helloWorld)
```

Now we're ready to test that functionality.

### Writing the first test

Using the generate subcommand we can create a new test file:
```bash
$ ( cd $PROJECT_ROOT/src/server ; ginkgo generate main )
```

After adopting the test file to our needs, it looks like this:
```go
package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	Context(".helloWorld", func() {
	
	})
})
```

Now for the test code, first we have to initialize the context of the webservice:
```go
e := echo.New()
req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
rec := httptest.NewRecorder()
c := e.NewContext(req, rec)
```
This creates a new instance of echo, as we did in the webservice itself and from that a new `Context`, that includes a mock-request and a recording response.

Now we can check for the correct behavior:
```go
Expect(helloWorld(c)).To(BeNil())
Expect(rec.Code).To(Equal(http.StatusOK))
Expect(rec.Body.String()).To(Equal("Hello CLΛRK!"))
```
As `helloWorld` returns an error, we expect it to return nil, when everything went fine. The `Response` should then contain the specified values.

