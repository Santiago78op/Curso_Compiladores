/* ============================================================
 * Lexer.flex — Analizador lexico de CompScript (OLC1, PT1 VD2024)
 *
 * SECCION 1: codigo de usuario
 * ============================================================ */
package compscript.analisis;

import java_cup.runtime.Symbol;

%%

/* ============================================================
 * SECCION 2: opciones y definiciones regulares
 * ============================================================ */
%class Lexer
%public
%unicode
%cup
%line
%column
%ignorecase

%{
  /* El Interprete conecta el contexto tras crear el lexer: permite
     registrar cada token (reporte 6.1) y mandar los errores lexicos a
     la tabla de errores (6.2) en vez de a stderr. */
  public compscript.interprete.Contexto contexto;

  private Symbol symbol(int type) {
    registrar(type);
    return new Symbol(type, yyline, yycolumn, yytext());
  }
  private Symbol symbol(int type, Object value) {
    registrar(type);
    return new Symbol(type, yyline, yycolumn, value);
  }
  private void registrar(int type) {
    if (contexto != null) contexto.registrarToken(yytext(), type, yyline, yycolumn);
  }

  /* Procesa las secuencias de escape del enunciado (5.4) dentro de
     cadenas y caracteres: \n \t \\ \" \' */
  private String procesarEscapes(String s) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char ch = s.charAt(i);
      if (ch == '\\' && i + 1 < s.length()) {
        char sig = s.charAt(++i);
        switch (sig) {
          case 'n': sb.append('\n'); break;
          case 't': sb.append('\t'); break;
          case '\\': sb.append('\\'); break;
          case '"': sb.append('"'); break;
          case '\'': sb.append('\''); break;
          default: sb.append(sig);   // escape desconocido: literal
        }
      } else {
        sb.append(ch);
      }
    }
    return sb.toString();
  }
%}

/* --- definiciones regulares --- */
Letra        = [a-zA-Z_]
Digito       = [0-9]
Id           = {Letra}({Letra}|{Digito})*
Entero       = {Digito}+
Decimal      = {Digito}+"."{Digito}+
Cadena       = \"([^\"\\\n]|\\.)*\"
Caracter     = '([^'\\\n]|\\.)'
ComentLinea  = "//"[^\r\n]*
ComentMulti  = "/*"~"*/"
Blancos      = [ \t\r\n\f]+

%%

/* ============================================================
 * SECCION 3: reglas
 * ============================================================ */

/* ---- 1. Se reconoce y se DESCARTA ---- */
{ComentLinea}   { /* comentario de linea: ignorar */ }
{ComentMulti}   { /* comentario multilinea: ignorar */ }
{Blancos}       { /* espacios: solo separan */ }

/* ---- 2. Palabras reservadas (case-insensitive por %ignorecase).
        SIEMPRE antes que {Id} ---- */
"int"       { return symbol(sym.T_INT); }
"double"    { return symbol(sym.T_DOUBLE); }
"bool"      { return symbol(sym.T_BOOL); }
"char"      { return symbol(sym.T_CHAR); }
"string"    { return symbol(sym.T_STRING); }
"void"      { return symbol(sym.VOID); }

"let"       { return symbol(sym.LET); }
"const"     { return symbol(sym.CONST); }
"true"      { return symbol(sym.TRUE); }
"false"     { return symbol(sym.FALSE); }

"if"        { return symbol(sym.IF); }
"else"      { return symbol(sym.ELSE); }
"match"     { return symbol(sym.MATCH); }
"default"   { return symbol(sym.DEFAULT); }

"while"     { return symbol(sym.WHILE); }
"for"       { return symbol(sym.FOR); }
"do"        { return symbol(sym.DO); }

"break"     { return symbol(sym.BREAK); }
"continue"  { return symbol(sym.CONTINUE); }
"return"    { return symbol(sym.RETURN); }

"struct"    { return symbol(sym.STRUCT); }
"list"      { return symbol(sym.LIST); }
"cast"      { return symbol(sym.CAST); }
"as"        { return symbol(sym.AS); }

"console"   { return symbol(sym.CONSOLE); }
"log"       { return symbol(sym.LOG); }
"push"      { return symbol(sym.PUSH); }
"get"       { return symbol(sym.GET); }
"set"       { return symbol(sym.SET); }
"remove"    { return symbol(sym.REMOVE); }
"pop"       { return symbol(sym.POP); }
"reverse"   { return symbol(sym.REVERSE); }

"round"     { return symbol(sym.ROUND); }
"length"    { return symbol(sym.LENGTH); }
"tostring"  { return symbol(sym.TOSTRING); }
"run_main"  { return symbol(sym.RUN_MAIN); }

/* ---- 3. Simbolos (longest match resuelve == vs =, => vs =, etc.) ---- */
"++"        { return symbol(sym.INCR); }
"--"        { return symbol(sym.DECR); }
"=="        { return symbol(sym.IGUAL_IGUAL); }
"!="        { return symbol(sym.DIFERENTE); }
"<="        { return symbol(sym.MENOR_IGUAL); }
">="        { return symbol(sym.MAYOR_IGUAL); }
"&&"        { return symbol(sym.AND); }
"||"        { return symbol(sym.OR); }
"=>"        { return symbol(sym.FLECHA); }

"+"         { return symbol(sym.MAS); }
"-"         { return symbol(sym.MENOS); }
"*"         { return symbol(sym.POR); }
"/"         { return symbol(sym.DIV); }
"^"         { return symbol(sym.POT); }
"$"         { return symbol(sym.RAIZ); }
"%"         { return symbol(sym.MOD); }
"!"         { return symbol(sym.NOT); }
"<"         { return symbol(sym.MENOR); }
">"         { return symbol(sym.MAYOR); }
"="         { return symbol(sym.IGUAL); }

"("         { return symbol(sym.PAR_IZQ); }
")"         { return symbol(sym.PAR_DER); }
"["         { return symbol(sym.COR_IZQ); }
"]"         { return symbol(sym.COR_DER); }
"{"         { return symbol(sym.LLAVE_IZQ); }
"}"         { return symbol(sym.LLAVE_DER); }
","         { return symbol(sym.COMA); }
";"         { return symbol(sym.PUNTO_COMA); }
":"         { return symbol(sym.DOS_PUNTOS); }
"."         { return symbol(sym.PUNTO); }

/* ---- 4. Tokens con patron (despues de reservadas) ---- */
{Decimal}   { return symbol(sym.DECIMAL, Double.valueOf(yytext())); }
{Entero}    { return symbol(sym.ENTERO, Integer.valueOf(yytext())); }
{Cadena}    { String v = procesarEscapes(yytext().substring(1, yylength()-1));
              return symbol(sym.CADENA, v); }
{Caracter}  { String v = procesarEscapes(yytext().substring(1, yylength()-1));
              return symbol(sym.CARACTER, Character.valueOf(v.charAt(0))); }
{Id}        { return symbol(sym.ID, yytext()); }

/* ---- 5. Cualquier otro caracter = error lexico. El analisis CONTINUA:
        se descarta el caracter intruso y se sigue escaneando ---- */
[^]         { if (contexto != null) {
                contexto.error("Lexico", "el caracter '" + yytext()
                    + "' no pertenece al lenguaje", yyline, yycolumn);
              } else {
                System.err.printf("Error lexico: '%s' (linea %d, col %d)%n",
                    yytext(), yyline + 1, yycolumn + 1);
              } }
