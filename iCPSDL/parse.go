package iCPSDL

import "iCPSDL/parser"
import "github.com/antlr/antlr4/runtime/Go/antlr"

import "iCPSDL/util"

/******************************************************************************
 * Custom syntax error Listener
 ******************************************************************************/

type customErrorListener struct {
    *antlr.DefaultErrorListener
    util.BasicErrorLog
    file string
}

func newCustomErrorListener(file string) (*customErrorListener) {
    this := new(customErrorListener)
    this.file = file
    return this
}

func (this *customErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
    customError := util.NewSyntaxError(msg, this.file, line, column)
    this.Add(customError)
}

/******************************************************************************
 * Parser
 ******************************************************************************/

func Parse(file string) ([]expression, util.ErrorLog) {
    input, err := antlr.NewFileStream(file)

    if err != nil {
        return nil, util.NewErrorLog().Add(util.NewSystemError(err.Error()))
    }

    lexerLog := newCustomErrorListener(file)
    lexer := parser.NewontologyLexer(input)
    lexer.RemoveErrorListeners()
    lexer.AddErrorListener(lexerLog)

    stream := antlr.NewCommonTokenStream(lexer, 0)

    parserLog := newCustomErrorListener(file)
    p := parser.NewontologyParser(stream)
    p.RemoveErrorListeners()
    p.AddErrorListener(parserLog)
    tree := p.Program()

    if lexerLog.HasErrors() {
        return nil, lexerLog
    }

    if parserLog.HasErrors() {
        return nil, parserLog
    }

    return new(ontologyVisitor).VisitProgram(tree.(*parser.ProgramContext)), nil
}

func parse(input antlr.CharStream, log *customErrorListener) []expression {
    lexer := parser.NewontologyLexer(input)
    lexer.RemoveErrorListeners()
    lexer.AddErrorListener(log)

    stream := antlr.NewCommonTokenStream(lexer, 0)

    p := parser.NewontologyParser(stream)
    p.RemoveErrorListeners()
    p.AddErrorListener(log)
    tree := p.Program()

    if log.HasErrors() {
        return nil
    }

    return new(ontologyVisitor).VisitProgram(tree.(*parser.ProgramContext))
}

func parseFile(file string, log *customErrorListener) []expression {
    input, err := antlr.NewFileStream(file)

    if err != nil {
        log.Add(util.NewSystemError(err.Error()))
        return nil
    }
    return parse(input, log)
}


func parseString(inputString string, log *customErrorListener) []expression {
    input:= antlr.NewInputStream(inputString)
    return parse(input, log)
}
