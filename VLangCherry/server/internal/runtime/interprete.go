// Package runtime es el interprete de VLangCherry: tipos y valores en
// tiempo de ejecucion (tipos.go, valores.go), tabla de simbolos anidada
// (entorno.go), aritmetica/comparaciones (operaciones.go), funciones
// embebidas y conversiones (nativas.go), recoleccion de errores
// (errores.go) y el evaluador de dos pasadas (interprete.go).
package runtime

import (
	"fmt"
	"sort"

	"vlangcherry/internal/ast"
)

type TipoSenal int

const (
	SenalNinguna TipoSenal = iota
	SenalBreak
	SenalContinue
	SenalReturn
)

type Senal struct {
	Tipo  TipoSenal
	Valor Valor
	// TieneValor distingue "return;" (false) de "return expr;" (true), para
	// validar el retorno contra el tipo declarado de la funcion.
	TieneValor bool
}

// Guardas de estabilidad: no las pide el enunciado, pero el interprete corre
// dentro del servidor. Un "for true {}" colgaria la peticion para siempre y
// una recursion sin caso base desbordaria la pila de Go (fatal error que el
// recover() de analizar.go NO atrapa -> caeria el proceso entero). Ambas se
// reportan como error semantico en su lugar. (mismos valores que CompInterpreter)
const (
	maxIteraciones = 1000000
	maxProfundidad = 2000
)

// FilaSimbolo es una fila del reporte 8.2 (Tabla de simbolos).
type FilaSimbolo struct {
	ID          string `json:"id"`
	TipoSimbolo string `json:"categoria"`
	TipoDato    string `json:"tipoDato"`
	Ambito      string `json:"entorno"`
	Valor       string `json:"valor"`
	Linea       int    `json:"linea"`
	Columna     int    `json:"columna"`
}

// Interprete es el "ambiente fresco por ejecucion" (analogo a los otros
// proyectos del curso): dos pasadas (2.2 del enunciado) - registrar
// structs/funciones primero (permite llamadas hacia adelante, 7.1), luego
// declarar globales y correr main().
type Interprete struct {
	errores   *ListaErrores
	global    *Global
	entGlobal *Entorno
	consola   []string
	simbolos  map[string]*FilaSimbolo
	// profLoop/profSwitch cuentan cuantos "for"/"switch" anidados rodean la
	// sentencia que se esta ejecutando ahora mismo (4.8.1/4.8.2: break solo
	// es valido dentro de un ciclo o un switch, continue solo dentro de un
	// ciclo). invocarFuncion los resetea a 0 al entrar a una funcion/metodo
	// y los restaura al salir, para que un break/continue en el cuerpo de
	// una funcion llamada desde dentro de un ciclo no "herede" el contexto
	// del llamador.
	profLoop   int
	profSwitch int
	// profundidad de llamadas anidadas (recursion): guarda anti-desborde.
	profundidad int
}

func NuevoInterprete(errores *ListaErrores) *Interprete {
	return &Interprete{errores: errores, global: NuevoGlobal(), simbolos: make(map[string]*FilaSimbolo)}
}

func (in *Interprete) Consola() []string { return in.consola }

func (in *Interprete) Simbolos() []FilaSimbolo {
	filas := make([]FilaSimbolo, 0, len(in.simbolos))
	for _, f := range in.simbolos {
		filas = append(filas, *f)
	}
	// Orden de lectura del reporte 8.2: por ambito y luego por posicion en el
	// fuente (linea, columna), no por la clave interna del map.
	sort.Slice(filas, func(i, j int) bool {
		if filas[i].Ambito != filas[j].Ambito {
			return filas[i].Ambito < filas[j].Ambito
		}
		if filas[i].Linea != filas[j].Linea {
			return filas[i].Linea < filas[j].Linea
		}
		if filas[i].Columna != filas[j].Columna {
			return filas[i].Columna < filas[j].Columna
		}
		return filas[i].ID < filas[j].ID
	})
	return filas
}

func (in *Interprete) errorSemantico(desc string, linea, col int) {
	in.errores.Semantico(desc, linea, col)
}

// Interpretar es el punto de entrada (3.2: localizar y ejecutar main()).
func (in *Interprete) Interpretar(programa *ast.Programa) {
	for _, s := range programa.Structs {
		if _, dup := in.global.Structs[s.Nombre]; dup {
			in.errorSemantico("el struct \""+s.Nombre+"\" ya fue definido", s.Linea, s.Columna)
			continue
		}
		in.global.Structs[s.Nombre] = s
	}

	for _, f := range programa.Funciones {
		if f.ReceptorTipo != "" {
			if in.global.Metodos[f.ReceptorTipo] == nil {
				in.global.Metodos[f.ReceptorTipo] = make(map[string]*ast.DeclFuncion)
			}
			if _, dup := in.global.Metodos[f.ReceptorTipo][f.Nombre]; dup {
				in.errorSemantico("el struct \""+f.ReceptorTipo+"\" ya tiene un método \""+f.Nombre+"\"", f.Linea, f.Columna)
				continue
			}
			in.global.Metodos[f.ReceptorTipo][f.Nombre] = f
			in.registrarSimboloGlobal(f.Nombre, "Método", f.Linea, f.Columna)
		} else {
			if _, dup := in.global.Funciones[f.Nombre]; dup {
				in.errorSemantico("la función \""+f.Nombre+"\" ya fue declarada", f.Linea, f.Columna)
				continue
			}
			// 7.1: "las funciones, variables o structs no pueden tener el
			// mismo nombre" - los structs ya estan todos registrados en
			// este punto (se registran en la pasada anterior), asi que
			// alcanza con chequear en esta direccion.
			if _, dup := in.global.Structs[f.Nombre]; dup {
				in.errorSemantico("el nombre \""+f.Nombre+"\" ya está en uso por un struct", f.Linea, f.Columna)
				continue
			}
			in.global.Funciones[f.Nombre] = f
			in.registrarSimboloGlobal(f.Nombre, "Función", f.Linea, f.Columna)
		}
	}

	in.entGlobal = NuevoEntorno("global", nil)
	for _, g := range programa.Globales {
		in.ejecutarDeclaracion(g, in.entGlobal)
	}

	mainFn, ok := in.global.Funciones["main"]
	if !ok {
		in.errorSemantico("no se encontró la función main", 0, 0)
		return
	}
	if len(mainFn.Parametros) > 0 {
		in.errorSemantico("la función main no debe declarar parámetros", mainFn.Linea, mainFn.Columna)
		return
	}
	in.invocarFuncion(mainFn, nil, nil, mainFn.Linea, mainFn.Columna)
}

func (in *Interprete) registrarSimboloGlobal(nombre, categoria string, linea, col int) {
	clave := "global::" + nombre
	in.simbolos[clave] = &FilaSimbolo{ID: nombre, TipoSimbolo: categoria, TipoDato: "", Ambito: "global", Linea: linea, Columna: col}
}

// registrarSimbolo agrega (o reemplaza) la fila 8.2 de una declaracion.
// M4 — decision de diseño documentada: la tabla de simbolos es un registro de
// DECLARACIONES, una fila por SITIO de declaracion en el fuente. Por eso la
// clave incluye la posicion (linea:col) ademas de ambito y nombre:
//   - dos bloques hermanos de la misma funcion que declaran el mismo nombre en
//     lineas distintas ya NO se pisan (antes colisionaban por ambito::nombre);
//   - una declaracion dentro de un ciclo o de una funcion recursiva ocupa
//     SIEMPRE la misma posicion -> queda UNA fila (la ultima activacion), no
//     177: la tabla describe el programa, no cada activacion en runtime.
// La columna Valor es el valor AL DECLARAR (snapshot): una reasignacion
// posterior NO actualiza la fila. Es una eleccion consistente (declaraciones,
// no estado final); para ver valores finales se usa print.
func (in *Interprete) registrarSimbolo(nombre, categoria string, ent *Entorno, valor Valor, linea, col int) {
	clave := fmt.Sprintf("%s::%s::%d:%d", ent.Nombre, nombre, linea, col)
	in.simbolos[clave] = &FilaSimbolo{
		ID: nombre, TipoSimbolo: categoria, TipoDato: valor.Tipo.String(),
		Ambito: ent.Nombre, Valor: AValorReporte(valor), Linea: linea, Columna: col,
	}
}

// ---------------- Tipos y compatibilidad ----------------

func (in *Interprete) resolverTipoAST(t ast.TipoAST, linea, col int) Tipo {
	if t.EsSlice {
		el := in.resolverTipoAST(*t.Elemento, linea, col)
		return TipoSliceDe(el)
	}
	switch t.NombreBase {
	case "int":
		return TipoInt()
	case "float64":
		return TipoFloat()
	case "string":
		return TipoString()
	case "bool":
		return TipoBool()
	case "rune":
		return TipoRune()
	default:
		if _, ok := in.global.Structs[t.NombreBase]; !ok {
			in.errorSemantico("el tipo \""+t.NombreBase+"\" no fue definido", linea, col)
		}
		return TipoStructDe(t.NombreBase)
	}
}

// tipoCompatible/coercionar: 3.6 exige tipo exacto, pero se relaja la
// mezcla int<->float64 (documentado en docs/gramatica.txt): el enunciado
// muestra ambos casos (el ejemplo de conversion explicita f64(10+1) Y
// tablas aritmeticas con promocion automatica int+float64->float64
// asignada de vuelta a variables float64), asi que se prioriza usabilidad.
func tipoCompatible(declarado Tipo, v Valor) bool {
	if declarado.Igual(v.Tipo) {
		return true
	}
	if EsNumerico(declarado) && EsNumerico(v.Tipo) {
		return true
	}
	if v.Tipo.Base == TNil && (declarado.Base == TSlice || declarado.Base == TStruct) {
		return true
	}
	return false
}

func coercionar(declarado Tipo, v Valor) Valor {
	if EsNumerico(declarado) && EsNumerico(v.Tipo) && declarado.Base != v.Tipo.Base {
		if declarado.Base == TFloat64 {
			return Valor{Tipo: TipoFloat(), F: aFloat(v)}
		}
		return Valor{Tipo: TipoInt(), I: int64(aFloat(v))}
	}
	return v
}

// ---------------- Sentencias ----------------

func (in *Interprete) ejecutarBloque(b *ast.Bloque, ent *Entorno) Senal {
	for _, s := range b.Sentencias {
		senal := in.ejecutarSentencia(s, ent)
		if senal.Tipo != SenalNinguna {
			return senal
		}
	}
	return Senal{Tipo: SenalNinguna}
}

func (in *Interprete) ejecutarListaSentencias(sentencias []ast.Nodo, ent *Entorno) Senal {
	for _, s := range sentencias {
		senal := in.ejecutarSentencia(s, ent)
		if senal.Tipo != SenalNinguna {
			return senal
		}
	}
	return Senal{Tipo: SenalNinguna}
}

// ejecutarSentencia despacha por tipo concreto de nodo. Devuelve una Senal
// no-Ninguna solo para las sentencias que interrumpen el flujo normal
// (break/continue/return, o un bloque/if/switch/for que las propague),
// asi ejecutarBloque puede cortar la ejecucion del resto de sentencias.
func (in *Interprete) ejecutarSentencia(nodo ast.Nodo, ent *Entorno) Senal {
	switch n := nodo.(type) {
	case *ast.DeclVariable:
		in.ejecutarDeclaracion(n, ent)
	case *ast.Asignacion:
		in.ejecutarAsignacion(n, ent)
	case *ast.IncrementoDecremento:
		in.ejecutarIncDec(n, ent)
	case *ast.ExpresionSentencia:
		in.evaluar(n.Expr, ent)
	case *ast.SentenciaIf:
		return in.ejecutarIf(n, ent)
	case *ast.SentenciaSwitch:
		return in.ejecutarSwitch(n, ent)
	case *ast.SentenciaFor:
		return in.ejecutarFor(n, ent)
	case *ast.SentenciaBreak:
		// 4.8.1: "la sentencia break solo se puede usar dentro de un bucle
		// (for) o un switch... si se encuentra fuera se considerará un error".
		if in.profLoop == 0 && in.profSwitch == 0 {
			in.errorSemantico("la sentencia \"break\" no puede usarse fuera de un ciclo o un switch", n.Linea, n.Columna)
			return Senal{Tipo: SenalNinguna}
		}
		return Senal{Tipo: SenalBreak}
	case *ast.SentenciaContinue:
		// 4.8.2: "la sentencia continue solo se puede usar dentro de un
		// bucle (for)... si se encuentra fuera se considerará un error".
		if in.profLoop == 0 {
			in.errorSemantico("la sentencia \"continue\" no puede usarse fuera de un ciclo", n.Linea, n.Columna)
			return Senal{Tipo: SenalNinguna}
		}
		return Senal{Tipo: SenalContinue}
	case *ast.SentenciaReturn:
		if n.Valor == nil {
			return Senal{Tipo: SenalReturn}
		}
		v, _ := in.evaluar(n.Valor, ent)
		return Senal{Tipo: SenalReturn, Valor: v, TieneValor: true}
	case *ast.Bloque:
		return in.ejecutarBloque(n, NuevoEntorno(ent.Nombre, ent))
	}
	return Senal{Tipo: SenalNinguna}
}

// ejecutarDeclaracion resuelve las dos formas de declaracion (4.2): con
// tipo explicito ("var x tipo = expr", donde el valor debe ser compatible
// con el tipo declarado) o inferido (":=", donde el tipo lo determina la
// expresion). Sin valor inicial, la variable toma ValorPorDefecto(tipo).
func (in *Interprete) ejecutarDeclaracion(n *ast.DeclVariable, ent *Entorno) {
	if ent.DeclaradoLocal(n.Nombre) {
		in.errorSemantico("la variable \""+n.Nombre+"\" ya fue declarada en este ámbito", n.Linea, n.Columna)
		return
	}
	// 7.1: "las funciones, variables o structs no pueden tener el mismo
	// nombre" - solo aplica al ambito global (structs y funciones libres
	// solo pueden declararse ahi, 6.1/7.1); para cuando esta declaracion
	// corre, structs y funciones ya estan completamente registrados.
	if ent == in.entGlobal {
		if _, dup := in.global.Structs[n.Nombre]; dup {
			in.errorSemantico("el nombre \""+n.Nombre+"\" ya está en uso por un struct", n.Linea, n.Columna)
			return
		}
		if _, dup := in.global.Funciones[n.Nombre]; dup {
			in.errorSemantico("el nombre \""+n.Nombre+"\" ya está en uso por una función", n.Linea, n.Columna)
			return
		}
	}

	var valor Valor
	var tipoDeclarado Tipo

	if n.Inferido {
		v, ok := in.evaluar(n.Valor, ent)
		if !ok {
			return
		}
		valor = v
		tipoDeclarado = v.Tipo
	} else {
		tipoDeclarado = in.resolverTipoAST(*n.TipoVar, n.Linea, n.Columna)
		if n.Valor != nil {
			v, ok := in.evaluar(n.Valor, ent)
			if !ok {
				return
			}
			if !tipoCompatible(tipoDeclarado, v) {
				in.errorSemantico(fmt.Sprintf("no se puede asignar %s a variable de tipo %s", v.Tipo, tipoDeclarado), n.Linea, n.Columna)
				return
			}
			valor = coercionar(tipoDeclarado, v)
		} else {
			valor = ValorPorDefecto(tipoDeclarado)
		}
	}

	ent.DeclararMut(n.Nombre, valor, n.Mutable)
	in.registrarSimbolo(n.Nombre, "Variable", ent, valor, n.Linea, n.Columna)
}

// verificarReasignable aplica A1: reasignar una variable declarada sin 'mut'
// es error semantico. Solo aplica cuando el destino es la variable ENTERA
// (un identificador simple); mutar un campo o elemento (lugar.campo / lugar[i])
// se permite siempre, porque no reenlaza la variable. Devuelve false si debe
// abortarse la asignacion. Si la variable no existe, deja que resolverLugar
// reporte "no definida" (no duplica el error).
func (in *Interprete) verificarReasignable(lugar ast.Nodo, ent *Entorno, linea, col int) bool {
	id, ok := lugar.(*ast.Identificador)
	if !ok {
		return true
	}
	mutable, existe := ent.EsMutable(id.Nombre)
	if existe && !mutable {
		in.errorSemantico("no se puede reasignar \""+id.Nombre+"\": fue declarada sin \"mut\"", linea, col)
		return false
	}
	return true
}

// ejecutarAsignacion resuelve la celda mutable del lugar destino y aplica
// "=" directamente o "+="/"-=" combinando el valor actual con el nuevo via
// Suma/Resta antes de verificar compatibilidad de tipo con la variable.
func (in *Interprete) ejecutarAsignacion(n *ast.Asignacion, ent *Entorno) {
	if !in.verificarReasignable(n.Lugar, ent, n.Linea, n.Columna) {
		return
	}
	ptr, ok := in.resolverLugar(n.Lugar, ent)
	if !ok {
		return
	}
	v, ok := in.evaluar(n.Valor, ent)
	if !ok {
		return
	}

	switch n.Operador {
	case "+=":
		r, err := Suma(*ptr, v, n.Linea, n.Columna)
		if err != nil {
			in.errorSemantico(err.Error(), n.Linea, n.Columna)
			return
		}
		if !tipoCompatible(ptr.Tipo, r) {
			in.errorSemantico(fmt.Sprintf("+= produciría %s, incompatible con la variable de tipo %s", r.Tipo, ptr.Tipo), n.Linea, n.Columna)
			return
		}
		*ptr = coercionar(ptr.Tipo, r)
	case "-=":
		r, err := Resta(*ptr, v, n.Linea, n.Columna)
		if err != nil {
			in.errorSemantico(err.Error(), n.Linea, n.Columna)
			return
		}
		if !tipoCompatible(ptr.Tipo, r) {
			in.errorSemantico(fmt.Sprintf("-= produciría %s, incompatible con la variable de tipo %s", r.Tipo, ptr.Tipo), n.Linea, n.Columna)
			return
		}
		*ptr = coercionar(ptr.Tipo, r)
	default:
		if !tipoCompatible(ptr.Tipo, v) {
			in.errorSemantico(fmt.Sprintf("no se puede asignar %s a variable de tipo %s", v.Tipo, ptr.Tipo), n.Linea, n.Columna)
			return
		}
		*ptr = coercionar(ptr.Tipo, v)
	}
}

func (in *Interprete) ejecutarIncDec(n *ast.IncrementoDecremento, ent *Entorno) {
	if !in.verificarReasignable(n.Lugar, ent, n.Linea, n.Columna) {
		return
	}
	ptr, ok := in.resolverLugar(n.Lugar, ent)
	if !ok {
		return
	}
	delta := int64(1)
	if n.Operador == "--" {
		delta = -1
	}
	switch ptr.Tipo.Base {
	case TInt:
		ptr.I += delta
	case TFloat64:
		ptr.F += float64(delta)
	default:
		in.errorSemantico("++/-- solo se puede aplicar a valores numéricos", n.Linea, n.Columna)
	}
}

// resolverLugar da un puntero mutable a la celda referenciada por un
// "lugar" (ID | lugar[expr] | lugar.campo), reusando los mismos nodos
// ExprIndice/ExprCampo que las expresiones de lectura.
func (in *Interprete) resolverLugar(nodo ast.Nodo, ent *Entorno) (*Valor, bool) {
	switch n := nodo.(type) {
	case *ast.Identificador:
		ptr, ok := ent.Buscar(n.Nombre)
		if !ok {
			in.errorSemantico("la variable \""+n.Nombre+"\" no está definida en este contexto", n.Linea, n.Columna)
			return nil, false
		}
		return ptr, true
	case *ast.ExprIndice:
		base, ok := in.resolverLugar(n.Base, ent)
		if !ok {
			return nil, false
		}
		idx, ok := in.evaluar(n.Indice, ent)
		if !ok {
			return nil, false
		}
		if idx.Tipo.Base != TInt {
			in.errorSemantico("el índice debe ser de tipo int", n.Linea, n.Columna)
			return nil, false
		}
		if base.Tipo.Base != TSlice || base.Slice == nil {
			in.errorSemantico("no se puede indexar un valor que no es slice", n.Linea, n.Columna)
			return nil, false
		}
		i := int(idx.I)
		if i < 0 || i >= len(base.Slice.Elems) {
			in.errorSemantico(fmt.Sprintf("índice %d fuera de rango (tamaño %d)", i, len(base.Slice.Elems)), n.Linea, n.Columna)
			return nil, false
		}
		return &base.Slice.Elems[i], true
	case *ast.ExprCampo:
		base, ok := in.resolverLugar(n.Base, ent)
		if !ok {
			return nil, false
		}
		if base.Tipo.Base != TStruct || base.Struct == nil {
			in.errorSemantico("no se puede acceder al campo \""+n.Nombre+"\" de un valor que no es struct", n.Linea, n.Columna)
			return nil, false
		}
		campo, ok := base.Struct.Campos[n.Nombre]
		if !ok {
			in.errorSemantico("el struct \""+base.Struct.NombreTipo+"\" no tiene el campo \""+n.Nombre+"\"", n.Linea, n.Columna)
			return nil, false
		}
		return campo, true
	}
	// B1: un nodo que no es un "lugar" asignable (ID | lugar[i] | lugar.campo)
	// falla con mensaje en vez de en silencio. No deberia ocurrir con un AST
	// bien formado; es una red de seguridad ante nodos inesperados.
	l, c := nodo.Pos()
	in.errorSemantico("destino de asignación no válido", l, c)
	return nil, false
}

// ejecutarIf evalua las ramas if/else-if en orden y ejecuta la primera
// cuya condicion de bool sea verdadera (cada una en su propio ambito
// anidado); si ninguna aplica, ejecuta el else si existe.
func (in *Interprete) ejecutarIf(n *ast.SentenciaIf, ent *Entorno) Senal {
	for _, rama := range n.Ramas {
		cond, ok := in.evaluar(rama.Condicion, ent)
		if !ok {
			return Senal{Tipo: SenalNinguna}
		}
		if cond.Tipo.Base != TBool {
			l, c := rama.Condicion.Pos()
			in.errorSemantico("la condición del if debe ser de tipo bool", l, c)
			return Senal{Tipo: SenalNinguna}
		}
		if cond.B {
			return in.ejecutarBloque(rama.Cuerpo, NuevoEntorno(ent.Nombre, ent))
		}
	}
	if n.Else != nil {
		return in.ejecutarBloque(n.Else, NuevoEntorno(ent.Nombre, ent))
	}
	return Senal{Tipo: SenalNinguna}
}

// ejecutarSwitch evalua la expresion del switch y la compara (==) contra
// cada caso en orden; el primer caso igual gana. Sin "fallthrough" (no
// contemplado en el enunciado): un break explicito o implicito al final
// del caso corta el switch. Si ningun caso matchea, corre el default.
func (in *Interprete) ejecutarSwitch(n *ast.SentenciaSwitch, ent *Entorno) Senal {
	in.profSwitch++
	defer func() { in.profSwitch-- }()

	val, ok := in.evaluar(n.Expr, ent)
	if !ok {
		return Senal{Tipo: SenalNinguna}
	}
	for _, caso := range n.Casos {
		if caso.Valor == nil {
			continue
		}
		cv, ok := in.evaluar(caso.Valor, ent)
		if !ok {
			continue
		}
		igual, err := valoresIguales(val, cv)
		if err != nil {
			// 4.4.1: "cualquier otra combinación será inválida y se deberá
			// reportar el error" - antes este caso se descartaba en
			// silencio (como si simplemente no coincidiera), inconsistente
			// con como evaluarBinaria reporta el mismo error para == fuera
			// de un switch.
			l, c := caso.Valor.Pos()
			in.errorSemantico(err.Error(), l, c)
			continue
		}
		if !igual {
			continue
		}
		senal := in.ejecutarListaSentencias(caso.Sentencias, NuevoEntorno(ent.Nombre, ent))
		if senal.Tipo == SenalBreak {
			return Senal{Tipo: SenalNinguna}
		}
		return senal
	}
	for _, caso := range n.Casos {
		if caso.Valor != nil {
			continue
		}
		senal := in.ejecutarListaSentencias(caso.Sentencias, NuevoEntorno(ent.Nombre, ent))
		if senal.Tipo == SenalBreak {
			return Senal{Tipo: SenalNinguna}
		}
		return senal
	}
	return Senal{Tipo: SenalNinguna}
}

// ejecutarFor cubre las 3 formas del for (4.7.3): "condicion" (equivalente
// a un while), "clasico" (init; condicion; actualizacion, con su propio
// ambito compartido entre iteraciones) y "rango" (for i, v in slice, con
// un ambito nuevo por iteracion que declara indice y valor).
func (in *Interprete) ejecutarFor(n *ast.SentenciaFor, ent *Entorno) Senal {
	in.profLoop++
	defer func() { in.profLoop-- }()

	switch n.Forma {
	case "condicion":
		iter := 0
		for {
			cond, ok := in.evaluar(n.Condicion, ent)
			if !ok {
				return Senal{Tipo: SenalNinguna}
			}
			if cond.Tipo.Base != TBool {
				l, c := n.Condicion.Pos()
				in.errorSemantico("la condición del for debe ser de tipo bool", l, c)
				return Senal{Tipo: SenalNinguna}
			}
			if !cond.B {
				break
			}
			if iter++; iter > maxIteraciones {
				l, c := n.Condicion.Pos()
				in.errorSemantico("posible ciclo infinito en el for (se superó el máximo de iteraciones)", l, c)
				break
			}
			senal := in.ejecutarBloque(n.Cuerpo, NuevoEntorno(ent.Nombre, ent))
			if senal.Tipo == SenalBreak {
				break
			}
			if senal.Tipo == SenalReturn {
				return senal
			}
		}

	case "clasico":
		entFor := NuevoEntorno(ent.Nombre, ent)
		if n.Init != nil {
			in.ejecutarSentencia(n.Init, entFor)
		}
		iter := 0
		for {
			if iter++; iter > maxIteraciones {
				l, c := n.Pos()
				in.errorSemantico("posible ciclo infinito en el for (se superó el máximo de iteraciones)", l, c)
				break
			}
			if n.Condicion != nil {
				cond, ok := in.evaluar(n.Condicion, entFor)
				if !ok {
					return Senal{Tipo: SenalNinguna}
				}
				if cond.Tipo.Base != TBool {
					l, c := n.Condicion.Pos()
					in.errorSemantico("la condición del for debe ser de tipo bool", l, c)
					return Senal{Tipo: SenalNinguna}
				}
				if !cond.B {
					break
				}
			}
			senal := in.ejecutarBloque(n.Cuerpo, NuevoEntorno(ent.Nombre, entFor))
			if senal.Tipo == SenalBreak {
				break
			}
			if senal.Tipo == SenalReturn {
				return senal
			}
			if n.Actualizacion != nil {
				in.ejecutarSentencia(n.Actualizacion, entFor)
			}
		}

	case "rango":
		iterable, ok := in.evaluar(n.Iterable, ent)
		if !ok {
			return Senal{Tipo: SenalNinguna}
		}
		if iterable.Tipo.Base != TSlice {
			l, c := n.Iterable.Pos()
			in.errorSemantico("for ... in solo puede iterar sobre un slice", l, c)
			return Senal{Tipo: SenalNinguna}
		}
		var elems []Valor
		if iterable.Slice != nil {
			elems = iterable.Slice.Elems
		}
		for i, elem := range elems {
			entIter := NuevoEntorno(ent.Nombre, ent)
			entIter.Declarar(n.VarIndice, Valor{Tipo: TipoInt(), I: int64(i)})
			entIter.Declarar(n.VarValor, elem)
			senal := in.ejecutarBloque(n.Cuerpo, entIter)
			if senal.Tipo == SenalBreak {
				break
			}
			if senal.Tipo == SenalReturn {
				return senal
			}
		}
	}
	return Senal{Tipo: SenalNinguna}
}

// ---------------- Expresiones ----------------

func (in *Interprete) evaluar(nodo ast.Nodo, ent *Entorno) (Valor, bool) {
	switch n := nodo.(type) {
	case *ast.LiteralEntero:
		return Valor{Tipo: TipoInt(), I: n.Valor}, true
	case *ast.LiteralDecimal:
		return Valor{Tipo: TipoFloat(), F: n.Valor}, true
	case *ast.LiteralCadena:
		return Valor{Tipo: TipoString(), S: n.Valor}, true
	case *ast.LiteralRune:
		return Valor{Tipo: TipoRune(), R: n.Valor}, true
	case *ast.LiteralBool:
		return Valor{Tipo: TipoBool(), B: n.Valor}, true
	case *ast.LiteralNil:
		return Valor{Tipo: TipoNil()}, true
	case *ast.Identificador:
		ptr, ok := ent.Buscar(n.Nombre)
		if !ok {
			in.errorSemantico("la variable \""+n.Nombre+"\" no está definida en este contexto", n.Linea, n.Columna)
			return Valor{}, false
		}
		return *ptr, true
	case *ast.ExprUnaria:
		v, ok := in.evaluar(n.Operando, ent)
		if !ok {
			return Valor{}, false
		}
		var r Valor
		var err error
		if n.Operador == "!" {
			r, err = NotLogico(v)
		} else {
			r, err = NegacionUnaria(v, n.Linea, n.Columna)
		}
		if err != nil {
			in.errorSemantico(err.Error(), n.Linea, n.Columna)
			return Valor{}, false
		}
		return r, true
	case *ast.ExprBinaria:
		iz, ok := in.evaluar(n.Izquierda, ent)
		if !ok {
			return Valor{}, false
		}
		// Cortocircuito de && y ||: si el operando izquierdo ya decide el
		// resultado, la derecha NO se evalua. Importa cuando la derecha
		// podria fallar (p.ej. la guarda "x != 0 && 10/x > 1") o tener
		// efectos; ademas es la semantica habitual de estos operadores.
		if n.Operador == "&&" && iz.Tipo.Base == TBool && !iz.B {
			return Valor{Tipo: TipoBool(), B: false}, true
		}
		if n.Operador == "||" && iz.Tipo.Base == TBool && iz.B {
			return Valor{Tipo: TipoBool(), B: true}, true
		}
		der, ok := in.evaluar(n.Derecha, ent)
		if !ok {
			return Valor{}, false
		}
		return in.evaluarBinaria(n.Operador, iz, der, n.Linea, n.Columna)
	case *ast.ExprIndice:
		base, ok := in.evaluar(n.Base, ent)
		if !ok {
			return Valor{}, false
		}
		idx, ok := in.evaluar(n.Indice, ent)
		if !ok {
			return Valor{}, false
		}
		if idx.Tipo.Base != TInt {
			in.errorSemantico("el índice debe ser de tipo int", n.Linea, n.Columna)
			return Valor{}, false
		}
		if base.Tipo.Base != TSlice || base.Slice == nil {
			in.errorSemantico("no se puede indexar un valor que no es slice", n.Linea, n.Columna)
			return Valor{}, false
		}
		i := int(idx.I)
		if i < 0 || i >= len(base.Slice.Elems) {
			in.errorSemantico(fmt.Sprintf("índice %d fuera de rango (tamaño %d)", i, len(base.Slice.Elems)), n.Linea, n.Columna)
			return Valor{}, false
		}
		return base.Slice.Elems[i], true
	case *ast.ExprCampo:
		base, ok := in.evaluar(n.Base, ent)
		if !ok {
			return Valor{}, false
		}
		if base.Tipo.Base != TStruct || base.Struct == nil {
			in.errorSemantico("no se puede acceder al campo \""+n.Nombre+"\" de un valor que no es struct", n.Linea, n.Columna)
			return Valor{}, false
		}
		campo, ok := base.Struct.Campos[n.Nombre]
		if !ok {
			in.errorSemantico("el struct \""+base.Struct.NombreTipo+"\" no tiene el campo \""+n.Nombre+"\"", n.Linea, n.Columna)
			return Valor{}, false
		}
		return *campo, true
	case *ast.ExprLlamada:
		return in.evaluarLlamada(n, ent)
	case *ast.LiteralSlice:
		return in.evaluarLiteralSlice(n, ent)
	case *ast.LiteralStruct:
		return in.evaluarLiteralStruct(n, ent)
	}
	// B1: expresion de un tipo de nodo no contemplado -> error con posicion,
	// en vez de un {} silencioso que se propaga como si fuera un valor valido.
	l, c := nodo.Pos()
	in.errorSemantico("expresión no soportada", l, c)
	return Valor{}, false
}

func (in *Interprete) evaluarBinaria(op string, a, b Valor, linea, col int) (Valor, bool) {
	var r Valor
	var err error
	switch op {
	case "+":
		r, err = Suma(a, b, linea, col)
	case "-":
		r, err = Resta(a, b, linea, col)
	case "*":
		r, err = Multiplicacion(a, b, linea, col)
	case "/":
		r, err = Division(a, b, linea, col)
	case "%":
		r, err = Modulo(a, b, linea, col)
	case "<", "<=", ">", ">=":
		r, err = Relacional(op, a, b)
	case "==":
		r, err = Igualdad(a, b, false)
	case "!=":
		r, err = Igualdad(a, b, true)
	case "&&":
		r, err = Y(a, b)
	case "||":
		r, err = O(a, b)
	default:
		err = fmt.Errorf("operador desconocido %s", op)
	}
	if err != nil {
		in.errorSemantico(err.Error(), linea, col)
		return Valor{}, false
	}
	return r, true
}

// evaluarLlamada resuelve el callee en 3 casos: funcion nativa (por
// nombre, ver nativas.go), funcion libre (identificador registrado en
// global.Funciones) o metodo (Base.Metodo(), busca el metodo por el tipo
// de struct del receptor evaluado). Los argumentos se evaluan antes de
// saber cual de los 3 casos aplica, como en cualquier evaluacion de call.
func (in *Interprete) evaluarLlamada(n *ast.ExprLlamada, ent *Entorno) (Valor, bool) {
	args := make([]Valor, 0, len(n.Argumentos))
	for _, a := range n.Argumentos {
		v, ok := in.evaluar(a, ent)
		if !ok {
			return Valor{}, false
		}
		args = append(args, v)
	}

	switch callee := n.Callee.(type) {
	case *ast.Identificador:
		if EsNombreNativa(callee.Nombre) {
			return in.llamarNativa(callee.Nombre, args, n.Linea, n.Columna)
		}
		fn, ok := in.global.Funciones[callee.Nombre]
		if !ok {
			in.errorSemantico("la función \""+callee.Nombre+"\" no está declarada", n.Linea, n.Columna)
			return Valor{}, false
		}
		return in.invocarFuncion(fn, args, nil, n.Linea, n.Columna)

	case *ast.ExprCampo:
		recVal, ok := in.evaluar(callee.Base, ent)
		if !ok {
			return Valor{}, false
		}
		if recVal.Tipo.Base != TStruct || recVal.Struct == nil {
			in.errorSemantico("no se puede llamar \""+callee.Nombre+"\" sobre un valor que no es struct", n.Linea, n.Columna)
			return Valor{}, false
		}
		metodos := in.global.Metodos[recVal.Struct.NombreTipo]
		fn, ok := metodos[callee.Nombre]
		if !ok {
			in.errorSemantico("el struct \""+recVal.Struct.NombreTipo+"\" no tiene el método \""+callee.Nombre+"\"", n.Linea, n.Columna)
			return Valor{}, false
		}
		return in.invocarFuncion(fn, args, &recVal, n.Linea, n.Columna)
	}

	in.errorSemantico("la expresión no es invocable", n.Linea, n.Columna)
	return Valor{}, false
}

// invocarFuncion arma un ambito nuevo con el nombre de la funcion (para
// que el reporte 8.2 agrupe sus variables locales), declara el receptor
// (si es un metodo) y los parametros ya validados/coercionados, ejecuta
// el cuerpo y valida que la senal final sea consistente con el tipo de
// retorno declarado (o su ausencia).
func (in *Interprete) invocarFuncion(fn *ast.DeclFuncion, args []Valor, receptor *Valor, linea, col int) (Valor, bool) {
	in.profundidad++
	defer func() { in.profundidad-- }()
	if in.profundidad > maxProfundidad {
		in.errorSemantico(fmt.Sprintf("recursión demasiado profunda al llamar \"%s\" (posible recursión sin caso base)", fn.Nombre), linea, col)
		return Valor{}, false
	}
	if len(args) != len(fn.Parametros) {
		in.errorSemantico(fmt.Sprintf("\"%s\" espera %d argumento(s), se recibieron %d", fn.Nombre, len(fn.Parametros), len(args)), linea, col)
		return Valor{}, false
	}

	entFn := NuevoEntorno(fn.Nombre, in.entGlobal)
	if receptor != nil && fn.ReceptorNombre != "" {
		rec := *receptor
		if !fn.ReceptorPuntero {
			// Receptor por valor (7.1): el metodo trabaja sobre una COPIA del
			// struct; mutar sus campos NO se refleja en el llamador. Con
			// receptor por puntero, en cambio, se comparte el StructVal.
			rec = ClonarPorValor(rec)
		}
		entFn.Declarar(fn.ReceptorNombre, rec)
	}
	for i, p := range fn.Parametros {
		tipoParam := in.resolverTipoAST(p.Tipo, fn.Linea, fn.Columna)
		if !tipoCompatible(tipoParam, args[i]) {
			in.errorSemantico(fmt.Sprintf("el parámetro \"%s\" espera %s, se recibió %s", p.Nombre, tipoParam, args[i].Tipo), linea, col)
			return Valor{}, false
		}
		entFn.Declarar(p.Nombre, coercionar(tipoParam, args[i]))
	}

	// Un break/continue "perdido" en el cuerpo de fn no debe heredar el
	// contexto de ciclo/switch del llamador (p.ej. una funcion invocada
	// desde dentro de un for no vuelve validos sus propios break sueltos).
	loopPrevio, switchPrevio := in.profLoop, in.profSwitch
	in.profLoop, in.profSwitch = 0, 0
	senal := in.ejecutarBloque(fn.Cuerpo, entFn)
	in.profLoop, in.profSwitch = loopPrevio, switchPrevio

	if fn.TipoRetorno == nil {
		// Sin tipo de retorno declarado: un "return <valor>" es un error
		// semantico (no puede devolver nada), en vez de descartarse en silencio.
		if senal.Tipo == SenalReturn && senal.TieneValor {
			in.errorSemantico("\""+fn.Nombre+"\" no declara tipo de retorno y no puede retornar un valor", fn.Linea, fn.Columna)
		}
		return Valor{Tipo: TipoVoid()}, true
	}
	tipoRet := in.resolverTipoAST(*fn.TipoRetorno, fn.Linea, fn.Columna)
	if senal.Tipo != SenalReturn {
		in.errorSemantico("la función \""+fn.Nombre+"\" debe retornar un valor de tipo "+tipoRet.String(), fn.Linea, fn.Columna)
		return ValorPorDefecto(tipoRet), false
	}
	if !senal.TieneValor {
		// "return;" en una funcion que declara tipo de retorno
		in.errorSemantico("la función \""+fn.Nombre+"\" declara retorno "+tipoRet.String()+" pero tiene un \"return\" sin valor", fn.Linea, fn.Columna)
		return ValorPorDefecto(tipoRet), false
	}
	if !tipoCompatible(tipoRet, senal.Valor) {
		in.errorSemantico(fmt.Sprintf("\"%s\" debe retornar %s, se retornó %s", fn.Nombre, tipoRet, senal.Valor.Tipo), fn.Linea, fn.Columna)
		return Valor{}, false
	}
	return coercionar(tipoRet, senal.Valor), true
}

// evaluarLiteralSlice arma un []T o un [][]T (slice de slices) segun si
// TipoElem.EsSlice: en el caso 2D, cada fila del literal se evalua y
// tipa como un []T independiente antes de envolverlas en el slice externo.
func (in *Interprete) evaluarLiteralSlice(n *ast.LiteralSlice, ent *Entorno) (Valor, bool) {
	if n.TipoElem.EsSlice {
		tipoFila := in.resolverTipoAST(*n.TipoElem.Elemento, n.Linea, n.Columna)
		var filas []Valor
		for _, fila := range n.Filas {
			elems, ok := in.evaluarListaConTipo(fila, tipoFila, ent, n.Linea, n.Columna)
			if !ok {
				return Valor{}, false
			}
			filas = append(filas, NuevoSlice(tipoFila, elems))
		}
		tipoDeFila := in.resolverTipoAST(n.TipoElem, n.Linea, n.Columna)
		return NuevoSlice(tipoDeFila, filas), true
	}

	tipoElem := in.resolverTipoAST(n.TipoElem, n.Linea, n.Columna)
	var expresiones []ast.Nodo
	if len(n.Filas) > 0 {
		expresiones = n.Filas[0]
	}
	elems, ok := in.evaluarListaConTipo(expresiones, tipoElem, ent, n.Linea, n.Columna)
	if !ok {
		return Valor{}, false
	}
	return NuevoSlice(tipoElem, elems), true
}

func (in *Interprete) evaluarListaConTipo(nodos []ast.Nodo, tipoElem Tipo, ent *Entorno, linea, col int) ([]Valor, bool) {
	elems := make([]Valor, 0, len(nodos))
	for _, e := range nodos {
		v, ok := in.evaluar(e, ent)
		if !ok {
			return nil, false
		}
		if !tipoCompatible(tipoElem, v) {
			in.errorSemantico(fmt.Sprintf("elemento de tipo %s no es compatible con []%s", v.Tipo, tipoElem), linea, col)
			return nil, false
		}
		elems = append(elems, coercionar(tipoElem, v))
	}
	return elems, true
}

// evaluarLiteralStruct instancia un struct: primero llena todos los
// campos con ValorPorDefecto segun la declaracion (para que un literal
// parcial "Persona{Nombre: "Ana"}" deje el resto en su default), y luego
// sobreescribe con los campos que el literal especifico trae.
func (in *Interprete) evaluarLiteralStruct(n *ast.LiteralStruct, ent *Entorno) (Valor, bool) {
	decl, ok := in.global.Structs[n.NombreStruct]
	if !ok {
		in.errorSemantico("el struct \""+n.NombreStruct+"\" no fue definido", n.Linea, n.Columna)
		return Valor{}, false
	}

	campos := make(map[string]*Valor)
	var orden []string
	for _, cs := range decl.Campos {
		t := in.resolverTipoAST(cs.Tipo, n.Linea, n.Columna)
		v := ValorPorDefecto(t)
		campos[cs.Nombre] = &v
		orden = append(orden, cs.Nombre)
	}

	vistos := make(map[string]bool)
	for _, cv := range n.Campos {
		// B2: el mismo campo dos veces en el literal (Persona{Nombre:"a",
		// Nombre:"b"}) es error, en vez de que el ultimo gane en silencio.
		if vistos[cv.Nombre] {
			in.errorSemantico("el campo \""+cv.Nombre+"\" aparece repetido en el literal de \""+n.NombreStruct+"\"", n.Linea, n.Columna)
			return Valor{}, false
		}
		vistos[cv.Nombre] = true
		campoDecl := buscarCampoStruct(decl, cv.Nombre)
		if campoDecl == nil {
			in.errorSemantico("el struct \""+n.NombreStruct+"\" no tiene el campo \""+cv.Nombre+"\"", n.Linea, n.Columna)
			return Valor{}, false
		}
		v, ok := in.evaluar(cv.Valor, ent)
		if !ok {
			return Valor{}, false
		}
		tipoCampo := in.resolverTipoAST(campoDecl.Tipo, n.Linea, n.Columna)
		if !tipoCompatible(tipoCampo, v) {
			in.errorSemantico(fmt.Sprintf("el campo \"%s\" espera %s, se recibió %s", cv.Nombre, tipoCampo, v.Tipo), n.Linea, n.Columna)
			return Valor{}, false
		}
		*campos[cv.Nombre] = coercionar(tipoCampo, v)
	}

	return NuevaInstanciaStruct(n.NombreStruct, orden, campos), true
}

func buscarCampoStruct(decl *ast.DeclStruct, nombre string) *ast.CampoStruct {
	for i := range decl.Campos {
		if decl.Campos[i].Nombre == nombre {
			return &decl.Campos[i]
		}
	}
	return nil
}
