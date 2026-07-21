// Package traductor convierte el parse tree que genera ANTLR en el AST
// propio de internal/ast. No usa el patron Visitor generado: como cada
// alternativa etiquetada de la gramatica produce un tipo de contexto Go
// concreto, alcanza con un type-switch directo (mas corto que implementar
// las ~60 firmas de BaseVLangCherryVisitor).
package traductor

import (
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"vlangcherry/internal/ast"
	"vlangcherry/internal/parser"
)

func pos(ctx antlr.ParserRuleContext) (int, int) {
	t := ctx.GetStart()
	return t.GetLine(), t.GetColumn() + 1
}

// Traducir es el punto de entrada: parse tree -> ast.Programa.
func Traducir(ctx parser.IProgramaContext) *ast.Programa {
	p := &ast.Programa{}
	l, c := pos(ctx)
	p.Linea, p.Columna = l, c
	for _, dg := range ctx.AllDeclaracionGlobal() {
		dgc := dg.(*parser.DeclaracionGlobalContext)
		switch {
		case dgc.DeclaracionStruct() != nil:
			p.Structs = append(p.Structs, traducirStruct(dgc.DeclaracionStruct().(*parser.DeclaracionStructContext)))
		case dgc.DeclaracionFuncion() != nil:
			p.Funciones = append(p.Funciones, traducirFuncion(dgc.DeclaracionFuncion().(*parser.DeclaracionFuncionContext)))
		case dgc.DeclaracionVariable() != nil:
			p.Globales = append(p.Globales, traducirDeclVariable(dgc.DeclaracionVariable()))
		}
	}
	return p
}

func traducirStruct(ctx *parser.DeclaracionStructContext) *ast.DeclStruct {
	l, c := pos(ctx)
	d := &ast.DeclStruct{Nombre: ctx.ID().GetText()}
	d.Linea, d.Columna = l, c
	for _, cs := range ctx.AllCampoStruct() {
		csc := cs.(*parser.CampoStructContext)
		d.Campos = append(d.Campos, ast.CampoStruct{
			Nombre: csc.ID().GetText(),
			Tipo:   traducirTipo(csc.Tipo()),
		})
	}
	return d
}

func traducirFuncion(ctx *parser.DeclaracionFuncionContext) *ast.DeclFuncion {
	l, c := pos(ctx)
	f := &ast.DeclFuncion{Nombre: ctx.ID().GetText()}
	f.Linea, f.Columna = l, c

	if r := ctx.Receptor(); r != nil {
		rc := r.(*parser.ReceptorContext)
		ids := rc.AllID()
		f.ReceptorTipo = ids[0].GetText()
		f.ReceptorNombre = ""
		if len(ids) == 2 {
			// (p *Persona): ID0=p (nombre), ID1=Persona (tipo)
			f.ReceptorNombre = ids[0].GetText()
			f.ReceptorTipo = ids[1].GetText()
		}
		f.ReceptorPuntero = rc.POR() != nil
	}

	if lp := ctx.ListaParametros(); lp != nil {
		lpc := lp.(*parser.ListaParametrosContext)
		for _, pr := range lpc.AllParametro() {
			prc := pr.(*parser.ParametroContext)
			f.Parametros = append(f.Parametros, ast.Parametro{
				Nombre: prc.ID().GetText(),
				Tipo:   traducirTipo(prc.Tipo()),
			})
		}
	}

	if t := ctx.Tipo(); t != nil {
		tt := traducirTipo(t)
		f.TipoRetorno = &tt
	}

	f.Cuerpo = traducirBloque(ctx.Bloque())
	return f
}

func traducirTipo(ctx parser.ITipoContext) ast.TipoAST {
	tc := ctx.(*parser.TipoContext)
	if ts := tc.TipoSlice(); ts != nil {
		tsc := ts.(*parser.TipoSliceContext)
		el := traducirTipo(tsc.Tipo())
		return ast.TipoAST{EsSlice: true, Elemento: &el}
	}
	if tp := tc.TipoPrimitivo(); tp != nil {
		return ast.TipoAST{NombreBase: tp.GetText()}
	}
	return ast.TipoAST{NombreBase: tc.ID().GetText()}
}

func traducirBloque(ctx parser.IBloqueContext) *ast.Bloque {
	bc := ctx.(*parser.BloqueContext)
	l, c := pos(bc)
	b := &ast.Bloque{}
	b.Linea, b.Columna = l, c
	for _, s := range bc.AllSentencia() {
		b.Sentencias = append(b.Sentencias, traducirSentencia(s))
	}
	return b
}

func traducirSentencia(ctx parser.ISentenciaContext) ast.Nodo {
	switch sc := ctx.(type) {
	case *parser.SentDeclaracionContext:
		return traducirDeclVariable(sc.DeclaracionVariable())
	case *parser.SentAsignacionContext:
		return traducirAsignacion(sc.Asignacion().(*parser.AsignacionContext))
	case *parser.SentIncDecContext:
		return traducirIncDec(sc.IncrementoDecremento().(*parser.IncrementoDecrementoContext))
	case *parser.SentExpresionContext:
		l, c := pos(sc)
		n := &ast.ExpresionSentencia{Expr: traducirExpr(sc.Expr())}
		n.Linea, n.Columna = l, c
		return n
	case *parser.SentIfContext:
		return traducirIf(sc.SentenciaIf().(*parser.SentenciaIfContext))
	case *parser.SentSwitchContext:
		return traducirSwitch(sc.SentenciaSwitch().(*parser.SentenciaSwitchContext))
	case *parser.SentForContext:
		return traducirFor(sc.SentenciaFor())
	case *parser.SentBreakContext:
		l, c := pos(sc)
		n := &ast.SentenciaBreak{}
		n.Linea, n.Columna = l, c
		return n
	case *parser.SentContinueContext:
		l, c := pos(sc)
		n := &ast.SentenciaContinue{}
		n.Linea, n.Columna = l, c
		return n
	case *parser.SentReturnContext:
		l, c := pos(sc)
		n := &ast.SentenciaReturn{}
		n.Linea, n.Columna = l, c
		if e := sc.Expr(); e != nil {
			n.Valor = traducirExpr(e)
		}
		return n
	case *parser.SentBloqueContext:
		return traducirBloque(sc.Bloque())
	}
	return nil
}

func traducirDeclVariable(ctx parser.IDeclaracionVariableContext) *ast.DeclVariable {
	switch dc := ctx.(type) {
	case *parser.DeclTipadaContext:
		l, c := pos(dc)
		n := &ast.DeclVariable{Nombre: dc.ID().GetText()}
		n.Linea, n.Columna = l, c
		tt := traducirTipo(dc.Tipo())
		n.TipoVar = &tt
		if e := dc.Expr(); e != nil {
			n.Valor = traducirExpr(e)
		}
		return n
	case *parser.DeclInferidaContext:
		l, c := pos(dc)
		n := &ast.DeclVariable{Nombre: dc.ID().GetText(), Inferido: true}
		n.Linea, n.Columna = l, c
		n.Valor = traducirExpr(dc.Expr())
		return n
	}
	return nil
}

func traducirLugar(ctx parser.ILugarContext) ast.Nodo {
	switch lc := ctx.(type) {
	case *parser.LugarIdContext:
		l, c := pos(lc)
		n := &ast.Identificador{Nombre: lc.ID().GetText()}
		n.Linea, n.Columna = l, c
		return n
	case *parser.LugarIndiceContext:
		l, c := pos(lc)
		n := &ast.ExprIndice{Base: traducirLugar(lc.Lugar()), Indice: traducirExpr(lc.Expr())}
		n.Linea, n.Columna = l, c
		return n
	case *parser.LugarCampoContext:
		l, c := pos(lc)
		n := &ast.ExprCampo{Base: traducirLugar(lc.Lugar()), Nombre: lc.ID().GetText()}
		n.Linea, n.Columna = l, c
		return n
	}
	return nil
}

func traducirAsignacion(ctx *parser.AsignacionContext) *ast.Asignacion {
	l, c := pos(ctx)
	n := &ast.Asignacion{
		Lugar:    traducirLugar(ctx.Lugar()),
		Operador: ctx.GetOp().GetText(),
		Valor:    traducirExpr(ctx.Expr()),
	}
	n.Linea, n.Columna = l, c
	return n
}

func traducirIncDec(ctx *parser.IncrementoDecrementoContext) *ast.IncrementoDecremento {
	l, c := pos(ctx)
	n := &ast.IncrementoDecremento{
		Lugar:    traducirLugar(ctx.Lugar()),
		Operador: ctx.GetOp().GetText(),
	}
	n.Linea, n.Columna = l, c
	return n
}

func traducirIf(ctx *parser.SentenciaIfContext) *ast.SentenciaIf {
	l, c := pos(ctx)
	n := &ast.SentenciaIf{}
	n.Linea, n.Columna = l, c

	exprs := ctx.AllExpr()
	bloques := ctx.AllBloque()
	// Ramas if/else-if: una rama por cada expr (el bloque con mismo indice).
	for i, e := range exprs {
		n.Ramas = append(n.Ramas, ast.RamaIf{
			Condicion: traducirExpr(e),
			Cuerpo:    traducirBloque(bloques[i]),
		})
	}
	// Si hay un bloque extra (sin expr asociada) es el else final.
	if len(bloques) > len(exprs) {
		n.Else = traducirBloque(bloques[len(bloques)-1])
	}
	return n
}

func traducirSwitch(ctx *parser.SentenciaSwitchContext) *ast.SentenciaSwitch {
	l, c := pos(ctx)
	n := &ast.SentenciaSwitch{Expr: traducirExpr(ctx.Expr())}
	n.Linea, n.Columna = l, c

	for _, cs := range ctx.AllCasoSwitch() {
		csc := cs.(*parser.CasoSwitchContext)
		caso := ast.CasoSwitch{Valor: traducirExpr(csc.Expr())}
		for _, s := range csc.AllSentencia() {
			caso.Sentencias = append(caso.Sentencias, traducirSentencia(s))
		}
		n.Casos = append(n.Casos, caso)
	}
	if def := ctx.DefaultSwitch(); def != nil {
		defc := def.(*parser.DefaultSwitchContext)
		caso := ast.CasoSwitch{Valor: nil}
		for _, s := range defc.AllSentencia() {
			caso.Sentencias = append(caso.Sentencias, traducirSentencia(s))
		}
		n.Casos = append(n.Casos, caso)
	}
	return n
}

func traducirFor(ctx parser.ISentenciaForContext) *ast.SentenciaFor {
	switch fc := ctx.(type) {
	case *parser.ForCondicionContext:
		l, c := pos(fc)
		n := &ast.SentenciaFor{Forma: "condicion", Condicion: traducirExpr(fc.Expr()), Cuerpo: traducirBloque(fc.Bloque())}
		n.Linea, n.Columna = l, c
		return n
	case *parser.ForClasicoContext:
		l, c := pos(fc)
		n := &ast.SentenciaFor{Forma: "clasico", Cuerpo: traducirBloque(fc.Bloque())}
		n.Linea, n.Columna = l, c
		if init := fc.ForInit(); init != nil {
			n.Init = traducirForInit(init.(*parser.ForInitContext))
		}
		if e := fc.Expr(); e != nil {
			n.Condicion = traducirExpr(e)
		}
		if act := fc.ForActualizacion(); act != nil {
			n.Actualizacion = traducirForActualizacion(act.(*parser.ForActualizacionContext))
		}
		return n
	case *parser.ForRangoContext:
		l, c := pos(fc)
		ids := fc.AllID()
		n := &ast.SentenciaFor{
			Forma:     "rango",
			VarIndice: ids[0].GetText(),
			VarValor:  ids[1].GetText(),
			Iterable:  traducirExpr(fc.Expr()),
			Cuerpo:    traducirBloque(fc.Bloque()),
		}
		n.Linea, n.Columna = l, c
		return n
	}
	return nil
}

func traducirForInit(ctx *parser.ForInitContext) ast.Nodo {
	if dv := ctx.DeclaracionVariable(); dv != nil {
		return traducirDeclVariable(dv)
	}
	return traducirAsignacion(ctx.Asignacion().(*parser.AsignacionContext))
}

func traducirForActualizacion(ctx *parser.ForActualizacionContext) ast.Nodo {
	if a := ctx.Asignacion(); a != nil {
		return traducirAsignacion(a.(*parser.AsignacionContext))
	}
	return traducirIncDec(ctx.IncrementoDecremento().(*parser.IncrementoDecrementoContext))
}

// ---------------- Expresiones ----------------

func traducirExpr(ctx parser.IExprContext) ast.Nodo {
	switch ec := ctx.(type) {
	case *parser.ExprLlamadaContext:
		l, c := pos(ec)
		n := &ast.ExprLlamada{Callee: traducirExpr(ec.Expr())}
		n.Linea, n.Columna = l, c
		if la := ec.ListaArgumentos(); la != nil {
			n.Argumentos = traducirListaArgumentos(la.(*parser.ListaArgumentosContext))
		}
		return n
	case *parser.ExprIndiceContext:
		l, c := pos(ec)
		exprs := ec.AllExpr()
		n := &ast.ExprIndice{Base: traducirExpr(exprs[0]), Indice: traducirExpr(exprs[1])}
		n.Linea, n.Columna = l, c
		return n
	case *parser.ExprCampoContext:
		l, c := pos(ec)
		n := &ast.ExprCampo{Base: traducirExpr(ec.Expr()), Nombre: ec.ID().GetText()}
		n.Linea, n.Columna = l, c
		return n
	case *parser.ExprParentesisContext:
		return traducirExpr(ec.Expr())
	case *parser.ExprUnarioContext:
		l, c := pos(ec)
		n := &ast.ExprUnaria{Operador: ec.GetOp().GetText(), Operando: traducirExpr(ec.Expr())}
		n.Linea, n.Columna = l, c
		return n
	case *parser.ExprMultiplicativaContext:
		return traducirBinaria(ec, ec.AllExpr(), ec.GetOp().GetText())
	case *parser.ExprAditivaContext:
		return traducirBinaria(ec, ec.AllExpr(), ec.GetOp().GetText())
	case *parser.ExprRelacionalContext:
		return traducirBinaria(ec, ec.AllExpr(), ec.GetOp().GetText())
	case *parser.ExprIgualdadContext:
		return traducirBinaria(ec, ec.AllExpr(), ec.GetOp().GetText())
	case *parser.ExprAndContext:
		return traducirBinaria(ec, ec.AllExpr(), "&&")
	case *parser.ExprOrContext:
		return traducirBinaria(ec, ec.AllExpr(), "||")
	case *parser.ExprSliceLitContext:
		return traducirLiteralSlice(ec.LiteralSlice().(*parser.LiteralSliceContext))
	case *parser.ExprStructLitContext:
		return traducirLiteralStruct(ec.LiteralStruct().(*parser.LiteralStructContext))
	case *parser.ExprLiteralContext:
		return traducirLiteral(ec.Literal().(*parser.LiteralContext))
	case *parser.ExprIdentificadorContext:
		l, c := pos(ec)
		n := &ast.Identificador{Nombre: ec.ID().GetText()}
		n.Linea, n.Columna = l, c
		return n
	}
	return nil
}

func traducirBinaria(ctx antlr.ParserRuleContext, exprs []parser.IExprContext, op string) *ast.ExprBinaria {
	l, c := pos(ctx)
	n := &ast.ExprBinaria{Operador: op, Izquierda: traducirExpr(exprs[0]), Derecha: traducirExpr(exprs[1])}
	n.Linea, n.Columna = l, c
	return n
}

func traducirListaArgumentos(ctx *parser.ListaArgumentosContext) []ast.Nodo {
	var args []ast.Nodo
	for _, e := range ctx.AllExpr() {
		args = append(args, traducirExpr(e))
	}
	return args
}

func traducirLiteralSlice(ctx *parser.LiteralSliceContext) *ast.LiteralSlice {
	l, c := pos(ctx)
	n := &ast.LiteralSlice{TipoElem: traducirTipo(ctx.Tipo())}
	n.Linea, n.Columna = l, c

	if la := ctx.ListaArgumentos(); la != nil {
		n.Filas = [][]ast.Nodo{traducirListaArgumentos(la.(*parser.ListaArgumentosContext))}
		return n
	}
	for _, fs := range ctx.AllFilaSlice() {
		fsc := fs.(*parser.FilaSliceContext)
		var fila []ast.Nodo
		if la := fsc.ListaArgumentos(); la != nil {
			fila = traducirListaArgumentos(la.(*parser.ListaArgumentosContext))
		}
		n.Filas = append(n.Filas, fila)
	}
	return n
}

func traducirLiteralStruct(ctx *parser.LiteralStructContext) *ast.LiteralStruct {
	l, c := pos(ctx)
	n := &ast.LiteralStruct{NombreStruct: ctx.ID().GetText()}
	n.Linea, n.Columna = l, c
	for _, cv := range ctx.AllCampoValor() {
		cvc := cv.(*parser.CampoValorContext)
		n.Campos = append(n.Campos, ast.CampoValorLiteral{
			Nombre: cvc.ID().GetText(),
			Valor:  traducirExpr(cvc.Expr()),
		})
	}
	return n
}

func traducirLiteral(ctx *parser.LiteralContext) ast.Nodo {
	l, c := pos(ctx)
	switch {
	case ctx.ENTERO() != nil:
		v, _ := strconv.ParseInt(ctx.ENTERO().GetText(), 10, 64)
		n := &ast.LiteralEntero{Valor: v}
		n.Linea, n.Columna = l, c
		return n
	case ctx.DECIMAL() != nil:
		v, _ := strconv.ParseFloat(ctx.DECIMAL().GetText(), 64)
		n := &ast.LiteralDecimal{Valor: v}
		n.Linea, n.Columna = l, c
		return n
	case ctx.CADENA() != nil:
		n := &ast.LiteralCadena{Valor: desescaparCadena(ctx.CADENA().GetText())}
		n.Linea, n.Columna = l, c
		return n
	case ctx.RUNE() != nil:
		n := &ast.LiteralRune{Valor: desescaparRune(ctx.RUNE().GetText())}
		n.Linea, n.Columna = l, c
		return n
	case ctx.TRUE() != nil:
		n := &ast.LiteralBool{Valor: true}
		n.Linea, n.Columna = l, c
		return n
	case ctx.FALSE() != nil:
		n := &ast.LiteralBool{Valor: false}
		n.Linea, n.Columna = l, c
		return n
	case ctx.NIL() != nil:
		n := &ast.LiteralNil{}
		n.Linea, n.Columna = l, c
		return n
	}
	return nil
}

func desescaparCadena(lexema string) string {
	interior := lexema[1 : len(lexema)-1]
	return aplicarEscapes(interior)
}

func desescaparRune(lexema string) rune {
	interior := lexema[1 : len(lexema)-1]
	texto := aplicarEscapes(interior)
	for _, r := range texto {
		return r
	}
	return 0
}

func aplicarEscapes(s string) string {
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			i++
			switch s[i] {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case 'r':
				sb.WriteByte('\r')
			case '"':
				sb.WriteByte('"')
			case '\'':
				sb.WriteByte('\'')
			case '\\':
				sb.WriteByte('\\')
			default:
				sb.WriteByte(s[i])
			}
			continue
		}
		sb.WriteByte(s[i])
	}
	return sb.String()
}
