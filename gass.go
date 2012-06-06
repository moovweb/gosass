package main

import (
  "os"
  // "fmt"
  sass "gosass"
)

func main() {
  args := os.Args

  ctx := sass.FileContext {
    Options: sass.Options {
      OutputStyle: sass.NESTED_STYLE,
      IncludePaths: nil,
    },
    InputPath:    args[1],
    OutputString: "",
    ErrorStatus:  0,
    ErrorMessage: "",
  }

  sass.CompileFile(&ctx)

}
