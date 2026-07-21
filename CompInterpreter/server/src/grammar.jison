/* ============================================================
   CompInterpreter - Gramatica Jison (lexico + sintactico)
   OLC1 Proyecto 2 - USAC
   Un solo archivo: %lex (tokens) + gramatica que construye el AST.
   Case-insensitive (seccion 5.1) via %options case-insensitive.
   ============================================================ */

/* ---------------------- LEXICO ---------------------- */
%lex
%options case-insensitive flex

%%

/* --- comentarios y espacios (se descartan) --- */
"//".*                              /* comentario de linea */
"/*"([^*]|\*+[^*/])*\*+"/"          /* comentario multilinea */
\s+                                 /* espacios en blanco */

/* --- literales --- */
[0-9]+"."[0-9]+                     return 'DECIMAL';
[0-9]+                              return 'ENTERO';

\"(\\.|[^\\"])*\"    %{
                       yytext = yytext.slice(1, -1).replace(/\\(.)/g, function (m, c) {
                         switch (c) { case 'n': return '\n'; case 't': return '\t';
                           case '"': return '"'; case "'": return "'"; case '\\': return '\\';
                           default: return c; }
                       });
                       return 'CADENA';
                    %}

\'(\\.|[^\\'\n])\'   %{
                       var raw = yytext.slice(1, -1);
                       if (raw.length === 2 && raw[0] === '\\') {
                         switch (raw[1]) { case 'n': raw = '\n'; break; case 't': raw = '\t'; break;
                           case '"': raw = '"'; break; case "'": raw = "'"; break;
                           case '\\': raw = '\\'; break; default: raw = raw[1]; }
                       }
                       yytext = raw;
                       return 'CARACTER';
                    %}

/* --- palabras reservadas (ANTES del identificador) --- */
"true"                              return 'TRUE';
"false"                             return 'FALSE';
"null"                              return 'NULL';
"let"                               return 'LET';
"const"                             return 'CONST';
"int"                               return 'TIPO_INT';
"double"                            return 'TIPO_DOUBLE';
"bool"                              return 'TIPO_BOOL';
"char"                              return 'TIPO_CHAR';
"string"                            return 'TIPO_STRING';
"void"                              return 'VOID';
"new"                               return 'NEW';
"vector"                            return 'VECTOR';
"cast"                              return 'CAST';
"as"                                return 'AS';
"if"                                return 'IF';
"else"                              return 'ELSE';
"switch"                            return 'SWITCH';
"case"                              return 'CASE';
"default"                           return 'DEFAULT';
"while"                             return 'WHILE';
"for"                               return 'FOR';
"do"                                return 'DO';
"until"                             return 'UNTIL';
"loop"                              return 'LOOP';
"break"                             return 'BREAK';
"continue"                          return 'CONTINUE';
"return"                            return 'RETURN';
"function"                          return 'FUNCTION';
"echo"                              return 'ECHO';
"ejecutar"                          return 'EJECUTAR';
"is"                                return 'IS';
"lower"                             return 'LOWER';
"upper"                             return 'UPPER';
"round"                             return 'ROUND';
"len"                               return 'LEN';
"truncate"                          return 'TRUNCATE';
"tostring"                          return 'TOSTRING';
"tochararray"                       return 'TOCHARARRAY';
"reverse"                           return 'REVERSE';
"max"                               return 'MAX';
"min"                               return 'MIN';
"sum"                               return 'SUM';
"average"                           return 'AVERAGE';

/* --- identificador --- */
[a-zA-Z_][a-zA-Z0-9_]*              return 'ID';

/* --- operadores (multi-caracter ANTES de los simples) --- */
"++"                                return 'INCR';
"--"                                return 'DECR';
"=="                                return 'IGUAL';
"!="                                return 'DIFERENTE';
"<="                                return 'MENORIGUAL';
">="                                return 'MAYORIGUAL';
"&&"                                return 'AND';
"||"                                return 'OR';
"+"                                 return 'MAS';
"-"                                 return 'MENOS';
"*"                                 return 'POR';
"/"                                 return 'DIV';
"^"                                 return 'POTENCIA';
"$"                                 return 'RAIZ';
"%"                                 return 'MODULO';
"<"                                 return 'MENOR';
">"                                 return 'MAYOR';
"!"                                 return 'NOT';
"="                                 return 'ASIGNA';
"("                                 return 'PARIZQ';
")"                                 return 'PARDER';
"{"                                 return 'LLAVEIZQ';
"}"                                 return 'LLAVEDER';
"["                                 return 'CORIZQ';
"]"                                 return 'CORDER';
";"                                 return 'PUNTOCOMA';
":"                                 return 'DOSPUNTOS';
","                                 return 'COMA';

<<EOF>>                             return 'EOF';

/* Recuperacion de error lexico: descarta el caracter y sigue (seccion 6.1) */
.                                   %{
                                       if (yy.errores) {
                                         yy.errores.push({ tipo: 'Léxico',
                                           descripcion: 'El carácter "' + yytext + '" no pertenece al lenguaje',
                                           linea: yylloc.first_line, columna: yylloc.first_column + 1 });
                                       }
                                    %}

/lex

/* ---------------------- PRECEDENCIA (seccion 5.10) ----------------------
   Se declara de MENOR a MAYOR importancia. nivel 7 (||) es el mas bajo,
   nivel 0 (negacion unaria) el mas alto. */
%right IF DOSPUNTOS          /* operador ternario: mas bajo que todo */
%left OR                     /* nivel 7 */
%left AND                    /* nivel 6 */
%right NOT                   /* nivel 5 */
%left IGUAL DIFERENTE MENOR MENORIGUAL MAYOR MAYORIGUAL IS   /* nivel 4 */
%left MAS MENOS              /* nivel 3 */
%left POR DIV MODULO         /* nivel 2 */
%nonassoc POTENCIA RAIZ      /* nivel 1 */
%right UMINUS                /* nivel 0 (mayor importancia) */

%start inicio

%% /* ---------------------- SINTACTICO ---------------------- */

inicio
    : lista_global EOF        { $$ = nodo('PROGRAMA', { globales: $1 }); return $$; }
    | EOF                     { $$ = nodo('PROGRAMA', { globales: [] }); return $$; }
    ;

/* A nivel global solo: declaraciones, funciones/metodos y ejecutar (seccion 5.27) */
lista_global
    : lista_global elemento_global   { $1.push($2); $$ = $1; }
    | elemento_global                { $$ = [$1]; }
    ;

elemento_global
    : declaracion PUNTOCOMA          { $$ = $1; }
    | funcion                        { $$ = $1; }
    | ejecutar_stmt                  { $$ = $1; }
    | error PUNTOCOMA                { $$ = nodo('ERROR', {}); }
    ;

funcion
    : FUNCTION tipo ID PARIZQ params PARDER bloque
        { $$ = nodo('FUNCION', { retorno: $2, id: $3, params: $5, cuerpo: $7 }, @3); }
    | FUNCTION VOID ID PARIZQ params PARDER bloque
        { $$ = nodo('METODO', { id: $3, params: $5, cuerpo: $7 }, @3); }
    ;

params
    : lista_params            { $$ = $1; }
    |                         { $$ = []; }
    ;

lista_params
    : lista_params COMA param { $1.push($3); $$ = $1; }
    | param                   { $$ = [$1]; }
    ;

param
    : ID DOSPUNTOS tipo                     { $$ = nodo('PARAM', { id: $1, tipoParam: $3, porDefecto: null }, @1); }
    | ID DOSPUNTOS tipo ASIGNA expresion    { $$ = nodo('PARAM', { id: $1, tipoParam: $3, porDefecto: $5 }, @1); }
    ;

bloque
    : LLAVEIZQ instrucciones LLAVEDER   { $$ = $2; }
    | LLAVEIZQ LLAVEDER                 { $$ = []; }
    ;

instrucciones
    : instrucciones instruccion   { if ($2) $1.push($2); $$ = $1; }
    | instruccion                 { $$ = $1 ? [$1] : []; }
    ;

instruccion
    : declaracion PUNTOCOMA       { $$ = $1; }
    | asignacion PUNTOCOMA        { $$ = $1; }
    | incdec PUNTOCOMA            { $$ = $1; }
    | llamada PUNTOCOMA           { $$ = $1; }
    | nativa_stmt PUNTOCOMA       { $$ = $1; }
    | echo_stmt PUNTOCOMA         { $$ = $1; }
    | if_stmt                     { $$ = $1; }
    | switch_stmt                 { $$ = $1; }
    | while_stmt                  { $$ = $1; }
    | for_stmt                    { $$ = $1; }
    | dountil_stmt                { $$ = $1; }
    | loop_stmt                   { $$ = $1; }
    | BREAK PUNTOCOMA             { $$ = nodo('BREAK', {}, @1); }
    | CONTINUE PUNTOCOMA          { $$ = nodo('CONTINUE', {}, @1); }
    | RETURN PUNTOCOMA            { $$ = nodo('RETURN', { valor: null }, @1); }
    | RETURN expresion PUNTOCOMA  { $$ = nodo('RETURN', { valor: $2 }, @1); }
    | error PUNTOCOMA             { $$ = nodo('ERROR', {}); }
    ;

/* ----- declaraciones (sin punto y coma: lo agregan quien las use) ----- */
/* Las cuatro formas comparten el prefijo `mutabilidad lista_ids : tipo`
   para no generar conflicto shift/reduce. Los vectores usan lista_ids pero
   el interprete valida que sea un solo id (mas de uno => error semantico). */
declaracion
    : mutabilidad lista_ids DOSPUNTOS tipo
        { $$ = nodo('DECLARACION', { mutable: $1, ids: $2, tipoDato: $4, valor: null }, @2); }
    | mutabilidad lista_ids DOSPUNTOS tipo ASIGNA expresion
        { $$ = nodo('DECLARACION', { mutable: $1, ids: $2, tipoDato: $4, valor: $6 }, @2); }
    | mutabilidad lista_ids DOSPUNTOS tipo CORIZQ CORDER ASIGNA vector_init1
        { $$ = nodo('DECLARACION_VECTOR', { mutable: $1, ids: $2, tipoDato: $4, dimension: 1, init: $8 }, @2); }
    | mutabilidad lista_ids DOSPUNTOS tipo CORIZQ CORDER CORIZQ CORDER ASIGNA vector_init2
        { $$ = nodo('DECLARACION_VECTOR', { mutable: $1, ids: $2, tipoDato: $4, dimension: 2, init: $10 }, @2); }
    ;

mutabilidad
    : LET     { $$ = true; }
    | CONST   { $$ = false; }
    ;

lista_ids
    : lista_ids COMA ID   { $1.push($3); $$ = $1; }
    | ID                  { $$ = [$1]; }
    ;

/* La alternativa `expresion` es una extension propia (documentada): permite
   inicializar un vector con el resultado de una expresion que produzca un
   vector (p.ej. reverse(v) o toCharArray(s)). No aparece en 5.15.1 pero es
   necesaria para poder almacenar el resultado de esas nativas, que el propio
   enunciado define como devolviendo un vector. No hay conflicto: expresion
   nunca empieza con NEW ni con CORIZQ. */
vector_init1
    : NEW VECTOR tipo CORIZQ expresion CORDER
        { $$ = nodo('VECTOR_NEW', { tipoBase: $3, tam: $5 }, @1); }
    | CORIZQ lista_expresiones CORDER
        { $$ = nodo('VECTOR_LISTA', { valores: $2 }, @1); }
    | expresion
        { $$ = nodo('VECTOR_EXPR', { expr: $1 }, @1); }
    ;

vector_init2
    : NEW VECTOR tipo CORIZQ expresion CORDER CORIZQ expresion CORDER
        { $$ = nodo('VECTOR_NEW2', { tipoBase: $3, filas: $5, cols: $8 }, @1); }
    | CORIZQ lista_filas CORDER
        { $$ = nodo('VECTOR_LISTA2', { filas: $2 }, @1); }
    | expresion
        { $$ = nodo('VECTOR_EXPR', { expr: $1 }, @1); }
    ;

lista_filas
    : lista_filas COMA CORIZQ lista_expresiones CORDER   { $1.push($4); $$ = $1; }
    | CORIZQ lista_expresiones CORDER                    { $$ = [$2]; }
    ;

lista_expresiones
    : lista_expresiones COMA expresion   { $1.push($3); $$ = $1; }
    | expresion                          { $$ = [$1]; }
    ;

/* ----- asignaciones ----- */
asignacion
    : ID ASIGNA expresion
        { $$ = nodo('ASIGNACION', { id: $1, valor: $3 }, @1); }
    | ID CORIZQ expresion CORDER ASIGNA expresion
        { $$ = nodo('ASIGNACION_VECTOR', { id: $1, idx1: $3, idx2: null, valor: $6 }, @1); }
    | ID CORIZQ expresion CORDER CORIZQ expresion CORDER ASIGNA expresion
        { $$ = nodo('ASIGNACION_VECTOR', { id: $1, idx1: $3, idx2: $6, valor: $9 }, @1); }
    ;

incdec
    : ID INCR   { $$ = nodo('INCREMENTO', { id: $1 }, @1); }
    | ID DECR   { $$ = nodo('DECREMENTO', { id: $1 }, @1); }
    ;

echo_stmt
    : ECHO expresion   { $$ = nodo('ECHO', { valor: $2 }, @1); }
    ;

/* ----- llamada a funcion/metodo (parametros nombrados) ----- */
llamada
    : ID PARIZQ args PARDER   { $$ = nodo('LLAMADA', { id: $1, args: $3 }, @1); }
    ;

args
    : lista_args   { $$ = $1; }
    |              { $$ = []; }
    ;

lista_args
    : lista_args COMA arg   { $1.push($3); $$ = $1; }
    | arg                   { $$ = [$1]; }
    ;

arg
    : ID ASIGNA expresion   { $$ = nodo('ARG', { nombre: $1, valor: $3 }, @1); }
    ;

ejecutar_stmt
    : EJECUTAR ID PARIZQ args PARDER PUNTOCOMA
        { $$ = nodo('EJECUTAR', { id: $2, args: $4 }, @1); }
    ;

/* ----- sentencias de control ----- */
if_stmt
    : IF PARIZQ expresion PARDER bloque
        { $$ = nodo('IF', { cond: $3, entonces: $5, sino: null }, @1); }
    | IF PARIZQ expresion PARDER bloque ELSE bloque
        { $$ = nodo('IF', { cond: $3, entonces: $5, sino: $7 }, @1); }
    | IF PARIZQ expresion PARDER bloque ELSE if_stmt
        { $$ = nodo('IF', { cond: $3, entonces: $5, sino: [$7] }, @1); }
    ;

switch_stmt
    : SWITCH PARIZQ expresion PARDER LLAVEIZQ switch_body LLAVEDER
        { $$ = nodo('SWITCH', { expr: $3, casos: $6.casos, defecto: $6.defecto }, @1); }
    ;

switch_body
    : lista_casos                 { $$ = { casos: $1, defecto: null }; }
    | lista_casos caso_default    { $$ = { casos: $1, defecto: $2 }; }
    | caso_default                { $$ = { casos: [], defecto: $1 }; }
    ;

lista_casos
    : lista_casos caso   { $1.push($2); $$ = $1; }
    | caso               { $$ = [$1]; }
    ;

caso
    : CASE expresion DOSPUNTOS instrucciones   { $$ = nodo('CASE', { expr: $2, cuerpo: $4 }, @1); }
    | CASE expresion DOSPUNTOS                 { $$ = nodo('CASE', { expr: $2, cuerpo: [] }, @1); }
    ;

caso_default
    : DEFAULT DOSPUNTOS instrucciones   { $$ = nodo('DEFAULT', { cuerpo: $3 }, @1); }
    | DEFAULT DOSPUNTOS                 { $$ = nodo('DEFAULT', { cuerpo: [] }, @1); }
    ;

/* ----- ciclos ----- */
while_stmt
    : WHILE PARIZQ expresion PARDER bloque
        { $$ = nodo('WHILE', { cond: $3, cuerpo: $5 }, @1); }
    ;

for_stmt
    : FOR PARIZQ for_init PUNTOCOMA expresion PUNTOCOMA for_update PARDER bloque
        { $$ = nodo('FOR', { init: $3, cond: $5, update: $7, cuerpo: $9 }, @1); }
    ;

for_init
    : declaracion   { $$ = $1; }
    | asignacion    { $$ = $1; }
    ;

for_update
    : incdec        { $$ = $1; }
    | asignacion    { $$ = $1; }
    ;

dountil_stmt
    : DO bloque UNTIL PARIZQ expresion PARDER PUNTOCOMA
        { $$ = nodo('DOUNTIL', { cuerpo: $2, cond: $5 }, @1); }
    ;

loop_stmt
    : LOOP bloque   { $$ = nodo('LOOP', { cuerpo: $2 }, @1); }
    ;

/* ----- tipos ----- */
tipo
    : TIPO_INT      { $$ = 'int'; }
    | TIPO_DOUBLE   { $$ = 'double'; }
    | TIPO_BOOL     { $$ = 'bool'; }
    | TIPO_CHAR     { $$ = 'char'; }
    | TIPO_STRING   { $$ = 'string'; }
    ;

/* ----- expresiones ----- */
expresion
    : expresion MAS expresion         { $$ = nodo('BINARIA', { op: '+', izq: $1, der: $3 }, @2); }
    | expresion MENOS expresion       { $$ = nodo('BINARIA', { op: '-', izq: $1, der: $3 }, @2); }
    | expresion POR expresion         { $$ = nodo('BINARIA', { op: '*', izq: $1, der: $3 }, @2); }
    | expresion DIV expresion         { $$ = nodo('BINARIA', { op: '/', izq: $1, der: $3 }, @2); }
    | expresion POTENCIA expresion    { $$ = nodo('BINARIA', { op: '^', izq: $1, der: $3 }, @2); }
    | expresion RAIZ expresion        { $$ = nodo('BINARIA', { op: '$', izq: $1, der: $3 }, @2); }
    | expresion MODULO expresion      { $$ = nodo('BINARIA', { op: '%', izq: $1, der: $3 }, @2); }
    | MENOS expresion %prec UMINUS    { $$ = nodo('UNARIA', { op: '-', expr: $2 }, @1); }
    | expresion IGUAL expresion       { $$ = nodo('RELACIONAL', { op: '==', izq: $1, der: $3 }, @2); }
    | expresion DIFERENTE expresion   { $$ = nodo('RELACIONAL', { op: '!=', izq: $1, der: $3 }, @2); }
    | expresion MENOR expresion       { $$ = nodo('RELACIONAL', { op: '<', izq: $1, der: $3 }, @2); }
    | expresion MENORIGUAL expresion  { $$ = nodo('RELACIONAL', { op: '<=', izq: $1, der: $3 }, @2); }
    | expresion MAYOR expresion       { $$ = nodo('RELACIONAL', { op: '>', izq: $1, der: $3 }, @2); }
    | expresion MAYORIGUAL expresion  { $$ = nodo('RELACIONAL', { op: '>=', izq: $1, der: $3 }, @2); }
    | expresion AND expresion         { $$ = nodo('LOGICA', { op: '&&', izq: $1, der: $3 }, @2); }
    | expresion OR expresion          { $$ = nodo('LOGICA', { op: '||', izq: $1, der: $3 }, @2); }
    | NOT expresion                   { $$ = nodo('LOGICA', { op: '!', izq: $2, der: null }, @1); }
    | expresion IS tipo               { $$ = nodo('IS', { expr: $1, tipoConsulta: $3 }, @2); }
    | PARIZQ expresion PARDER         { $$ = $2; }
    | CAST PARIZQ expresion AS tipo PARDER   { $$ = nodo('CAST', { expr: $3, tipoDestino: $5 }, @1); }
    | ternario                        { $$ = $1; }
    | llamada                         { $$ = $1; }
    | nativa                          { $$ = $1; }
    | acceso_vector                   { $$ = $1; }
    | ID                              { $$ = nodo('IDENTIFICADOR', { nombre: $1 }, @1); }
    | ENTERO                          { $$ = nodo('LITERAL', { tipoLit: 'int', valor: $1 }, @1); }
    | DECIMAL                         { $$ = nodo('LITERAL', { tipoLit: 'double', valor: $1 }, @1); }
    | CADENA                          { $$ = nodo('LITERAL', { tipoLit: 'string', valor: $1 }, @1); }
    | CARACTER                        { $$ = nodo('LITERAL', { tipoLit: 'char', valor: $1 }, @1); }
    | TRUE                            { $$ = nodo('LITERAL', { tipoLit: 'bool', valor: true }, @1); }
    | FALSE                           { $$ = nodo('LITERAL', { tipoLit: 'bool', valor: false }, @1); }
    | NULL                            { $$ = nodo('LITERAL', { tipoLit: 'null', valor: null }, @1); }
    ;

ternario
    : IF PARIZQ expresion PARDER expresion DOSPUNTOS expresion %prec IF
        { $$ = nodo('TERNARIO', { cond: $3, verdadero: $5, falso: $7 }, @1); }
    ;

acceso_vector
    : ID CORIZQ expresion CORDER
        { $$ = nodo('ACCESO_VECTOR', { id: $1, idx1: $3, idx2: null }, @1); }
    | ID CORIZQ expresion CORDER CORIZQ expresion CORDER
        { $$ = nodo('ACCESO_VECTOR', { id: $1, idx1: $3, idx2: $6 }, @1); }
    ;

/* ----- funciones nativas (usadas como expresion) ----- */
nativa
    : LOWER PARIZQ expresion PARDER        { $$ = nodo('NATIVA', { fn: 'lower', arg: $3 }, @1); }
    | UPPER PARIZQ expresion PARDER        { $$ = nodo('NATIVA', { fn: 'upper', arg: $3 }, @1); }
    | ROUND PARIZQ expresion PARDER        { $$ = nodo('NATIVA', { fn: 'round', arg: $3 }, @1); }
    | LEN PARIZQ expresion PARDER          { $$ = nodo('NATIVA', { fn: 'len', arg: $3 }, @1); }
    | TRUNCATE PARIZQ expresion PARDER     { $$ = nodo('NATIVA', { fn: 'truncate', arg: $3 }, @1); }
    | TOSTRING PARIZQ expresion PARDER     { $$ = nodo('NATIVA', { fn: 'tostring', arg: $3 }, @1); }
    | TOCHARARRAY PARIZQ expresion PARDER  { $$ = nodo('NATIVA', { fn: 'tochararray', arg: $3 }, @1); }
    | REVERSE PARIZQ expresion PARDER      { $$ = nodo('NATIVA', { fn: 'reverse', arg: $3 }, @1); }
    | MAX PARIZQ expresion PARDER          { $$ = nodo('NATIVA', { fn: 'max', arg: $3 }, @1); }
    | MIN PARIZQ expresion PARDER          { $$ = nodo('NATIVA', { fn: 'min', arg: $3 }, @1); }
    | SUM PARIZQ expresion PARDER          { $$ = nodo('NATIVA', { fn: 'sum', arg: $3 }, @1); }
    | AVERAGE PARIZQ expresion PARDER      { $$ = nodo('NATIVA', { fn: 'average', arg: $3 }, @1); }
    ;

/* nativa como instruccion suelta (ej. reverse(v);) */
nativa_stmt
    : nativa   { $$ = $1; }
    ;

%%

/* ---------------------- CODIGO DE SOPORTE ---------------------- */
/* Fabrica de nodos del AST. loc = objeto @N de jison (first_line, first_column). */
function nodo(tipo, props, loc) {
  var n = { tipo: tipo };
  if (props) { for (var k in props) { if (props.hasOwnProperty(k)) n[k] = props[k]; } }
  if (loc) { n.linea = loc.first_line; n.columna = loc.first_column + 1; }
  return n;
}
