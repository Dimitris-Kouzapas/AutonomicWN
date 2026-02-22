grammar sessions;

module          :   imports declaration* EOF
                ;

imports         :   (IMPORT path SEMICOLON?)*
                ;

path            :   ID (DOT ID)*
                ;

declaration     :   VAL name AS abstraction
                |   VAL name AS primaryTerm
                |   TYPE name session
                |   TYPE name configuration
                |   TYPE name recordType
                |   TYPE name brokerType
                |   name COLON sessionDef
                |   declarationSugar
                ;

declarationSugar:   PROC name participant* DOT process
                |   RECORD name primaryRecord
                |   BROKER name broker
                ;

name            :   ID
                ;

participant     :   ID
                ;

variableDef     :   variable type
                ;

variable        :   ID
                ;

application     :   term participant*
                ;

abstraction     :   PROC participant* DOT process
                ;

concurrent      :	CONC WITH (ROLE name AS CURRENT SEMICOLON?)? (ROLE participant AS application SEMICOLON?)+
				|	CONC WITH LBRA (ROLE name AS CURRENT SEMICOLON?)? (ROLE participant AS application SEMICOLON?)+ RBRA
                ;

acceptProc      :   ACC LPAR variable RPAR WITH (ROLE participant AS application SEMICOLON?)+
				|	ACC LPAR variable RPAR WITH LBRA (ROLE participant AS application SEMICOLON?)+ RBRA
                ;

request         :	REQ LPAR variable RPAR WITH ROLE name AS REQUESTER SEMICOLON? (ROLE participant AS application SEMICOLON?)*
				|	REQ LPAR variable RPAR WITH LBRA ROLE name AS REQUESTER SEMICOLON? (ROLE participant AS application SEMICOLON?)* RBRA
                ;

recurse         :	REC LPAR variable RPAR WITH RECURSION ROLE participant AS application SEMICOLON? (ROLE participant AS application SEMICOLON?)*
				|	REC LPAR variable RPAR WITH LBRA RECURSION ROLE participant AS application SEMICOLON? (ROLE participant AS application SEMICOLON?)* RBRA
                ;

process         :   blockProc
                |   sequenceProc
                |   ifThenElseProc
                |   selectProc
                |   branchProc
                |   acceptProc
                |   recurse
                |   terminateProc
                ;

blockProc       :   LBRA process RBRA
                ;

sequenceProc    :   prefix SEMICOLON? process?
                |   concurrent SEMICOLON? process?
                |   request SEMICOLON? process?
                |   call SEMICOLON? process?
                ;

prefix          :   receive
                |   out
                |   inp
                |   let
                |   closeBroker
                ;

receive         :   variableDef ARROW participant
                |   variableDef ARROW send
                |   send
                ;

send            :   participant (ARROW logicalOr)+
                |   LPAR receive RPAR (ARROW logicalOr)*
                ;

out             :   OUT (LSQ term RSQ)? (ARROW logicalOr)+
                ;

inp             :   variableDef ARROW INP (LSQ term RSQ)?
                ;

let             :   LET variable type? AS logicalOr
                ;

closeBroker     :   variable ARROW CLOSE
                ;

call            :   variableDef (COMMA variableDef)* EQUALS term LPAR (logicalOr (COMMA logicalOr)*)? RPAR
                ;

ifThenElseProc  :   IF logicalOr THEN process ELSE process
                ;

selectProc      :   SELECT participant OF (logicalOr ARROW labelExpression COLON process)* labelExpression COLON process
				|	SELECT participant OF LBRA (logicalOr ARROW labelExpression COLON process)* labelExpression COLON process RBRA
                |   participant ARROW labelExpression process?
                ;

branchProc      :   BRANCH participant OF (labelExpression COLON process)* labelExpression COLON process
				|	BRANCH participant OF LBRA (labelExpression COLON process)* labelExpression COLON process RBRA
                ;

labelExpression :   LABEL
                ;

terminateProc   :   TERMINATE SEMICOLON?
                ;

logicalOr       :   logicalOr LOR logicalAnd
                |   logicalAnd
                ;

logicalAnd      :   logicalAnd LAND equalityExpr
                |   equalityExpr
                ;

equalityExpr    :   equalityExpr EQOP relationalExpr
                |   relationalExpr
                ;

relationalExpr  :   relationalExpr RELOP sumExpr
                |   sumExpr
                ;

sumExpr         :   sumExpr SUMOP multExpr
                |   multExpr
                ;

multExpr        :   multExpr MULOP unary
                |   unary
                ;

unary           :   NOT unary
                |   SUMOP unary
                |   term
                ;

term            :   literal
                |   variable                                      // variable
                |   listExpr                                      // list
                |   term listSlice
                |   term LSQ logicalOr RSQ                        // list access
                |   conditionalExpr
                |	abstraction
                |   record
                |   broker
                |   term DOT ID
                |   LPAR logicalOr RPAR
                ;

literal         :   INT_LIT
                |   FLOAT_LIT
                |   TRUE | FALSE
                |   STRING_LIT
                ;

listExpr        :   LSQ (logicalOr (COMMA logicalOr)*)? COMMA? RSQ
                ;

listSlice       :   LSQ left=logicalOr? COLON right=logicalOr? RSQ
                ;

conditionalExpr :   IF logicalOr THEN logicalOr ELSE logicalOr
                ;

record          :   LBRA (ID COLON logicalOr (SEMICOLON | COMMA)?)* RBRA
                ;

broker          :   SYNC LBRA CHANNEL COLON configurationDef (SEMICOLON | COMMA)? (REQUESTER COLON participant (SEMICOLON | COMMA)?)? RBRA 
                |   ASYNC LBRA CHANNEL COLON configurationDef RBRA
                ;

primaryTerm     :   literal
                |   primaryList
                |   primaryRecord
                |   broker
                ;

primaryList     :   LSQ (primaryTerm (COMMA primaryTerm)*)? COMMA? RSQ
                ;

primaryRecord   :   LBRA (ID COLON primaryTerm (SEMICOLON | COMMA)?)* RBRA
                ;

////////////////////////////////////////////////////////////////////////////////

type            :   primitiveType
                |   LSQ type RSQ
                |   recordType
                |   brokerType
                |   session
                |   nameType
                |   ioType
                |   LPAR type RPAR
                ;

primitiveType   :   BOOLEAN
                |   INT
                |   FLOAT
                |   STRING
                ;

nameType        :   name
                ;

ioType          :   IO
                ;

session         :   localAbstraction
                |   projection
                ;

configuration   :   globalDef
                |   localContext
                ;

recordType      :   RECORD LBRA (ID type (SEMICOLON | COMMA)?)* RBRA
                ;

brokerType      :   BROKER LBRA configurationDef CHANNEL (SEMICOLON | COMMA)? (ID REQUESTER (SEMICOLON | COMMA)?)? RBRA 
                ;

sessionDef      :   session
                |   nameType
                ;

configurationDef:   configuration
                |   nameType
                ;

/*********************
 *   local types
 *********************/


projection      :   configurationDef LPAR ID RPAR
                ;

localAbstraction:   LOCAL ID* DOT local
                ;

labelType       :   LABEL
                ;

local           :   sendLocal
                |   receiveLocal
                |   selectLocal
                |   branchLocal
                |   end
                ;

sendLocal       :   ID SEND type (DOT local)?
                ;

receiveLocal    :   ID RECEIVE type (DOT local)?
                ;

selectLocal     :   SELECT ID LBRA labelLocal RBRA (OR LBRA labelLocal RBRA)*
                ;

branchLocal     :   BRANCH ID LBRA labelLocal RBRA (OR LBRA labelLocal RBRA)*
                ;

labelLocal      :   labelType COLON local
                ;

end             :   END
                ;

/*********************
 *   local context
 *********************/

localContext    :   CONTEXT LBRA (ID COLON local COMMA?)+ RBRA
                ;

/*********************
 *   global types
 *********************/

globalDef       :   GLOBAL ID* DOT global
                ;

global          :   pass
                |   choice
                |   end
                ;

pass            :   ID PASS ID COLON type (DOT global)?
                ;

choice          :   ID PASS ID LBRA labelGlobal RBRA (OR LBRA labelGlobal RBRA)*
                ;

labelGlobal     :   labelType COLON global
                ;

/******************************************************************************
 *  Lexer
 ******************************************************************************/

IMPORT:         'import';
PROC:           'proc';
VAL:            'val';
TYPE:           'type';
LET:			'let';
AS:				'as';
CONC:           'conc';
REQ:            'req';
ACC:            'acc';
REC:            'rec';
CURRENT:		'current';
RECURSION:		'recursion';
ROLE:			'role';
WITH:			'with';
SELECT:         'select';
BRANCH:         'branch';
OF:				'of';
IF:             'if';
THEN:           'then';
ELSE:           'else';
TERMINATE:      'term';

IO:             'io';

BROKER:         'broker';
CHANNEL:        'chan';
REQUESTER:      'requester';
SYNC:			'sync';
ASYNC:			'async';

COLON:          ':';
SEMICOLON:      ';';
DOT:            '.';
COMMA:          ',';
ARROW:          '<-';

LPAR:           '(';
RPAR:           ')';
LSQ:            '[';
RSQ:            ']';
LBRA:           '{';
RBRA:           '}';

LOR:            '||';
LAND:           '&&';
EQOP:           ('=='|'!=');
RELOP:          ('<'|'>'|'<='|'>=');
SUMOP:          ('+'|'-');
MULOP:          ('*'|'/'|'%'|'^');
NOT:            'not';

EQUALS:         '=';

OUT:            'out';
INP:            'inp';
CLOSE:          '.close';

fragment DIGITS:
                [0-9]+ ;
fragment EXP:   [eE] [+\-]? DIGITS ;

FLOAT_LIT:      DIGITS '.' DIGITS (EXP)?
         |      DIGITS '.' (EXP)?
         |      '.' DIGITS (EXP)?
         |      DIGITS EXP
         ;

INT_LIT:        DIGITS ;

TRUE:           'true';
FALSE:          'false';
//STRING_LIT:     '"'~('"')*'"';
STRING_LIT : '"' ( '\\' . | ~["\\] )* '"' ;

// Types

INT:            'int';
FLOAT:          'float';
STRING:         'string';
BOOLEAN:        'bool';

RECORD:         'record';

LOCAL:          'local';
SEND:           '!';
RECEIVE:        '?';
OR:             'or';
END:            'end';

GLOBAL:         'global';
PASS:           '->';
CONTEXT:        'context';

// Identifiers

ID:             [a-z][a-zA-Z0-9_]*;
LABEL:          [A-Z][a-zA-Z0-9_]*;

COMMENTS:       '#' ~( '\n'|'\r' )* '\r'? '\n'  -> skip;
BLOCK_COMMENT:  '/*' ( . | '\r' | '\n' )*? '*/' -> skip ;
//BLOCK_COMMENT:  '/*' .*? '*/' -> skip ;
WS:             [ \t\n\r] -> skip;
