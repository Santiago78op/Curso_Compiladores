/* ============================================================
 * Lexer.flex — Analizador léxico de DataForge (OLC1, Proyecto 1)
 *
 * SECCIÓN 1: código de usuario (se copia tal cual arriba de la clase)
 * ============================================================ */
package dataforge.analisis;

import java_cup.runtime.Symbol;

%%

/* ============================================================
 * SECCIÓN 2: opciones y definiciones regulares (macros)
 * ============================================================ */
%class Lexer
%public
%unicode
%cup
%line
%column
%ignorecase

%{
  /* Conexión con el entorno (Etapa 6): el Interprete la setea tras
     crear el lexer. Permite registrar cada token para el reporte §6.1
     y mandar los errores léxicos a la tabla §6.2 en vez de a stderr. */
  public dataforge.interprete.Entorno entorno;

  /* Helpers: construyen el Symbol que consumirá el parser CUP.
     Guardamos yyline/yycolumn para los reportes (empiezan en 0). */
  private Symbol symbol(int type) {
    registrar(type);
    return new Symbol(type, yyline, yycolumn, yytext());
  }
  private Symbol symbol(int type, Object value) {
    registrar(type);
    return new Symbol(type, yyline, yycolumn, value);
  }
  private void registrar(int type) {
    if (entorno != null) entorno.registrarToken(yytext(), type, yyline, yycolumn);
  }
%}

/* --- definiciones regulares: la columna "Patrón" de nuestra tabla --- */
Letra        = [a-zA-Z]
Digito       = [0-9]
Id           = {Letra}({Letra}|{Digito})*
IdArreglo    = "@"{Id}
Numero       = {Digito}+("."{Digito}+)?
Cadena       = \"[^\"\r\n]*\"
ComentLinea  = "!"[^\r\n]*
ComentMulti  = "<!"~"!>"
Blancos      = [ \t\r\n]+

%%

/* ============================================================
 * SECCIÓN 3: reglas — "cuando veas esto → entregá este token"
 * Orden importante: reservadas ANTES que {Id} (a igual longitud
 * de match, gana la regla escrita primero).
 * ============================================================ */

/* ---- 1. Se reconoce pero se DESCARTA (no produce token) ---- */
{ComentLinea}   { /* comentario de línea: ignorar */ }
{ComentMulti}   { /* comentario multilínea: ignorar */ }
{Blancos}       { /* espacios, tabs y saltos: solo separan */ }

/* ---- 2. Palabras reservadas (con %ignorecase son case insensitive) ---- */
"program"                 { return symbol(sym.PROGRAM); }
"end"                     { return symbol(sym.END); }
"var"                     { return symbol(sym.VAR); }
"arr"                     { return symbol(sym.ARR); }
"double"                  { return symbol(sym.DOUBLE); }
"char"                    { return symbol(sym.CHAR); }

"sum"                     { return symbol(sym.SUM); }
"res"                     { return symbol(sym.RES); }
"mul"                     { return symbol(sym.MUL); }
"div"                     { return symbol(sym.DIV); }
"mod"                     { return symbol(sym.MOD); }

"media"                   { return symbol(sym.MEDIA); }
"mediana"                 { return symbol(sym.MEDIANA); }
"moda"                    { return symbol(sym.MODA); }
"varianza"                { return symbol(sym.VARIANZA); }
"max"                     { return symbol(sym.MAX); }
"min"                     { return symbol(sym.MIN); }

"console"                 { return symbol(sym.CONSOLE); }
"print"                   { return symbol(sym.PRINT); }
"column"                  { return symbol(sym.COLUMN); }

/* el enunciado abre con graphX( pero ejecuta EXEC grapX (typo sistemático):
   aceptamos ambas formas para el mismo token */
"graphbar"  | "grapbar"   { return symbol(sym.GRAPH_BAR); }
"graphpie"  | "grappie"   { return symbol(sym.GRAPH_PIE); }
"graphline" | "grapline"  { return symbol(sym.GRAPH_LINE); }
"histogram"               { return symbol(sym.HISTOGRAM); }
"exec"                    { return symbol(sym.EXEC); }

/* ---- 3. Símbolos (el longest match resuelve :: vs : y <- solo) ---- */
"::"                      { return symbol(sym.DOBLE_DOS_PUNTOS); }
":"                       { return symbol(sym.DOS_PUNTOS); }
"<-"                      { return symbol(sym.ASIGNACION); }
"->"                      { return symbol(sym.FLECHA); }
"="                       { return symbol(sym.IGUAL); }
";"                       { return symbol(sym.PUNTO_COMA); }
","                       { return symbol(sym.COMA); }
"("                       { return symbol(sym.PAR_IZQ); }
")"                       { return symbol(sym.PAR_DER); }
"["                       { return symbol(sym.COR_IZQ); }
"]"                       { return symbol(sym.COR_DER); }

/* ---- 4. Tokens con patrón (después de las reservadas) ---- */
{Numero}                  { return symbol(sym.NUMERO, Double.valueOf(yytext())); }
{Cadena}                  { /* valor sin las comillas */
                            return symbol(sym.CADENA, yytext().substring(1, yylength()-1)); }
{IdArreglo}               { return symbol(sym.ID_ARREGLO); }
{Id}                      { return symbol(sym.ID); }

/* ---- 5. Cualquier otro carácter = error léxico ----
   Va a la tabla de errores del entorno (reporte §6.2). El análisis
   CONTINÚA: el carácter intruso se descarta y se sigue escaneando. */
[^]                       { if (entorno != null) {
                              entorno.error("Léxico", "el carácter '" + yytext()
                                  + "' no pertenece al lenguaje", yyline, yycolumn);
                            } else {
                              System.err.printf(
                                  "Error léxico: '%s' no pertenece al lenguaje (línea %d, columna %d)%n",
                                  yytext(), yyline + 1, yycolumn + 1);
                            } }
