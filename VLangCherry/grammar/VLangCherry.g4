grammar VLangCherry;

// ============================================================
// VLangCherry (OLC2 - Proyecto 1). Gramática propia, con fines
// descriptivos el enunciado permite diseñarla libremente (11.9).
// Decisiones de desambiguación tomadas frente a inconsistencias
// del enunciado (documentadas en docs/gramatica.txt):
//   - palabra reservada para declarar funciones: "func" (no "fn"):
//     "fn" solo aparece una vez en un ejemplo, "func" es consistente
//     en el resto del documento y en la seccion 7 completa.
//   - literales rune con comilla simple 'x' (3.7), no backtick.
// ============================================================

programa
    : declaracionGlobal* EOF
    ;

declaracionGlobal
    : declaracionStruct
    | declaracionFuncion
    | declaracionVariable
    ;

// ---------------- Structs ----------------

declaracionStruct
    : 'struct' ID '{' campoStruct+ '}'
    ;

campoStruct
    : tipo ID
    ;

// ---------------- Funciones y metodos ----------------

declaracionFuncion
    : 'func' receptor? ID '(' listaParametros? ')' tipo? bloque
    ;

receptor
    : '(' ID '*'? ID ')'
    ;

listaParametros
    : parametro (',' parametro)*
    ;

parametro
    : ID tipo
    ;

// ---------------- Tipos ----------------

tipo
    : tipoSlice
    | tipoPrimitivo
    | ID
    ;

tipoPrimitivo
    : 'int'
    | 'float64'
    | 'string'
    | 'bool'
    | 'rune'
    ;

tipoSlice
    : '[' ']' tipo
    ;

// ---------------- Sentencias ----------------

bloque
    : '{' sentencia* '}'
    ;

sentencia
    : declaracionVariable    #sentDeclaracion
    | asignacion             #sentAsignacion
    | incrementoDecremento   #sentIncDec
    | expr ';'?              #sentExpresion
    | sentenciaIf            #sentIf
    | sentenciaSwitch        #sentSwitch
    | sentenciaFor           #sentFor
    | 'break' ';'?           #sentBreak
    | 'continue' ';'?        #sentContinue
    | 'return' expr? ';'?    #sentReturn
    | bloque                 #sentBloque
    ;

declaracionVariable
    : 'mut'? ID tipo ('=' expr)? ';'?     # declTipada
    | 'mut'? ID ':=' expr ';'?            # declInferida
    ;

asignacion
    : lugar op=('='|'+='|'-=') expr ';'?
    ;

lugar
    : ID                        # lugarId
    | lugar '[' expr ']'        # lugarIndice
    | lugar '.' ID              # lugarCampo
    ;

incrementoDecremento
    : lugar op=('++'|'--') ';'?
    ;

sentenciaIf
    : 'if' expr bloque ('else' 'if' expr bloque)* ('else' bloque)?
    ;

sentenciaSwitch
    : 'switch' expr '{' casoSwitch* defaultSwitch? '}'
    ;

casoSwitch
    : 'case' expr ':' sentencia*
    ;

defaultSwitch
    : 'default' ':' sentencia*
    ;

sentenciaFor
    : 'for' expr bloque                                        # forCondicion
    | 'for' forInit? ';' expr? ';' forActualizacion? bloque     # forClasico
    | 'for' ID ',' ID 'in' expr bloque                          # forRango
    ;

forInit
    : declaracionVariable
    | asignacion
    ;

forActualizacion
    : asignacion
    | incrementoDecremento
    ;

// ---------------- Expresiones ----------------
// Orden de alternativas = precedencia (mayor a menor), tabla 4.6:
//   ( ) [ ] . llamada   >  ! -(unario)  >  * / %  >  + -
//   >  < <= >= >  >  == !=  >  &&  >  ||

expr
    : expr '(' listaArgumentos? ')'            # exprLlamada
    | expr '[' expr ']'                        # exprIndice
    | expr '.' ID                              # exprCampo
    | '(' expr ')'                             # exprParentesis
    | op=('!'|'-') expr                        # exprUnario
    | expr op=('*'|'/'|'%') expr               # exprMultiplicativa
    | expr op=('+'|'-') expr                   # exprAditiva
    | expr op=('<'|'<='|'>='|'>') expr         # exprRelacional
    | expr op=('=='|'!=') expr                 # exprIgualdad
    | expr '&&' expr                           # exprAnd
    | expr '||' expr                           # exprOr
    | literalSlice                             # exprSliceLit
    | literalStruct                            # exprStructLit
    | literal                                  # exprLiteral
    | ID                                       # exprIdentificador
    ;

listaArgumentos
    : expr (',' expr)*
    ;

literalSlice
    : '[' ']' tipo '{' listaArgumentos? '}'
    | '[' ']' tipo '{' filaSlice (',' filaSlice)* ','? '}'
    ;

filaSlice
    : '{' listaArgumentos? '}'
    ;

literalStruct
    : ID '{' campoValor (',' campoValor)* ','? '}'
    ;

campoValor
    : ID ':' expr
    ;

literal
    : ENTERO
    | DECIMAL
    | CADENA
    | RUNE
    | 'true'
    | 'false'
    | 'nil'
    ;

// ============================================================
// LEXER
// ============================================================

MUT: 'mut';
STRUCT: 'struct';
FUNC: 'func';
IF: 'if';
ELSE: 'else';
SWITCH: 'switch';
CASE: 'case';
DEFAULT: 'default';
FOR: 'for';
IN: 'in';
BREAK: 'break';
CONTINUE: 'continue';
RETURN: 'return';
NIL: 'nil';
TRUE: 'true';
FALSE: 'false';
TIPO_INT: 'int';
TIPO_FLOAT: 'float64';
TIPO_STRING: 'string';
TIPO_BOOL: 'bool';
TIPO_RUNE: 'rune';

MASIGUAL: '+=';
MENOSIGUAL: '-=';
INCREMENTO: '++';
DECREMENTO: '--';
ASIGNAINFERIDA: ':=';
IGUAL: '==';
DIFERENTE: '!=';
MENORIGUAL: '<=';
MAYORIGUAL: '>=';
Y: '&&';
O: '||';
MENOR: '<';
MAYOR: '>';
MAS: '+';
MENOS: '-';
POR: '*';
DIV: '/';
MODULO: '%';
NOT: '!';
ASIGNA: '=';
PARIZQ: '(';
PARDER: ')';
CORIZQ: '[';
CORDER: ']';
LLAIZQ: '{';
LLADER: '}';
COMA: ',';
PUNTO: '.';
DOSPUNTOS: ':';
PUNTOCOMA: ';';

ID
    : [a-zA-Z_][a-zA-Z0-9_]*
    ;

ENTERO
    : [0-9]+
    ;

DECIMAL
    : [0-9]+ '.' [0-9]+
    ;

CADENA
    : '"' (ESC | ~["\\\r\n])* '"'
    ;

RUNE
    : '\'' (ESC | ~['\\\r\n]) '\''
    ;

fragment ESC
    : '\\' ["\\ntr]
    ;

COMENTARIO_LINEA
    : '//' ~[\r\n]* -> skip
    ;

COMENTARIO_BLOQUE
    : '/*' .*? '*/' -> skip
    ;

WS
    : [ \t\r\n]+ -> skip
    ;
