package util

type Stream struct {
  indent int
  indentSymbol string
  text string
  changeLine bool
}

func NewStream() (this *Stream) {
  this = new(Stream)
  this.indent = 0
  this.indentSymbol = "\t"
  this.text = ""
  this.changeLine = true
  return
}

func (this *Stream) Set(symbol string) {
  this.indentSymbol = symbol
}

func (this *Stream) Inc() {
  this.indent++
}

func (this *Stream) Dec() {
  if this.indent > 0 {
    this.indent--
  }
}

func (this *Stream) indentation() (indent string) {
  indent = "   "
  for i := 0; i < this.indent; i++ {
    indent += this.indentSymbol
  }
  return
}

//func (this *Stream) PrintlnNoIndent(s string) {
//  this.text += s
//}

func (this *Stream) Print(s string) {
  if this.changeLine == true {
    this.text += this.indentation()
  }
  this.text += s
  this.changeLine = false
}

func (this *Stream) Println(s string) {
  this.Print(s)
  this.text += "\n"
  this.changeLine = true
}

func (this *Stream) ToString() string {
  return this.text
}
