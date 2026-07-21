/* ============================================================
   Entorno: tabla de simbolos con cadena de ambitos.
   global -> funcion/metodo -> bloque (if/while/for/...).
   Claves en minuscula: el lenguaje es case insensitive (5.1),
   lo que alcanza tambien a los identificadores.
   ============================================================ */

class Simbolo {
  constructor(id, categoria, tipoDato, valor, mutable, linea, columna) {
    this.id = id;                 // nombre original
    this.categoria = categoria;   // 'Variable' | 'Vector' | 'Funcion' | 'Metodo' | 'Parametro'
    this.tipoDato = tipoDato;     // 'int' | 'double' | ... | 'vector'
    this.valor = valor;           // Valor
    this.mutable = mutable;       // let=true, const=false
    this.linea = linea;
    this.columna = columna;
  }
}

class Entorno {
  constructor(padre, nombre) {
    this.padre = padre || null;
    this.nombre = nombre || 'global';
    this.tabla = new Map();
  }

  clave(id) { return String(id).toLowerCase(); }

  // Declara en ESTE entorno. Devuelve false si ya existe en este mismo ambito.
  declarar(simbolo) {
    const k = this.clave(simbolo.id);
    if (this.tabla.has(k)) return false;
    this.tabla.set(k, simbolo);
    return true;
  }

  // Busca subiendo por la cadena de entornos.
  obtener(id) {
    const k = this.clave(id);
    let e = this;
    while (e) {
      if (e.tabla.has(k)) return e.tabla.get(k);
      e = e.padre;
    }
    return null;
  }

  // Solo en este entorno.
  obtenerLocal(id) {
    return this.tabla.get(this.clave(id)) || null;
  }

  existe(id) { return this.obtener(id) !== null; }
}

module.exports = { Entorno, Simbolo };
