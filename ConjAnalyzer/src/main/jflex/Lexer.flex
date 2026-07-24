/* ============================================================
 * Lexer.flex — Analizador lexico de ConjAnalyzer (OLC1, Proyecto 1)
 *
 * SECCION 1: codigo de usuario (se copia tal cual arriba de la clase)
 * ============================================================ */
package conjanalyzer.analisis;

import java_cup.runtime.Symbol;

%%

/* ============================================================
 * SECCION 2: opciones y definiciones regulares (macros)
 *   OJO: NO se usa %ignorecase — el lenguaje es CASE SENSITIVE (4.1)
 * ============================================================ */
%class Lexer
%public
%unicode
%cup
%line
%column

%{
  /* Conexion con el entorno: el Interprete la setea tras crear el lexer.
     Permite registrar cada token para el reporte de tokens (5.2) y mandar
     los errores lexicos a la tabla de errores (5.3) en vez de a stderr. */
  public conjanalyzer.interprete.Entorno entorno;

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

/* --- definiciones regulares --- */
Letra        = [a-zA-Z]
Digito       = [0-9]
Id           = {Letra}({Letra}|{Digito})*
Numero       = {Digito}+
ComentLinea  = "#"[^\r\n]*
ComentMulti  = "<!"~"!>"
Blancos      = [ \t\r\n]+

%%

/* ============================================================
 * SECCION 3: reglas — "cuando veas esto → entregá este token"
 * Dos mecanismos DISTINTOS, no confundirlos:
 *  - Longest match: entre lexemas de DISTINTA longitud gana el mas largo,
 *    sin importar el orden de las reglas ("->" le gana a "-", "<!...!>" a "<").
 *  - Orden de declaracion: solo desempata cuando dos reglas matchean la MISMA
 *    longitud. Por eso las reservadas y operadores van ANTES que {Id} y que el
 *    comodin de simbolo: asi "CONJ" (misma longitud que un Id de 4 letras) se
 *    reconoce como palabra reservada y no como identificador.
 * ============================================================ */

/* ---- 1. Se reconoce pero se DESCARTA (no produce token) ---- */
{ComentLinea}   { /* comentario de una linea (4.3.1): ignorar */ }
{ComentMulti}   { /* comentario multilinea (4.3.2): ignorar */ }
{Blancos}       { /* espacios, tabs y saltos: solo separan */ }

/* ---- 2. Palabras reservadas (case sensitive: solo en MAYUSCULAS) ---- */
"CONJ"                    { return symbol(sym.CONJ); }
"OPERA"                   { return symbol(sym.OPERA); }
"EVALUAR"                 { return symbol(sym.EVALUAR); }

/* ---- 3. Operadores de conjuntos (notacion prefija, 4.6) ----
   Se reservan como tokens: NO pueden usarse como elementos ni como
   nombres (restriccion documentada en docs/gramatica.txt). */
"U"                       { return symbol(sym.UNION); }
"&"                       { return symbol(sym.INTERSECCION); }
"^"                       { return symbol(sym.COMPLEMENTO); }
"-"                       { return symbol(sym.DIFERENCIA); }

/* ---- 4. Simbolos estructurales (el longest match resuelve "->" vs "-") ---- */
"->"                      { return symbol(sym.FLECHA); }
"~"                       { return symbol(sym.VIRGULILLA); }
"{"                       { return symbol(sym.LLAVE_IZQ); }
"}"                       { return symbol(sym.LLAVE_DER); }
"("                       { return symbol(sym.PAR_IZQ); }
")"                       { return symbol(sym.PAR_DER); }
":"                       { return symbol(sym.DOS_PUNTOS); }
";"                       { return symbol(sym.PUNTO_COMA); }
","                       { return symbol(sym.COMA); }

/* ---- 5. Tokens con patron (despues de reservadas y operadores) ---- */
{Id}                      { return symbol(sym.ID); }
{Numero}                  { return symbol(sym.NUMERO); }

/* ---- 6. Comodin: cualquier OTRO caracter del universo ASCII 33..126
        que no encajo arriba se toma como SIMBOLO (posible elemento de un
        conjunto: '!', '@', '$', '<', '>', etc.). ---- */
[\x21-\x7E]               { return symbol(sym.SIMBOLO); }

/* ---- 7. Cualquier caracter FUERA del universo = error lexico ----
   Va a la tabla de errores (5.3). El analisis CONTINUA: el caracter
   intruso se descarta y se sigue escaneando. */
[^]                       { if (entorno != null) {
                              entorno.error("Lexico", "el caracter '" + yytext()
                                  + "' no pertenece al lenguaje", yyline, yycolumn);
                            } else {
                              System.err.printf(
                                  "Error lexico: '%s' no pertenece al lenguaje (linea %d, columna %d)%n",
                                  yytext(), yyline + 1, yycolumn + 1);
                            } }
