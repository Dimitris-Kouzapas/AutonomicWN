package ast
import "fmt"

// monotonic counter
var counter int = 0

func newCall(exprs []expression, adef expression, variables []*variableExpr, proc process, line int) *introProc {
    name  := fmt.Sprintf("*(%v@line:%v)", adef.String(), line)
    counter += 1
    pexpr := newParticipantExpr(name, line)

    var lexpr expression = pexpr
    for _, rexpr := range exprs {
        lexpr = newSendExpr(lexpr, rexpr, line)
    }

    var expr seqExpr = nil
    for _, variable := range variables {
        expr = newReceiveExpr(variable, lexpr, line)
        lexpr = expr
    }

    proc = newSequentialProc(expr, proc, line)
    participant := newParticipantExpr("user", line)

    appl := newApplication(adef, []*participantExpr{ participant }, line)
    return newIntroProc(participant, []*participantExpr{pexpr}, []*application{appl}, proc, line)
}
