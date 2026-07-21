/* ============================================================
   Colector de errores (lexico / sintactico / semantico).
   Un solo arreglo para toda la ejecucion (seccion 6.1).
   ============================================================ */

class ListaErrores {
  constructor() {
    this.errores = [];
  }
  agregar(tipo, descripcion, linea, columna) {
    this.errores.push({
      tipo: tipo,
      descripcion: descripcion,
      linea: linea != null ? linea : 0,
      columna: columna != null ? columna : 0
    });
  }
  lexico(desc, l, c) { this.agregar('Léxico', desc, l, c); }
  sintactico(desc, l, c) { this.agregar('Sintáctico', desc, l, c); }
  semantico(desc, l, c) { this.agregar('Semántico', desc, l, c); }
  hay() { return this.errores.length > 0; }
}

/* Error semantico que corta la evaluacion de una expresion.
   Se captura arriba; el mecanismo principal es propagacion por null. */
class ErrorSemantico extends Error {
  constructor(descripcion, linea, columna) {
    super(descripcion);
    this.esSemantico = true;
    this.descripcion = descripcion;
    this.linea = linea;
    this.columna = columna;
  }
}

module.exports = { ListaErrores, ErrorSemantico };
