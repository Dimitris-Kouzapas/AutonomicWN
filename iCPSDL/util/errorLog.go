package util

import "fmt"
import "strings"

type ErrorLog interface {
    Add(e error) ErrorLog
    HasErrors() bool
    Iterator() func() (error, bool)
    PrintErrors()
    ToString() string
}

/*******************************************************************************
 *  ErrorLog
 ******************************************************************************/

type BasicErrorLog struct {
    log []error
}

func NewErrorLog() (this *BasicErrorLog) {
    this = new(BasicErrorLog)
    return
}

func (this *BasicErrorLog) Add(e error) ErrorLog {
    this.log = append(this.log, e)
    return this
}

func (this *BasicErrorLog) HasErrors() bool {
    if this.log != nil {
        return len(this.log) != 0
    }
    return false
}

func (this *BasicErrorLog) Iterator() func() (error, bool) {
    var size, index int = 0, 0
    if this.log != nil {
        size = len(this.log)
    }
    return  func() (error, bool) {
                if index < size {
                    elem := this.log[index]
                    index++
                    return elem, true
                }
                return nil, false
            }
}

func (this *BasicErrorLog) PrintErrors() {
    for _, e := range this.log {
        fmt.Println(e.Error())
    }
}


func (this *BasicErrorLog) ToString() string {
    // Use a Builder for efficient string concatenation
    var builder strings.Builder
    for _, e := range this.log {
        builder.WriteString(fmt.Sprintf("%v\n", e.Error()))
    }
    return builder.String()
}

/*******************************************************************************
 * basic Error
 ******************************************************************************/

type basicError struct {
    err string
    errorType string
    file string
    line int
}

func (this *basicError) init(err string, errorType string, file string, line int) {
    this.err = err
    this.errorType = errorType
    this.file = file
    this.line = line
    return
}

func (this *basicError) Error() string {
    return fmt.Sprintf("   %v:%v: %v error!\n\t%v", this.file, this.line, this.errorType, this.err)
}

/*******************************************************************************
 * Custom Syntax Error
 ******************************************************************************/

type syntaxError struct {
    basicError
    column int
}

func NewSyntaxError(msg string, file string, line int, column int) (this *syntaxError) {
    this = new(syntaxError)
    this.init(msg, "Syntax", file, line)
    this.column = column
    return
}

func (this *syntaxError) Error() string {
    return fmt.Sprintf("   %v:%v:%v Syntax error!\n\t%v", this.file, this.line, this.column, this.err)
}

/*******************************************************************************
 * type check Error
 ******************************************************************************/

type typeError struct {
    basicError
}

func NewTypeError(err string, file string, line int) (this *typeError) {
    this = new(typeError)
    this.init(err, "Type check", file, line)
    return
}

/*******************************************************************************
 * system Error
 ******************************************************************************/

type systemError struct {
    basicError
}

func NewSystemError(err string) (this *systemError) {
    this = new(systemError)
    this.init(err, "System", "", 0)
    return
}

func (this *systemError) Error() string {
    return fmt.Sprintf("   %v error!\n\t%v", this.errorType, this.err)
}


/*******************************************************************************
 * Command Line error
 ******************************************************************************/

type commandLineError struct {
    basicError
}

func NewCommandLineError(msg string, file string, line int) (this *commandLineError) {
    this = new(commandLineError)
    this.init(msg, "Command line", file, line)
    return
}

/*******************************************************************************
 * Execution error
 ******************************************************************************/

type executionException struct {
    basicError
}

func NewExecutionException(msg string) (this *executionException) {
  this = new(executionException)
  this.init(msg, "Execution exception", "", 0)
  return
}

func (this *executionException) Error() string {
  return fmt.Sprintf("   %v error!\n\t%v", this.errorType, this.err)
}
