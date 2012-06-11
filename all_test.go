package gosass

import (
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
const testFileName      = "test.scss"

func compileTest(t *testing.T, fileName string) (result string) {

  ctx := FileContext {
    Options: Options {
      OutputStyle: NESTED_STYLE,
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
      t.Error("ERROR: ", ctx.ErrorMessage)
    } else {
      t.Error("UNKNOWN ERROR")
    }
  } else {
    result = ctx.OutputString
  }

  return result
}

func TestConcurrent(t *testing.T) {
  testFunc := func(done chan bool) {
    done <- false
    compileTest(t, testFileName)
    done <- true
  }
  runParallel(testFunc, numConcurrentRuns)
}