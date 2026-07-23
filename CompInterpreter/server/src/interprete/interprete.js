/* ============================================================
   Interprete: recorre el AST en DOS PASADAS (seccion 5.21):
     Pasada 1 -> registrar todas las funciones/metodos (forward-reference).
     Pasada 2 -> ejecutar declaraciones globales y la sentencia `ejecutar`.
   Analisis semantico (tablas de tipos) + ejecucion en un solo recorrido.
   Control de flujo con senales: BREAK / CONTINUE / RETURN.
   ============================================================ */
const { TIPO, valorPorDefecto } = require('./tipos');
const { Valor, aNumero, aTexto } = require('./valor');
const { Entorno, Simbolo } = require('./entorno');
const ops = require('./operaciones');
const { ejecutarNativa } = require('./nativas');

const MAX_ITER = 1000000;   // guarda anti-bucle-infinito
const MAX_DEPTH = 2000;     // guarda anti-recursion desbordada

class Interprete {
  constructor(errores) {
    this.errores = errores;
    this.consola = [];
    this.funciones = new Map();      // id(lower) -> nodo FUNCION/METODO
    this.simbolos = new Map();       // "entorno::id" -> fila de reporte
    this.global = new Entorno(null, 'global');
    this.profundidad = 0;
  }

  imprimir(txt) { this.consola.push(txt); }

  // Registra/actualiza una fila del reporte de tabla de simbolos.
  // Usa el entorno DONDE fue declarado el simbolo (guardado en el propio
  // simbolo), de modo que reasignar un global dentro de una funcion no
  // duplica la fila bajo el nombre de esa funcion.
  registrarSimbolo(simbolo, entorno) {
    if (entorno && !simbolo.entornoNombre) simbolo.entornoNombre = entorno.nombre;
    const amb = simbolo.entornoNombre || (entorno ? entorno.nombre : 'global');
    const clave = amb + '::' + String(simbolo.id).toLowerCase();
    this.simbolos.set(clave, {
      id: simbolo.id,
      categoria: simbolo.categoria,
      tipoDato: simbolo.tipoDato,
      entorno: amb,
      valor: simbolo.valor ? aTexto(simbolo.valor) : '',
      linea: simbolo.linea,
      columna: simbolo.columna
    });
  }

  // ---------------- Punto de entrada ----------------
  interpretar(ast) {
    if (!ast || !ast.globales) return;

    // PASADA 1: funciones y metodos
    for (const g of ast.globales) {
      if (g.tipo === 'FUNCION' || g.tipo === 'METODO') {
        const clave = g.id.toLowerCase();
        if (this.funciones.has(clave)) {
          this.errores.semantico('Ya existe una función o método con el identificador "' + g.id +
            '" (no hay sobrecarga)', g.linea, g.columna);
          continue;
        }
        this.funciones.set(clave, g);
        const cat = g.tipo === 'FUNCION' ? 'Funcion' : 'Metodo';
        const sim = new Simbolo(g.id, cat, g.tipo === 'FUNCION' ? g.retorno : 'void', null, false, g.linea, g.columna);
        this.registrarSimbolo(sim, this.global);
      }
    }

    // PASADA 2: variables globales y ejecutar
    for (const g of ast.globales) {
      if (g.tipo === 'DECLARACION' || g.tipo === 'DECLARACION_VECTOR') {
        this.ejecutarInstruccion(g, this.global);
      } else if (g.tipo === 'EJECUTAR') {
        this.ejecutarEjecutar(g);
      }
    }
  }

  // Resuelve la sentencia `ejecutar <id>(...)`: busca el método por nombre
  // y lo invoca. Solo METODO es válido aquí (una FUNCION exige usar su
  // valor de retorno), pero igual se ejecuta para no cortar el programa.
  ejecutarEjecutar(nodo) {
    const fn = this.funciones.get(nodo.id.toLowerCase());
    if (!fn) {
      this.errores.semantico('El método "' + nodo.id + '" no está declarado', nodo.linea, nodo.columna);
      return;
    }
    if (fn.tipo !== 'METODO') {
      this.errores.semantico('ejecutar solo acepta métodos (void); "' + nodo.id + '" es una función', nodo.linea, nodo.columna);
      // aun asi lo ejecutamos para no cortar el programa
    }
    this.invocar(fn, nodo.args, nodo.linea, nodo.columna);
  }

  // ---------------- Ejecucion de bloques ----------------
  // Devuelve una senal ({tipo:'BREAK'|'CONTINUE'|'RETURN', valor}) o null.
  ejecutarBloque(instrucciones, entorno) {
    for (const inst of instrucciones) {
      const senal = this.ejecutarInstruccion(inst, entorno);
      if (senal) return senal;
    }
    return null;
  }

  // Despacha una instrucción por nodo.tipo. Devuelve la señal de control
  // que suba desde su ejecución (BREAK/CONTINUE/RETURN) o null.
  ejecutarInstruccion(nodo, entorno) {
    if (!nodo) return null;
    switch (nodo.tipo) {
      case 'DECLARACION': return this.execDeclaracion(nodo, entorno);
      case 'DECLARACION_VECTOR': return this.execDeclaracionVector(nodo, entorno);
      case 'ASIGNACION': return this.execAsignacion(nodo, entorno);
      case 'ASIGNACION_VECTOR': return this.execAsignacionVector(nodo, entorno);
      case 'INCREMENTO': return this.execIncDec(nodo, entorno, 1);
      case 'DECREMENTO': return this.execIncDec(nodo, entorno, -1);
      case 'ECHO': return this.execEcho(nodo, entorno);
      case 'IF': return this.execIf(nodo, entorno);
      case 'SWITCH': return this.execSwitch(nodo, entorno);
      case 'WHILE': return this.execWhile(nodo, entorno);
      case 'FOR': return this.execFor(nodo, entorno);
      case 'DOUNTIL': return this.execDoUntil(nodo, entorno);
      case 'LOOP': return this.execLoop(nodo, entorno);
      case 'BREAK': return { tipo: 'BREAK', linea: nodo.linea, columna: nodo.columna };
      case 'CONTINUE': return { tipo: 'CONTINUE', linea: nodo.linea, columna: nodo.columna };
      case 'RETURN': {
        const v = nodo.valor ? this.evaluar(nodo.valor, entorno) : null;
        return { tipo: 'RETURN', valor: nodo.valor ? v : null, tieneValor: !!nodo.valor, linea: nodo.linea, columna: nodo.columna };
      }
      case 'LLAMADA': this.evalLlamada(nodo, entorno, false); return null;    // llamada como instruccion (metodo void permitido)
      case 'NATIVA': this.evaluar(nodo, entorno); return null;
      case 'EJECUTAR': this.ejecutarEjecutar(nodo); return null;
      case 'ERROR': return null;
      default:
        return null;
    }
  }

  // ---------------- Declaraciones ----------------
  execDeclaracion(nodo, entorno) {
    let valorInicial = null;
    let tieneValor = nodo.valor != null;
    if (tieneValor) {
      valorInicial = this.evaluar(nodo.valor, entorno);
    }
    for (const id of nodo.ids) {
      let v;
      if (tieneValor) {
        v = valorInicial === null ? null : this.coercionar(nodo.tipoDato, valorInicial, nodo.linea, nodo.columna);
        if (v === null && valorInicial !== null) {
          // error de coercion ya reportado; se guarda con valor por defecto para continuar
          v = new Valor(nodo.tipoDato, valorPorDefecto(nodo.tipoDato));
        }
      } else {
        v = new Valor(nodo.tipoDato, valorPorDefecto(nodo.tipoDato));
      }
      const sim = new Simbolo(id, 'Variable', nodo.tipoDato, v, nodo.mutable, nodo.linea, nodo.columna);
      if (!entorno.declarar(sim)) {
        this.errores.semantico('La variable "' + id + '" ya fue declarada en este entorno', nodo.linea, nodo.columna);
      } else {
        this.registrarSimbolo(sim, entorno);
      }
    }
    return null;
  }

  // Declara un vector; `nodo.init.tipo` indica cuál de las 5 formas de
  // inicialización se usó: VECTOR_NEW (tamaño -> vector de nulos, 1D),
  // VECTOR_NEW2 (filas x columnas -> matriz de nulos, 2D), VECTOR_LISTA /
  // VECTOR_LISTA2 (literal explícito, 1D/2D) o VECTOR_EXPR (copia
  // superficial de otro vector ya existente, validando dimensión).
  execDeclaracionVector(nodo, entorno) {
    if (nodo.ids.length !== 1) {
      this.errores.semantico('Un vector solo puede declararse con un identificador', nodo.linea, nodo.columna);
    }
    const id = nodo.ids[0];
    const init = nodo.init;
    let vectorVal = null;

    if (init.tipo === 'VECTOR_NEW') {
      const tam = this.evaluar(init.tam, entorno);
      if (tam !== null) {
        if (tam.tipo !== TIPO.INT) this.errores.semantico('El tamaño del vector debe ser ENTERO', nodo.linea, nodo.columna);
        const n = Math.max(0, Math.trunc(aNumero(tam)));
        const arr = [];
        for (let i = 0; i < n; i++) arr.push(Valor.nulo());
        vectorVal = Valor.vector(nodo.tipoDato, 1, arr);
      }
    } else if (init.tipo === 'VECTOR_LISTA') {
      const arr = [];
      for (const e of init.valores) {
        const ev = this.evaluar(e, entorno);
        arr.push(ev === null ? Valor.nulo() : this.coercionarElementoOValorDefecto(nodo.tipoDato, ev, nodo.linea, nodo.columna));
      }
      vectorVal = Valor.vector(nodo.tipoDato, 1, arr);
    } else if (init.tipo === 'VECTOR_NEW2') {
      const f = this.evaluar(init.filas, entorno);
      const c = this.evaluar(init.cols, entorno);
      if (f !== null && c !== null) {
        const nf = Math.max(0, Math.trunc(aNumero(f)));
        const nc = Math.max(0, Math.trunc(aNumero(c)));
        const mat = [];
        for (let i = 0; i < nf; i++) {
          const fila = [];
          for (let j = 0; j < nc; j++) fila.push(Valor.nulo());
          mat.push(fila);
        }
        vectorVal = Valor.vector(nodo.tipoDato, 2, mat);
      }
    } else if (init.tipo === 'VECTOR_EXPR') {
      const ev = this.evaluar(init.expr, entorno);
      if (ev !== null) {
        if (ev.tipo !== TIPO.VECTOR) {
          this.errores.semantico('Se esperaba un VECTOR para inicializar "' + id + '"', nodo.linea, nodo.columna);
        } else if (ev.dimension !== nodo.dimension) {
          this.errores.semantico('El vector "' + id + '" es de ' + nodo.dimension + ' dimensión(es)', nodo.linea, nodo.columna);
        } else {
          // copia superficial conservando el tipo base declarado
          const contenido = nodo.dimension === 2 ? ev.valor.map(f => f.slice()) : ev.valor.slice();
          vectorVal = Valor.vector(nodo.tipoDato, nodo.dimension, contenido);
        }
      }
    } else if (init.tipo === 'VECTOR_LISTA2') {
      const mat = [];
      for (const fila of init.filas) {
        const filaArr = [];
        for (const e of fila) {
          const ev = this.evaluar(e, entorno);
          filaArr.push(ev === null ? Valor.nulo() : this.coercionarElementoOValorDefecto(nodo.tipoDato, ev, nodo.linea, nodo.columna));
        }
        mat.push(filaArr);
      }
      vectorVal = Valor.vector(nodo.tipoDato, 2, mat);
    }

    if (vectorVal === null) vectorVal = Valor.vector(nodo.tipoDato, nodo.dimension, []);
    const sim = new Simbolo(id, 'Vector', nodo.tipoDato, vectorVal, nodo.mutable, nodo.linea, nodo.columna);
    if (!entorno.declarar(sim)) {
      this.errores.semantico('El vector "' + id + '" ya fue declarado en este entorno', nodo.linea, nodo.columna);
    } else {
      this.registrarSimbolo(sim, entorno);
    }
    return null;
  }

  // ---------------- Asignaciones ----------------
  execAsignacion(nodo, entorno) {
    const sim = entorno.obtener(nodo.id);
    if (!sim) {
      this.errores.semantico('La variable "' + nodo.id + '" no está declarada', nodo.linea, nodo.columna);
      return null;
    }
    if (!sim.mutable) {
      this.errores.semantico('No se puede reasignar la constante "' + nodo.id + '"', nodo.linea, nodo.columna);
      return null;
    }
    const v = this.evaluar(nodo.valor, entorno);
    if (v === null) return null;
    const coer = this.coercionar(sim.tipoDato, v, nodo.linea, nodo.columna);
    if (coer === null) return null;
    sim.valor = coer;
    this.registrarSimbolo(sim, entorno);
    return null;
  }

  execAsignacionVector(nodo, entorno) {
    const sim = entorno.obtener(nodo.id);
    if (!sim) {
      this.errores.semantico('El vector "' + nodo.id + '" no está declarado', nodo.linea, nodo.columna);
      return null;
    }
    if (!sim.mutable) {
      this.errores.semantico('No se puede modificar la constante "' + nodo.id + '"', nodo.linea, nodo.columna);
      return null;
    }
    if (!sim.valor || sim.valor.tipo !== TIPO.VECTOR) {
      this.errores.semantico('"' + nodo.id + '" no es un vector', nodo.linea, nodo.columna);
      return null;
    }
    const i1 = this.evaluar(nodo.idx1, entorno);
    const nuevo = this.evaluar(nodo.valor, entorno);
    if (i1 === null || nuevo === null) return null;
    const idx1 = Math.trunc(aNumero(i1));

    if (nodo.idx2 === null) {
      if (sim.valor.dimension !== 1) { this.errores.semantico('El vector "' + nodo.id + '" es de 2 dimensiones', nodo.linea, nodo.columna); return null; }
      if (idx1 < 0 || idx1 >= sim.valor.valor.length) { this.errores.semantico('Índice fuera de rango en "' + nodo.id + '"', nodo.linea, nodo.columna); return null; }
      const coer = this.coercionarElemento(sim.valor.tipoBase, nuevo, nodo.linea, nodo.columna);
      if (coer !== null) sim.valor.valor[idx1] = coer;
    } else {
      const i2 = this.evaluar(nodo.idx2, entorno);
      if (i2 === null) return null;
      const idx2 = Math.trunc(aNumero(i2));
      if (sim.valor.dimension !== 2) { this.errores.semantico('El vector "' + nodo.id + '" no es de 2 dimensiones', nodo.linea, nodo.columna); return null; }
      const mat = sim.valor.valor;
      if (idx1 < 0 || idx1 >= mat.length || idx2 < 0 || idx2 >= mat[idx1].length) {
        this.errores.semantico('Índice fuera de rango en "' + nodo.id + '"', nodo.linea, nodo.columna); return null;
      }
      const coer = this.coercionarElemento(sim.valor.tipoBase, nuevo, nodo.linea, nodo.columna);
      if (coer !== null) mat[idx1][idx2] = coer;
    }
    this.registrarSimbolo(sim, entorno);
    return null;
  }

  execIncDec(nodo, entorno, delta) {
    const sim = entorno.obtener(nodo.id);
    if (!sim) { this.errores.semantico('La variable "' + nodo.id + '" no está declarada', nodo.linea, nodo.columna); return null; }
    if (!sim.mutable) { this.errores.semantico('No se puede modificar la constante "' + nodo.id + '"', nodo.linea, nodo.columna); return null; }
    if (sim.valor.tipo === TIPO.INT) sim.valor = Valor.int(sim.valor.valor + delta);
    else if (sim.valor.tipo === TIPO.DOUBLE) sim.valor = Valor.double(sim.valor.valor + delta);
    else { this.errores.semantico('Solo se puede incrementar/decrementar valores numéricos', nodo.linea, nodo.columna); return null; }
    this.registrarSimbolo(sim, entorno);
    return null;
  }

  execEcho(nodo, entorno) {
    const v = this.evaluar(nodo.valor, entorno);
    if (v !== null) this.imprimir(aTexto(v));
    return null;
  }

  // ---------------- Control ----------------
  execIf(nodo, entorno) {
    const cond = this.evaluar(nodo.cond, entorno);
    if (cond === null) return null;
    if (cond.tipo !== TIPO.BOOL) { this.errores.semantico('La condición del if debe ser BOOLEAN', nodo.linea, nodo.columna); return null; }
    if (cond.valor) {
      return this.ejecutarBloque(nodo.entonces, new Entorno(entorno, entorno.nombre));
    } else if (nodo.sino) {
      return this.ejecutarBloque(nodo.sino, new Entorno(entorno, entorno.nombre));
    }
    return null;
  }

  execSwitch(nodo, entorno) {
    const expr = this.evaluar(nodo.expr, entorno);
    if (expr === null) return null;
    const local = new Entorno(entorno, entorno.nombre);
    let coincidio = false;
    // fall-through (seccion 5.16.2): una vez que un case coincide, se siguen
    // ejecutando los siguientes hasta un break.
    for (const caso of nodo.casos) {
      if (!coincidio) {
        const ce = this.evaluar(caso.expr, local);
        if (ce === null) continue;
        const igual = ops.relacional('==', expr, ce, nodo.linea, nodo.columna, this.errores);
        if (igual !== null && igual.valor === true) coincidio = true;
      }
      if (coincidio) {
        const senal = this.ejecutarBloque(caso.cuerpo, local);
        if (senal) {
          if (senal.tipo === 'BREAK') return null;
          return senal;   // RETURN / CONTINUE se propagan
        }
      }
    }
    if (nodo.defecto) {
      // si no coincidio ningun case, o se llego por fall-through, se ejecuta default
      const senal = this.ejecutarBloque(nodo.defecto.cuerpo, local);
      if (senal) {
        if (senal.tipo === 'BREAK') return null;
        return senal;
      }
    }
    return null;
  }

  execWhile(nodo, entorno) {
    let n = 0;
    while (true) {
      const cond = this.evaluar(nodo.cond, entorno);
      if (cond === null) return null;
      if (cond.tipo !== TIPO.BOOL) { this.errores.semantico('La condición del while debe ser BOOLEAN', nodo.linea, nodo.columna); return null; }
      if (!cond.valor) break;
      if (++n > MAX_ITER) { this.errores.semantico('Posible ciclo infinito (while)', nodo.linea, nodo.columna); break; }
      const senal = this.ejecutarBloque(nodo.cuerpo, new Entorno(entorno, entorno.nombre));
      if (senal) {
        if (senal.tipo === 'BREAK') break;
        if (senal.tipo === 'CONTINUE') continue;
        return senal;   // RETURN
      }
    }
    return null;
  }

  execFor(nodo, entorno) {
    const local = new Entorno(entorno, entorno.nombre);
    this.ejecutarInstruccion(nodo.init, local);
    let n = 0;
    while (true) {
      const cond = this.evaluar(nodo.cond, local);
      if (cond === null) return null;
      if (cond.tipo !== TIPO.BOOL) { this.errores.semantico('La condición del for debe ser BOOLEAN', nodo.linea, nodo.columna); return null; }
      if (!cond.valor) break;
      if (++n > MAX_ITER) { this.errores.semantico('Posible ciclo infinito (for)', nodo.linea, nodo.columna); break; }
      const senal = this.ejecutarBloque(nodo.cuerpo, new Entorno(local, local.nombre));
      if (senal) {
        if (senal.tipo === 'BREAK') break;
        if (senal.tipo === 'RETURN') return senal;
        // CONTINUE cae a la actualizacion
      }
      this.ejecutarInstruccion(nodo.update, local);
    }
    return null;
  }

  execDoUntil(nodo, entorno) {
    // Ejecuta al menos una vez; TERMINA cuando la condicion se vuelve verdadera.
    let n = 0;
    while (true) {
      const senal = this.ejecutarBloque(nodo.cuerpo, new Entorno(entorno, entorno.nombre));
      if (senal) {
        if (senal.tipo === 'BREAK') break;
        if (senal.tipo === 'RETURN') return senal;
        // CONTINUE cae a la evaluacion de condicion
      }
      const cond = this.evaluar(nodo.cond, entorno);
      if (cond === null) return null;
      if (cond.tipo !== TIPO.BOOL) { this.errores.semantico('La condición del do-until debe ser BOOLEAN', nodo.linea, nodo.columna); return null; }
      if (cond.valor) break;   // termina cuando es verdadera
      if (++n > MAX_ITER) { this.errores.semantico('Posible ciclo infinito (do-until)', nodo.linea, nodo.columna); break; }
    }
    return null;
  }

  execLoop(nodo, entorno) {
    let n = 0;
    while (true) {
      if (++n > MAX_ITER) { this.errores.semantico('Posible ciclo infinito (loop)', nodo.linea, nodo.columna); break; }
      const senal = this.ejecutarBloque(nodo.cuerpo, new Entorno(entorno, entorno.nombre));
      if (senal) {
        if (senal.tipo === 'BREAK') break;
        if (senal.tipo === 'CONTINUE') continue;
        return senal;   // RETURN
      }
    }
    return null;
  }

  // ---------------- Llamadas a funcion/metodo ----------------
  // Ejecuta el cuerpo de `fn` en un entorno nuevo colgado del global (sin
  // clausura sobre el llamador): liga cada parámetro por nombre con el
  // argumento provisto o su valor por defecto, corre el cuerpo y, si es
  // FUNCION, exige que haya terminado con RETURN y coerciona el valor al
  // tipo declarado. Falla con guarda anti-recursión (MAX_DEPTH).
  invocar(fn, args, linea, columna, entornoLlamador, comoExpresion = false) {
    if (++this.profundidad > MAX_DEPTH) {
      this.errores.semantico('Recursión demasiado profunda al llamar "' + fn.id + '"', linea, columna);
      this.profundidad--;
      return null;
    }
    const llamador = entornoLlamador || this.global;
    const local = new Entorno(this.global, fn.id);

    // mapa de argumentos por nombre. Las expresiones de los argumentos se
    // evaluan en el entorno DEL LLAMADOR (ahi viven sus variables locales).
    const provistos = new Map();
    for (const a of args) provistos.set(a.nombre.toLowerCase(), a.valor);

    for (const p of fn.params) {
      let valorParam = null;
      if (provistos.has(p.id.toLowerCase())) {
        valorParam = this.evaluar(provistos.get(p.id.toLowerCase()), llamador);
      } else if (p.porDefecto) {
        valorParam = this.evaluar(p.porDefecto, local);
      } else {
        this.errores.semantico('Falta el parámetro "' + p.id + '" en la llamada a "' + fn.id + '"', linea, columna);
      }
      let v;
      if (valorParam === null) v = new Valor(p.tipoParam, valorPorDefecto(p.tipoParam));
      else {
        v = this.coercionar(p.tipoParam, valorParam, linea, columna);
        if (v === null) v = new Valor(p.tipoParam, valorPorDefecto(p.tipoParam));
      }
      const sim = new Simbolo(p.id, 'Parametro', p.tipoParam, v, true, p.linea, p.columna);
      local.declarar(sim);
      this.registrarSimbolo(sim, local);
    }

    const senal = this.ejecutarBloque(fn.cuerpo, local);
    this.profundidad--;

    // 5.18.2: "continue" (y, por la misma razón, "break") deben estar
    // dentro de un ciclo; si la señal escapa hasta el cuerpo de la función o
    // método sin que ningún ciclo la consuma, se reporta como error semántico
    // en lugar de truncar en silencio el resto de la ejecución.
    if (senal && (senal.tipo === 'BREAK' || senal.tipo === 'CONTINUE')) {
      const cual = senal.tipo === 'CONTINUE' ? 'continue' : 'break';
      this.errores.semantico('La sentencia "' + cual + '" debe estar dentro de un ciclo',
        senal.linea != null ? senal.linea : linea, senal.columna != null ? senal.columna : columna);
      return null;
    }

    if (fn.tipo === 'FUNCION') {
      if (!senal || senal.tipo !== 'RETURN') {
        this.errores.semantico('La función "' + fn.id + '" debe retornar un valor de tipo ' + fn.retorno, linea, columna);
        return null;
      }
      if (!senal.tieneValor) {
        this.errores.semantico('La función "' + fn.id + '" debe retornar un valor de tipo ' + fn.retorno, linea, columna);
        return null;
      }
      if (senal.valor === null) return null;
      const coer = this.coercionar(fn.retorno, senal.valor, linea, columna);
      return coer;
    }
    // M3: un metodo es void; un "return <expr>;" en su cuerpo es error
    // semantico (igual que CompScript), en vez de descartarse en silencio.
    if (senal && senal.tipo === 'RETURN' && senal.tieneValor) {
      this.errores.semantico('El método "' + fn.id + '" es void; "return" no puede llevar una expresión',
        senal.linea != null ? senal.linea : linea, senal.columna != null ? senal.columna : columna);
    }
    // M2: usar un metodo (void) donde se espera un valor (p.ej. "let x: int
    // = m();") es error, en vez de propagar un null silencioso que aparenta
    // un error previo que nunca ocurrio.
    if (comoExpresion) {
      this.errores.semantico('El método "' + fn.id + '" no retorna valor y no puede usarse como expresión', linea, columna);
    }
    return null;   // metodo: sin valor de retorno
  }

  // ---------------- Evaluacion de expresiones ----------------
  // Despacha una expresión por nodo.tipo y devuelve su Valor (o null si
  // hubo un error semántico ya reportado más abajo en la recursión).
  evaluar(nodo, entorno) {
    if (!nodo) return null;
    switch (nodo.tipo) {
      case 'LITERAL': return this.evalLiteral(nodo);
      case 'IDENTIFICADOR': {
        const sim = entorno.obtener(nodo.nombre);
        if (!sim) { this.errores.semantico('La variable "' + nodo.nombre + '" no está declarada', nodo.linea, nodo.columna); return null; }
        return sim.valor;
      }
      case 'BINARIA': {
        const i = this.evaluar(nodo.izq, entorno), d = this.evaluar(nodo.der, entorno);
        return ops.aritmetica(nodo.op, i, d, nodo.linea, nodo.columna, this.errores);
      }
      case 'UNARIA':
        return ops.negacion(this.evaluar(nodo.expr, entorno), nodo.linea, nodo.columna, this.errores);
      case 'RELACIONAL': {
        const i = this.evaluar(nodo.izq, entorno), d = this.evaluar(nodo.der, entorno);
        return ops.relacional(nodo.op, i, d, nodo.linea, nodo.columna, this.errores);
      }
      case 'LOGICA': {
        const i = this.evaluar(nodo.izq, entorno);
        // Cortocircuito de && y ||: si el operando izquierdo (BOOL) ya decide
        // el resultado, el derecho NO se evalua. Importa cuando la derecha
        // podria fallar (p.ej. la guarda "x != 0 && 10/x > 1") o tener
        // efectos; ademas es la semantica habitual de estos operadores. El
        // caso ! es unario (nodo.der === null) y no entra aqui.
        if (i !== null && i.tipo === TIPO.BOOL) {
          if (nodo.op === '&&' && i.valor === false) return Valor.bool(false);
          if (nodo.op === '||' && i.valor === true) return Valor.bool(true);
        }
        const d = nodo.der ? this.evaluar(nodo.der, entorno) : null;
        return ops.logica(nodo.op, i, d, nodo.linea, nodo.columna, this.errores);
      }
      case 'IS': return this.evalIs(nodo, entorno);
      case 'CAST': return this.evalCast(nodo, entorno);
      case 'TERNARIO': return this.evalTernario(nodo, entorno);
      case 'LLAMADA': return this.evalLlamada(nodo, entorno, true);    // contexto expresion: un metodo void aqui es error (M2)
      case 'NATIVA': return ejecutarNativa(nodo.fn, this.evaluar(nodo.arg, entorno), nodo.linea, nodo.columna, this.errores);
      case 'ACCESO_VECTOR': return this.evalAcceso(nodo, entorno);
      default:
        return null;
    }
  }

  evalLiteral(nodo) {
    switch (nodo.tipoLit) {
      case 'int': return Valor.int(parseInt(nodo.valor, 10));
      case 'double': return Valor.double(parseFloat(nodo.valor));
      case 'string': return Valor.string(nodo.valor);
      case 'char': return Valor.char(nodo.valor);
      case 'bool': return Valor.bool(nodo.valor);
      case 'null': return Valor.nulo();
      default: return null;
    }
  }

  evalIs(nodo, entorno) {
    const v = this.evaluar(nodo.expr, entorno);
    if (v === null) return null;
    return Valor.bool(v.tipo === nodo.tipoConsulta);
  }

  evalTernario(nodo, entorno) {
    const cond = this.evaluar(nodo.cond, entorno);
    if (cond === null) return null;
    if (cond.tipo !== TIPO.BOOL) { this.errores.semantico('La condición del operador ternario debe ser BOOLEAN', nodo.linea, nodo.columna); return null; }
    return cond.valor ? this.evaluar(nodo.verdadero, entorno) : this.evaluar(nodo.falso, entorno);
  }

  evalLlamada(nodo, entorno, comoExpresion = false) {
    const fn = this.funciones.get(nodo.id.toLowerCase());
    if (!fn) { this.errores.semantico('La función o método "' + nodo.id + '" no está declarado', nodo.linea, nodo.columna); return null; }
    return this.invocar(fn, nodo.args, nodo.linea, nodo.columna, entorno, comoExpresion);
  }

  // Lee un elemento de vector/matriz: valida que el id sea un vector, que
  // la cantidad de índices coincida con su dimensión (1D/2D) y que el/los
  // índice(s) esté(n) en rango antes de devolver el Valor almacenado.
  evalAcceso(nodo, entorno) {
    const sim = entorno.obtener(nodo.id);
    if (!sim) { this.errores.semantico('El vector "' + nodo.id + '" no está declarado', nodo.linea, nodo.columna); return null; }
    if (!sim.valor || sim.valor.tipo !== TIPO.VECTOR) { this.errores.semantico('"' + nodo.id + '" no es un vector', nodo.linea, nodo.columna); return null; }
    const i1v = this.evaluar(nodo.idx1, entorno);
    if (i1v === null) return null;
    const idx1 = Math.trunc(aNumero(i1v));
    if (nodo.idx2 === null) {
      if (sim.valor.dimension !== 1) { this.errores.semantico('El vector "' + nodo.id + '" es de 2 dimensiones, falta un índice', nodo.linea, nodo.columna); return null; }
      if (idx1 < 0 || idx1 >= sim.valor.valor.length) { this.errores.semantico('Índice fuera de rango en "' + nodo.id + '"', nodo.linea, nodo.columna); return null; }
      return sim.valor.valor[idx1];
    } else {
      const i2v = this.evaluar(nodo.idx2, entorno);
      if (i2v === null) return null;
      const idx2 = Math.trunc(aNumero(i2v));
      if (sim.valor.dimension !== 2) { this.errores.semantico('El vector "' + nodo.id + '" no es de 2 dimensiones', nodo.linea, nodo.columna); return null; }
      const mat = sim.valor.valor;
      if (idx1 < 0 || idx1 >= mat.length || idx2 < 0 || idx2 >= mat[idx1].length) { this.errores.semantico('Índice fuera de rango en "' + nodo.id + '"', nodo.linea, nodo.columna); return null; }
      return mat[idx1][idx2];
    }
  }

  // ---------------- Casteos (seccion 5.13) ----------------
  evalCast(nodo, entorno) {
    const v = this.evaluar(nodo.expr, entorno);
    if (v === null) return null;
    const dest = nodo.tipoDestino;
    const o = v.tipo;
    if (o === dest) return v;
    // combinaciones permitidas
    if (o === TIPO.INT && dest === TIPO.DOUBLE) return Valor.double(v.valor);
    if (o === TIPO.DOUBLE && dest === TIPO.INT) return Valor.int(Math.trunc(v.valor));
    if (o === TIPO.INT && dest === TIPO.STRING) return Valor.string(String(v.valor));
    if (o === TIPO.INT && dest === TIPO.CHAR) return Valor.char(String.fromCharCode(v.valor));
    if (o === TIPO.DOUBLE && dest === TIPO.STRING) return Valor.string(aTexto(v));
    if (o === TIPO.CHAR && dest === TIPO.INT) return Valor.int(v.valor.charCodeAt(0));
    if (o === TIPO.CHAR && dest === TIPO.DOUBLE) return Valor.double(v.valor.charCodeAt(0));
    this.errores.semantico('No se permite el casteo de ' + ops.nombreTipo(o) + ' a ' + ops.nombreTipo(dest), nodo.linea, nodo.columna);
    return null;
  }

  // ---------------- Coercion de asignacion ----------------
  /* DECISION DE DISENO: los tipos numericos son mutuamente asignables con
     coercion (double->int trunca). Sin esto, el propio ejemplo oficial del
     Anexo 11.1 seria invalido: `let modulo:int = x % 3` donde `%` produce
     DECIMAL segun la tabla literal del enunciado (5.5.6). */
  coercionar(tipoDest, v, l, c) {
    if (v === null) return null;
    const o = v.tipo;
    if (o === tipoDest) return v;
    switch (tipoDest) {
      case TIPO.INT:
        if (o === TIPO.DOUBLE) return Valor.int(Math.trunc(v.valor));
        if (o === TIPO.BOOL) return Valor.int(v.valor ? 1 : 0);
        break;
      case TIPO.DOUBLE:
        if (o === TIPO.INT) return Valor.double(v.valor);
        if (o === TIPO.BOOL) return Valor.double(v.valor ? 1 : 0);
        if (o === TIPO.CHAR) return Valor.double(v.valor.charCodeAt(0));
        break;
      case TIPO.BOOL:
        break;   // solo bool acepta bool
      case TIPO.CHAR:
        break;   // solo char acepta char
      case TIPO.STRING:
        break;   // solo string acepta string
    }
    this.errores.semantico('No se puede asignar un valor de tipo ' + ops.nombreTipo(o) +
      ' a una variable de tipo ' + ops.nombreTipo(tipoDest), l, c);
    return null;
  }

  // coercion de un elemento de vector (mismo criterio)
  coercionarElemento(tipoBase, v, l, c) {
    return this.coercionar(tipoBase, v, l, c);
  }

  // Igual que coercionarElemento, pero para usarse al RELLENAR un vector
  // (listas literales 1D/2D): si la coercion falla, el error semantico ya
  // quedo reportado por coercionar(), pero esa posicion del vector NO puede
  // quedar en JS null (a diferencia de una expresion suelta que se propaga
  // por null): los elementos de un vector deben ser siempre instancias de
  // Valor, porque las nativas (sum/max/min/average) y aTexto() leen
  // directamente `.tipo` de cada elemento sin volver a chequear null. Se
  // rellena con el valor por defecto del tipo base para poder seguir.
  coercionarElementoOValorDefecto(tipoBase, v, l, c) {
    const coer = this.coercionar(tipoBase, v, l, c);
    return coer !== null ? coer : new Valor(tipoBase, valorPorDefecto(tipoBase));
  }
}

module.exports = { Interprete };
