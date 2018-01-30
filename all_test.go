package gosass

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"testing"
)

func runParallel(testFunc func(chan bool), concurrency int) {
	runtime.GOMAXPROCS(4)
	done := make(chan bool, concurrency)
	for i := 0; i < concurrency; i++ {
		go testFunc(done)
	}
	for i := 0; i < concurrency; i++ {
		<-done
		<-done
	}
	runtime.GOMAXPROCS(1)
}

const numConcurrentRuns = 200

// const testFileName1     = "test1.scss"
// const testFileName2     = "test2.scss"
// const desiredOutput     = "div {\n  color: black; }\n  div span {\n    color: blue; }\n"

func compileTest(t *testing.T, fileName string) (result string) {
	ctx := FileContext{
		Options: Options{
			OutputStyle:  NESTED_STYLE,
			IncludePaths: make([]string, 0),
		},
		InputPath:    fileName,
		OutputString: "",
		ErrorStatus:  0,
		ErrorMessage: "",
	}

	CompileFile(&ctx)

	if ctx.ErrorStatus != 0 {
		if ctx.ErrorMessage != "" {
			t.Fatal("ERROR: ", ctx.ErrorMessage)
		} else {
			t.Fatalf("Unknown error, status %d", ctx.ErrorStatus)
		}
	}

	return ctx.OutputString
}

const numTests = 3 // TO DO: read the test dir and set this dynamically

// NOTE: This test sometimes fails or deadlocks due to concurrency issues. It is unclear if
// the test was *meant* to test concurrent behavior of gosass or that it was simply
// a way to concurrently (and thus faster) test gosass.
func TestConcurrent(t *testing.T) {
	testFunc := func(done chan bool) {
		done <- false
		for i := 1; i <= numTests; i++ {
			inputFile := fmt.Sprintf("test/test%d.scss", i)
			desiredOutput, err := ioutil.ReadFile(fmt.Sprintf("test/test%d.css", i))
			if err != nil {
				t.Fatalf("ERROR: couldn't read test/test%d.css: %v", i, err)
			}
			result := compileTest(t, inputFile)
			if result != string(desiredOutput) {
				t.Fatalf("ERROR: incorrect output [%q] != expected [%q]", result, string(desiredOutput))
			}
		}
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)
}

var testSassFuncs = []struct {
	name     string
	context  Context
	expected string
}{
	{name: "image-url test with scheme-less uri path",
		context: Context{
			Options: Options{
				OutputStyle:    COMPRESSED_STYLE,
				SourceComments: false,
				IncludePaths:   []string{"//include-root"},
				ImagePath:      "//image-root"},
			SourceString: "bg { src: image-url('fonts/fancy.otf?whynot');}"},
		expected: "bg{src:url(\"//image-root/fonts/fancy.otf?whynot\");}"},
	{name: "image-url test with relative uri path",
		context: Context{
			Options: Options{
				OutputStyle:    COMPRESSED_STYLE,
				SourceComments: false,
				IncludePaths:   []string{"//include-root"},
				ImagePath:      "/image-root"},
			SourceString: "bg { src: image-url('fonts/fancy.otf?whynot');}"},
		expected: "bg{src:url(\"/image-root/fonts/fancy.otf?whynot\");}"},
}

func TestSassFunctions(t *testing.T) {
	for _, tobj := range testSassFuncs {
		Compile(&tobj.context)

		if tobj.context.ErrorStatus != 0 {
			if tobj.context.ErrorMessage != "" {
				t.Fatal("ERROR: ", tobj.context.ErrorMessage)
			} else {
				t.Fatalf("Unknown error, status %d", tobj.context.ErrorStatus)
			}
		} else if tobj.context.OutputString != tobj.expected {
			t.Fatalf("Test case %s failed.  Expected \"%s\" but received \"%s\".", tobj.name, tobj.expected, tobj.context.OutputString)
		}
	}
}
