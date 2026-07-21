// Code generated from VLangCherry.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // VLangCherry
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by VLangCherryParser.
type VLangCherryVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by VLangCherryParser#programa.
	VisitPrograma(ctx *ProgramaContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#declaracionGlobal.
	VisitDeclaracionGlobal(ctx *DeclaracionGlobalContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#declaracionStruct.
	VisitDeclaracionStruct(ctx *DeclaracionStructContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#campoStruct.
	VisitCampoStruct(ctx *CampoStructContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#declaracionFuncion.
	VisitDeclaracionFuncion(ctx *DeclaracionFuncionContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#receptor.
	VisitReceptor(ctx *ReceptorContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#listaParametros.
	VisitListaParametros(ctx *ListaParametrosContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#parametro.
	VisitParametro(ctx *ParametroContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#tipo.
	VisitTipo(ctx *TipoContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#tipoPrimitivo.
	VisitTipoPrimitivo(ctx *TipoPrimitivoContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#tipoSlice.
	VisitTipoSlice(ctx *TipoSliceContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#bloque.
	VisitBloque(ctx *BloqueContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#sentDeclaracion.
	VisitSentDeclaracion(ctx *SentDeclaracionContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#sentAsignacion.
	VisitSentAsignacion(ctx *SentAsignacionContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#sentIncDec.
	VisitSentIncDec(ctx *SentIncDecContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#sentExpresion.
	VisitSentExpresion(ctx *SentExpresionContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#sentIf.
	VisitSentIf(ctx *SentIfContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#sentSwitch.
	VisitSentSwitch(ctx *SentSwitchContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#sentFor.
	VisitSentFor(ctx *SentForContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#sentBreak.
	VisitSentBreak(ctx *SentBreakContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#sentContinue.
	VisitSentContinue(ctx *SentContinueContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#sentReturn.
	VisitSentReturn(ctx *SentReturnContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#sentBloque.
	VisitSentBloque(ctx *SentBloqueContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#declTipada.
	VisitDeclTipada(ctx *DeclTipadaContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#declInferida.
	VisitDeclInferida(ctx *DeclInferidaContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#asignacion.
	VisitAsignacion(ctx *AsignacionContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#lugarCampo.
	VisitLugarCampo(ctx *LugarCampoContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#lugarIndice.
	VisitLugarIndice(ctx *LugarIndiceContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#lugarId.
	VisitLugarId(ctx *LugarIdContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#incrementoDecremento.
	VisitIncrementoDecremento(ctx *IncrementoDecrementoContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#sentenciaIf.
	VisitSentenciaIf(ctx *SentenciaIfContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#sentenciaSwitch.
	VisitSentenciaSwitch(ctx *SentenciaSwitchContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#casoSwitch.
	VisitCasoSwitch(ctx *CasoSwitchContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#defaultSwitch.
	VisitDefaultSwitch(ctx *DefaultSwitchContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#forCondicion.
	VisitForCondicion(ctx *ForCondicionContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#forClasico.
	VisitForClasico(ctx *ForClasicoContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#forRango.
	VisitForRango(ctx *ForRangoContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#forInit.
	VisitForInit(ctx *ForInitContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#forActualizacion.
	VisitForActualizacion(ctx *ForActualizacionContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprIdentificador.
	VisitExprIdentificador(ctx *ExprIdentificadorContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprAditiva.
	VisitExprAditiva(ctx *ExprAditivaContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprStructLit.
	VisitExprStructLit(ctx *ExprStructLitContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprRelacional.
	VisitExprRelacional(ctx *ExprRelacionalContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprIndice.
	VisitExprIndice(ctx *ExprIndiceContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprParentesis.
	VisitExprParentesis(ctx *ExprParentesisContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprSliceLit.
	VisitExprSliceLit(ctx *ExprSliceLitContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprLiteral.
	VisitExprLiteral(ctx *ExprLiteralContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprOr.
	VisitExprOr(ctx *ExprOrContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprLlamada.
	VisitExprLlamada(ctx *ExprLlamadaContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprCampo.
	VisitExprCampo(ctx *ExprCampoContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprIgualdad.
	VisitExprIgualdad(ctx *ExprIgualdadContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprAnd.
	VisitExprAnd(ctx *ExprAndContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprUnario.
	VisitExprUnario(ctx *ExprUnarioContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#exprMultiplicativa.
	VisitExprMultiplicativa(ctx *ExprMultiplicativaContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#listaArgumentos.
	VisitListaArgumentos(ctx *ListaArgumentosContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#literalSlice.
	VisitLiteralSlice(ctx *LiteralSliceContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#filaSlice.
	VisitFilaSlice(ctx *FilaSliceContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#literalStruct.
	VisitLiteralStruct(ctx *LiteralStructContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#campoValor.
	VisitCampoValor(ctx *CampoValorContext) interface{}

	// Visit a parse tree produced by VLangCherryParser#literal.
	VisitLiteral(ctx *LiteralContext) interface{}
}
