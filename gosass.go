package gosass

/*
#cgo LDFLAGS: -Llibsass -lsass -lstdc++
#cgo CFLAGS: -Ilibsass

#include "sass_interface.h"
*/
import "C"

type Options struct {
  OutputStyle  int
  IncludePaths []string
  // eventually gonna' have things like callbacks and whatnot
}

/*
type SourceContext struct {
  Options
  SourceString string
  OutputString string
  ErrorStatus  int
  ErrorMessage string
}
*/

type FileContext struct {
  Options
  InputPath    string
  OutputString string
  ErrorStatus  int
  ErrorMessage string
}

const (
  NESTED_STYLE = iota
  EXPANDED_STYLE
  COMPACT_STYLE
  COMPRESSED_STYLE
)

/*
func NewSourceContext() (ctx *SourceContext) {
  // something
}

func NewFileContext() (ctx *FileContext) {
  // something
}
*/

func CompileFile(ctx *FileContext) string {
  c_ctx := C.sass_new_file_context()
  c_ctx.input_path = C.CString("some/path/I/guess")
  s := C.GoString(c_ctx.input_path)
  ctx.InputPath = s + "/whatever"
  return ctx.InputPath
}
