package iCPSDL

import "iCPSDL/util"

type baseNode struct {
    filename string
    line int
}

func (this *baseNode) init(line int) {
    this.line = line
    this.filename = ""
}

func (this *baseNode) getLine() int {
    return this.line
}

func (this *baseNode) setFilename(filename string) {
    this.filename = filename
}

func (this *baseNode) getFilename() string {
    return this.filename
}

func (this *baseNode) reportError(err string, log util.ErrorLog) {
    log.Add(util.NewTypeError(err, this.filename, this.line))
}
