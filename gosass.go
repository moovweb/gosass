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

func CompileFile(goCtx *FileContext) {
  // set up the underlying C context struct
  cCtx := C.sass_new_file_context()
  cCtx.input_path           = C.CString(goCtx.InputPath)
  cCtx.options.output_style = C.int(goCtx.Options.OutputStyle)
  // call the underlying C compile function to populate the C context
  C.sass_compile_file(cCtx)
  // extract values from the C context to populate the Go context object
  goCtx.OutputString = C.GoString(cCtx.output_string)
  goCtx.ErrorStatus  = int(cCtx.error_status)
  goCtx.ErrorMessage = C.GoString(cCtx.error_message)
}
