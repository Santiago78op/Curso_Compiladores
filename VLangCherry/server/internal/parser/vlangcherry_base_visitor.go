// Code generated from VLangCherry.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // VLangCherry
import "github.com/antlr4-go/antlr/v4"

type BaseVLangCherryVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseVLangCherryVisitor) VisitPrograma(ctx *ProgramaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitDeclaracionGlobal(ctx *DeclaracionGlobalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitDeclaracionStruct(ctx *DeclaracionStructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitCampoStruct(ctx *CampoStructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitDeclaracionFuncion(ctx *DeclaracionFuncionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitReceptor(ctx *ReceptorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitListaParametros(ctx *ListaParametrosContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitParametro(ctx *ParametroContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitTipo(ctx *TipoContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitTipoPrimitivo(ctx *TipoPrimitivoContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitTipoSlice(ctx *TipoSliceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitBloque(ctx *BloqueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitSentDeclaracion(ctx *SentDeclaracionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitSentAsignacion(ctx *SentAsignacionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitSentIncDec(ctx *SentIncDecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitSentExpresion(ctx *SentExpresionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitSentIf(ctx *SentIfContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitSentSwitch(ctx *SentSwitchContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitSentFor(ctx *SentForContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitSentBreak(ctx *SentBreakContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitSentContinue(ctx *SentContinueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitSentReturn(ctx *SentReturnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitSentBloque(ctx *SentBloqueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitDeclTipada(ctx *DeclTipadaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitDeclInferida(ctx *DeclInferidaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitAsignacion(ctx *AsignacionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitLugarCampo(ctx *LugarCampoContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitLugarIndice(ctx *LugarIndiceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitLugarId(ctx *LugarIdContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitIncrementoDecremento(ctx *IncrementoDecrementoContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitSentenciaIf(ctx *SentenciaIfContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitSentenciaSwitch(ctx *SentenciaSwitchContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitCasoSwitch(ctx *CasoSwitchContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitDefaultSwitch(ctx *DefaultSwitchContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitForCondicion(ctx *ForCondicionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitForClasico(ctx *ForClasicoContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitForRango(ctx *ForRangoContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitForInit(ctx *ForInitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitForActualizacion(ctx *ForActualizacionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprIdentificador(ctx *ExprIdentificadorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprAditiva(ctx *ExprAditivaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprStructLit(ctx *ExprStructLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprRelacional(ctx *ExprRelacionalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprIndice(ctx *ExprIndiceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprParentesis(ctx *ExprParentesisContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprSliceLit(ctx *ExprSliceLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprLiteral(ctx *ExprLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprOr(ctx *ExprOrContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprLlamada(ctx *ExprLlamadaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprCampo(ctx *ExprCampoContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprIgualdad(ctx *ExprIgualdadContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprAnd(ctx *ExprAndContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprUnario(ctx *ExprUnarioContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitExprMultiplicativa(ctx *ExprMultiplicativaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitListaArgumentos(ctx *ListaArgumentosContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitLiteralSlice(ctx *LiteralSliceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitFilaSlice(ctx *FilaSliceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitLiteralStruct(ctx *LiteralStructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitCampoValor(ctx *CampoValorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseVLangCherryVisitor) VisitLiteral(ctx *LiteralContext) interface{} {
	return v.VisitChildren(ctx)
}
