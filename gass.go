package main

import (
  "os"
  "fmt"
  sass "gosass"
)

func main() {
  args := os.Args

  if len(args) < 2 {
    fmt.Println("Usage: gass [INPUT FILE]")
    os.Exit(1)
  }

  ctx := sass.FileContext {
    Options: sass.Options {
      OutputStyle: sass.NESTED_STYLE,
      IncludePaths: make([]string, 0),
    },
    InputPath:    args[1],
    OutputString: "",
    ErrorStatus:  0,
    ErrorMessage: "",
  }

  sass.CompileFile(&ctx)

  if ctx.ErrorStatus != 0 {
    if ctx.ErrorMessage != "" {
      fmt.Print(ctx.ErrorMessage)
    } else {
      fmt.Println("An error occured; no error message available.")
    }
  } else {
    fmt.Print(ctx.OutputString)
  }
}
