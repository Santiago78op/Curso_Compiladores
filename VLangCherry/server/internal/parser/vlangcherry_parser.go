// Code generated from VLangCherry.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // VLangCherry
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type VLangCherryParser struct {
	*antlr.BaseParser
}

var VLangCherryParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func vlangcherryParserInit() {
	staticData := &VLangCherryParserStaticData
	staticData.LiteralNames = []string{
		"", "'mut'", "'struct'", "'func'", "'if'", "'else'", "'switch'", "'case'",
		"'default'", "'for'", "'in'", "'break'", "'continue'", "'return'", "'nil'",
		"'true'", "'false'", "'int'", "'float64'", "'string'", "'bool'", "'rune'",
		"'+='", "'-='", "'++'", "'--'", "':='", "'=='", "'!='", "'<='", "'>='",
		"'&&'", "'||'", "'<'", "'>'", "'+'", "'-'", "'*'", "'/'", "'%'", "'!'",
		"'='", "'('", "')'", "'['", "']'", "'{'", "'}'", "','", "'.'", "':'",
		"';'",
	}
	staticData.SymbolicNames = []string{
		"", "MUT", "STRUCT", "FUNC", "IF", "ELSE", "SWITCH", "CASE", "DEFAULT",
		"FOR", "IN", "BREAK", "CONTINUE", "RETURN", "NIL", "TRUE", "FALSE",
		"TIPO_INT", "TIPO_FLOAT", "TIPO_STRING", "TIPO_BOOL", "TIPO_RUNE", "MASIGUAL",
		"MENOSIGUAL", "INCREMENTO", "DECREMENTO", "ASIGNAINFERIDA", "IGUAL",
		"DIFERENTE", "MENORIGUAL", "MAYORIGUAL", "Y", "O", "MENOR", "MAYOR",
		"MAS", "MENOS", "POR", "DIV", "MODULO", "NOT", "ASIGNA", "PARIZQ", "PARDER",
		"CORIZQ", "CORDER", "LLAIZQ", "LLADER", "COMA", "PUNTO", "DOSPUNTOS",
		"PUNTOCOMA", "ID", "ENTERO", "DECIMAL", "CADENA", "RUNE", "COMENTARIO_LINEA",
		"COMENTARIO_BLOQUE", "WS",
	}
	staticData.RuleNames = []string{
		"programa", "declaracionGlobal", "declaracionStruct", "campoStruct",
		"declaracionFuncion", "receptor", "listaParametros", "parametro", "tipo",
		"tipoPrimitivo", "tipoSlice", "bloque", "sentencia", "declaracionVariable",
		"asignacion", "lugar", "incrementoDecremento", "sentenciaIf", "sentenciaSwitch",
		"casoSwitch", "defaultSwitch", "sentenciaFor", "forInit", "forActualizacion",
		"expr", "listaArgumentos", "literalSlice", "filaSlice", "literalStruct",
		"campoValor", "literal",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 59, 417, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 1, 0, 5,
		0, 64, 8, 0, 10, 0, 12, 0, 67, 9, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 3, 1,
		74, 8, 1, 1, 2, 1, 2, 1, 2, 1, 2, 4, 2, 80, 8, 2, 11, 2, 12, 2, 81, 1,
		2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 3, 4, 91, 8, 4, 1, 4, 1, 4, 1, 4,
		3, 4, 96, 8, 4, 1, 4, 1, 4, 3, 4, 100, 8, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1,
		5, 3, 5, 107, 8, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 5, 6, 115, 8, 6,
		10, 6, 12, 6, 118, 9, 6, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 3, 8, 126,
		8, 8, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 10, 1, 11, 1, 11, 5, 11, 136,
		8, 11, 10, 11, 12, 11, 139, 9, 11, 1, 11, 1, 11, 1, 12, 1, 12, 1, 12, 1,
		12, 1, 12, 3, 12, 148, 8, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 3, 12,
		155, 8, 12, 1, 12, 1, 12, 3, 12, 159, 8, 12, 1, 12, 1, 12, 3, 12, 163,
		8, 12, 1, 12, 3, 12, 166, 8, 12, 1, 12, 3, 12, 169, 8, 12, 1, 13, 3, 13,
		172, 8, 13, 1, 13, 1, 13, 1, 13, 1, 13, 3, 13, 178, 8, 13, 1, 13, 3, 13,
		181, 8, 13, 1, 13, 3, 13, 184, 8, 13, 1, 13, 1, 13, 1, 13, 1, 13, 3, 13,
		190, 8, 13, 3, 13, 192, 8, 13, 1, 14, 1, 14, 1, 14, 1, 14, 3, 14, 198,
		8, 14, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1,
		15, 1, 15, 5, 15, 211, 8, 15, 10, 15, 12, 15, 214, 9, 15, 1, 16, 1, 16,
		1, 16, 3, 16, 219, 8, 16, 1, 17, 1, 17, 1, 17, 1, 17, 1, 17, 1, 17, 1,
		17, 1, 17, 5, 17, 229, 8, 17, 10, 17, 12, 17, 232, 9, 17, 1, 17, 1, 17,
		3, 17, 236, 8, 17, 1, 18, 1, 18, 1, 18, 1, 18, 5, 18, 242, 8, 18, 10, 18,
		12, 18, 245, 9, 18, 1, 18, 3, 18, 248, 8, 18, 1, 18, 1, 18, 1, 19, 1, 19,
		1, 19, 1, 19, 5, 19, 256, 8, 19, 10, 19, 12, 19, 259, 9, 19, 1, 20, 1,
		20, 1, 20, 5, 20, 264, 8, 20, 10, 20, 12, 20, 267, 9, 20, 1, 21, 1, 21,
		1, 21, 1, 21, 1, 21, 1, 21, 3, 21, 275, 8, 21, 1, 21, 1, 21, 3, 21, 279,
		8, 21, 1, 21, 1, 21, 3, 21, 283, 8, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1,
		21, 1, 21, 1, 21, 1, 21, 1, 21, 3, 21, 294, 8, 21, 1, 22, 1, 22, 3, 22,
		298, 8, 22, 1, 23, 1, 23, 3, 23, 302, 8, 23, 1, 24, 1, 24, 1, 24, 1, 24,
		1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 3, 24, 315, 8, 24, 1,
		24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24,
		1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 3,
		24, 338, 8, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24,
		1, 24, 5, 24, 349, 8, 24, 10, 24, 12, 24, 352, 9, 24, 1, 25, 1, 25, 1,
		25, 5, 25, 357, 8, 25, 10, 25, 12, 25, 360, 9, 25, 1, 26, 1, 26, 1, 26,
		1, 26, 1, 26, 3, 26, 367, 8, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1,
		26, 1, 26, 1, 26, 1, 26, 5, 26, 378, 8, 26, 10, 26, 12, 26, 381, 9, 26,
		1, 26, 3, 26, 384, 8, 26, 1, 26, 1, 26, 3, 26, 388, 8, 26, 1, 27, 1, 27,
		3, 27, 392, 8, 27, 1, 27, 1, 27, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 5,
		28, 401, 8, 28, 10, 28, 12, 28, 404, 9, 28, 1, 28, 3, 28, 407, 8, 28, 1,
		28, 1, 28, 1, 29, 1, 29, 1, 29, 1, 29, 1, 30, 1, 30, 1, 30, 0, 2, 30, 48,
		31, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34,
		36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 0, 9, 1, 0, 17, 21,
		2, 0, 22, 23, 41, 41, 1, 0, 24, 25, 2, 0, 36, 36, 40, 40, 1, 0, 37, 39,
		1, 0, 35, 36, 2, 0, 29, 30, 33, 34, 1, 0, 27, 28, 2, 0, 14, 16, 53, 56,
		458, 0, 65, 1, 0, 0, 0, 2, 73, 1, 0, 0, 0, 4, 75, 1, 0, 0, 0, 6, 85, 1,
		0, 0, 0, 8, 88, 1, 0, 0, 0, 10, 103, 1, 0, 0, 0, 12, 111, 1, 0, 0, 0, 14,
		119, 1, 0, 0, 0, 16, 125, 1, 0, 0, 0, 18, 127, 1, 0, 0, 0, 20, 129, 1,
		0, 0, 0, 22, 133, 1, 0, 0, 0, 24, 168, 1, 0, 0, 0, 26, 191, 1, 0, 0, 0,
		28, 193, 1, 0, 0, 0, 30, 199, 1, 0, 0, 0, 32, 215, 1, 0, 0, 0, 34, 220,
		1, 0, 0, 0, 36, 237, 1, 0, 0, 0, 38, 251, 1, 0, 0, 0, 40, 260, 1, 0, 0,
		0, 42, 293, 1, 0, 0, 0, 44, 297, 1, 0, 0, 0, 46, 301, 1, 0, 0, 0, 48, 314,
		1, 0, 0, 0, 50, 353, 1, 0, 0, 0, 52, 387, 1, 0, 0, 0, 54, 389, 1, 0, 0,
		0, 56, 395, 1, 0, 0, 0, 58, 410, 1, 0, 0, 0, 60, 414, 1, 0, 0, 0, 62, 64,
		3, 2, 1, 0, 63, 62, 1, 0, 0, 0, 64, 67, 1, 0, 0, 0, 65, 63, 1, 0, 0, 0,
		65, 66, 1, 0, 0, 0, 66, 68, 1, 0, 0, 0, 67, 65, 1, 0, 0, 0, 68, 69, 5,
		0, 0, 1, 69, 1, 1, 0, 0, 0, 70, 74, 3, 4, 2, 0, 71, 74, 3, 8, 4, 0, 72,
		74, 3, 26, 13, 0, 73, 70, 1, 0, 0, 0, 73, 71, 1, 0, 0, 0, 73, 72, 1, 0,
		0, 0, 74, 3, 1, 0, 0, 0, 75, 76, 5, 2, 0, 0, 76, 77, 5, 52, 0, 0, 77, 79,
		5, 46, 0, 0, 78, 80, 3, 6, 3, 0, 79, 78, 1, 0, 0, 0, 80, 81, 1, 0, 0, 0,
		81, 79, 1, 0, 0, 0, 81, 82, 1, 0, 0, 0, 82, 83, 1, 0, 0, 0, 83, 84, 5,
		47, 0, 0, 84, 5, 1, 0, 0, 0, 85, 86, 3, 16, 8, 0, 86, 87, 5, 52, 0, 0,
		87, 7, 1, 0, 0, 0, 88, 90, 5, 3, 0, 0, 89, 91, 3, 10, 5, 0, 90, 89, 1,
		0, 0, 0, 90, 91, 1, 0, 0, 0, 91, 92, 1, 0, 0, 0, 92, 93, 5, 52, 0, 0, 93,
		95, 5, 42, 0, 0, 94, 96, 3, 12, 6, 0, 95, 94, 1, 0, 0, 0, 95, 96, 1, 0,
		0, 0, 96, 97, 1, 0, 0, 0, 97, 99, 5, 43, 0, 0, 98, 100, 3, 16, 8, 0, 99,
		98, 1, 0, 0, 0, 99, 100, 1, 0, 0, 0, 100, 101, 1, 0, 0, 0, 101, 102, 3,
		22, 11, 0, 102, 9, 1, 0, 0, 0, 103, 104, 5, 42, 0, 0, 104, 106, 5, 52,
		0, 0, 105, 107, 5, 37, 0, 0, 106, 105, 1, 0, 0, 0, 106, 107, 1, 0, 0, 0,
		107, 108, 1, 0, 0, 0, 108, 109, 5, 52, 0, 0, 109, 110, 5, 43, 0, 0, 110,
		11, 1, 0, 0, 0, 111, 116, 3, 14, 7, 0, 112, 113, 5, 48, 0, 0, 113, 115,
		3, 14, 7, 0, 114, 112, 1, 0, 0, 0, 115, 118, 1, 0, 0, 0, 116, 114, 1, 0,
		0, 0, 116, 117, 1, 0, 0, 0, 117, 13, 1, 0, 0, 0, 118, 116, 1, 0, 0, 0,
		119, 120, 5, 52, 0, 0, 120, 121, 3, 16, 8, 0, 121, 15, 1, 0, 0, 0, 122,
		126, 3, 20, 10, 0, 123, 126, 3, 18, 9, 0, 124, 126, 5, 52, 0, 0, 125, 122,
		1, 0, 0, 0, 125, 123, 1, 0, 0, 0, 125, 124, 1, 0, 0, 0, 126, 17, 1, 0,
		0, 0, 127, 128, 7, 0, 0, 0, 128, 19, 1, 0, 0, 0, 129, 130, 5, 44, 0, 0,
		130, 131, 5, 45, 0, 0, 131, 132, 3, 16, 8, 0, 132, 21, 1, 0, 0, 0, 133,
		137, 5, 46, 0, 0, 134, 136, 3, 24, 12, 0, 135, 134, 1, 0, 0, 0, 136, 139,
		1, 0, 0, 0, 137, 135, 1, 0, 0, 0, 137, 138, 1, 0, 0, 0, 138, 140, 1, 0,
		0, 0, 139, 137, 1, 0, 0, 0, 140, 141, 5, 47, 0, 0, 141, 23, 1, 0, 0, 0,
		142, 169, 3, 26, 13, 0, 143, 169, 3, 28, 14, 0, 144, 169, 3, 32, 16, 0,
		145, 147, 3, 48, 24, 0, 146, 148, 5, 51, 0, 0, 147, 146, 1, 0, 0, 0, 147,
		148, 1, 0, 0, 0, 148, 169, 1, 0, 0, 0, 149, 169, 3, 34, 17, 0, 150, 169,
		3, 36, 18, 0, 151, 169, 3, 42, 21, 0, 152, 154, 5, 11, 0, 0, 153, 155,
		5, 51, 0, 0, 154, 153, 1, 0, 0, 0, 154, 155, 1, 0, 0, 0, 155, 169, 1, 0,
		0, 0, 156, 158, 5, 12, 0, 0, 157, 159, 5, 51, 0, 0, 158, 157, 1, 0, 0,
		0, 158, 159, 1, 0, 0, 0, 159, 169, 1, 0, 0, 0, 160, 162, 5, 13, 0, 0, 161,
		163, 3, 48, 24, 0, 162, 161, 1, 0, 0, 0, 162, 163, 1, 0, 0, 0, 163, 165,
		1, 0, 0, 0, 164, 166, 5, 51, 0, 0, 165, 164, 1, 0, 0, 0, 165, 166, 1, 0,
		0, 0, 166, 169, 1, 0, 0, 0, 167, 169, 3, 22, 11, 0, 168, 142, 1, 0, 0,
		0, 168, 143, 1, 0, 0, 0, 168, 144, 1, 0, 0, 0, 168, 145, 1, 0, 0, 0, 168,
		149, 1, 0, 0, 0, 168, 150, 1, 0, 0, 0, 168, 151, 1, 0, 0, 0, 168, 152,
		1, 0, 0, 0, 168, 156, 1, 0, 0, 0, 168, 160, 1, 0, 0, 0, 168, 167, 1, 0,
		0, 0, 169, 25, 1, 0, 0, 0, 170, 172, 5, 1, 0, 0, 171, 170, 1, 0, 0, 0,
		171, 172, 1, 0, 0, 0, 172, 173, 1, 0, 0, 0, 173, 174, 5, 52, 0, 0, 174,
		177, 3, 16, 8, 0, 175, 176, 5, 41, 0, 0, 176, 178, 3, 48, 24, 0, 177, 175,
		1, 0, 0, 0, 177, 178, 1, 0, 0, 0, 178, 180, 1, 0, 0, 0, 179, 181, 5, 51,
		0, 0, 180, 179, 1, 0, 0, 0, 180, 181, 1, 0, 0, 0, 181, 192, 1, 0, 0, 0,
		182, 184, 5, 1, 0, 0, 183, 182, 1, 0, 0, 0, 183, 184, 1, 0, 0, 0, 184,
		185, 1, 0, 0, 0, 185, 186, 5, 52, 0, 0, 186, 187, 5, 26, 0, 0, 187, 189,
		3, 48, 24, 0, 188, 190, 5, 51, 0, 0, 189, 188, 1, 0, 0, 0, 189, 190, 1,
		0, 0, 0, 190, 192, 1, 0, 0, 0, 191, 171, 1, 0, 0, 0, 191, 183, 1, 0, 0,
		0, 192, 27, 1, 0, 0, 0, 193, 194, 3, 30, 15, 0, 194, 195, 7, 1, 0, 0, 195,
		197, 3, 48, 24, 0, 196, 198, 5, 51, 0, 0, 197, 196, 1, 0, 0, 0, 197, 198,
		1, 0, 0, 0, 198, 29, 1, 0, 0, 0, 199, 200, 6, 15, -1, 0, 200, 201, 5, 52,
		0, 0, 201, 212, 1, 0, 0, 0, 202, 203, 10, 2, 0, 0, 203, 204, 5, 44, 0,
		0, 204, 205, 3, 48, 24, 0, 205, 206, 5, 45, 0, 0, 206, 211, 1, 0, 0, 0,
		207, 208, 10, 1, 0, 0, 208, 209, 5, 49, 0, 0, 209, 211, 5, 52, 0, 0, 210,
		202, 1, 0, 0, 0, 210, 207, 1, 0, 0, 0, 211, 214, 1, 0, 0, 0, 212, 210,
		1, 0, 0, 0, 212, 213, 1, 0, 0, 0, 213, 31, 1, 0, 0, 0, 214, 212, 1, 0,
		0, 0, 215, 216, 3, 30, 15, 0, 216, 218, 7, 2, 0, 0, 217, 219, 5, 51, 0,
		0, 218, 217, 1, 0, 0, 0, 218, 219, 1, 0, 0, 0, 219, 33, 1, 0, 0, 0, 220,
		221, 5, 4, 0, 0, 221, 222, 3, 48, 24, 0, 222, 230, 3, 22, 11, 0, 223, 224,
		5, 5, 0, 0, 224, 225, 5, 4, 0, 0, 225, 226, 3, 48, 24, 0, 226, 227, 3,
		22, 11, 0, 227, 229, 1, 0, 0, 0, 228, 223, 1, 0, 0, 0, 229, 232, 1, 0,
		0, 0, 230, 228, 1, 0, 0, 0, 230, 231, 1, 0, 0, 0, 231, 235, 1, 0, 0, 0,
		232, 230, 1, 0, 0, 0, 233, 234, 5, 5, 0, 0, 234, 236, 3, 22, 11, 0, 235,
		233, 1, 0, 0, 0, 235, 236, 1, 0, 0, 0, 236, 35, 1, 0, 0, 0, 237, 238, 5,
		6, 0, 0, 238, 239, 3, 48, 24, 0, 239, 243, 5, 46, 0, 0, 240, 242, 3, 38,
		19, 0, 241, 240, 1, 0, 0, 0, 242, 245, 1, 0, 0, 0, 243, 241, 1, 0, 0, 0,
		243, 244, 1, 0, 0, 0, 244, 247, 1, 0, 0, 0, 245, 243, 1, 0, 0, 0, 246,
		248, 3, 40, 20, 0, 247, 246, 1, 0, 0, 0, 247, 248, 1, 0, 0, 0, 248, 249,
		1, 0, 0, 0, 249, 250, 5, 47, 0, 0, 250, 37, 1, 0, 0, 0, 251, 252, 5, 7,
		0, 0, 252, 253, 3, 48, 24, 0, 253, 257, 5, 50, 0, 0, 254, 256, 3, 24, 12,
		0, 255, 254, 1, 0, 0, 0, 256, 259, 1, 0, 0, 0, 257, 255, 1, 0, 0, 0, 257,
		258, 1, 0, 0, 0, 258, 39, 1, 0, 0, 0, 259, 257, 1, 0, 0, 0, 260, 261, 5,
		8, 0, 0, 261, 265, 5, 50, 0, 0, 262, 264, 3, 24, 12, 0, 263, 262, 1, 0,
		0, 0, 264, 267, 1, 0, 0, 0, 265, 263, 1, 0, 0, 0, 265, 266, 1, 0, 0, 0,
		266, 41, 1, 0, 0, 0, 267, 265, 1, 0, 0, 0, 268, 269, 5, 9, 0, 0, 269, 270,
		3, 48, 24, 0, 270, 271, 3, 22, 11, 0, 271, 294, 1, 0, 0, 0, 272, 274, 5,
		9, 0, 0, 273, 275, 3, 44, 22, 0, 274, 273, 1, 0, 0, 0, 274, 275, 1, 0,
		0, 0, 275, 276, 1, 0, 0, 0, 276, 278, 5, 51, 0, 0, 277, 279, 3, 48, 24,
		0, 278, 277, 1, 0, 0, 0, 278, 279, 1, 0, 0, 0, 279, 280, 1, 0, 0, 0, 280,
		282, 5, 51, 0, 0, 281, 283, 3, 46, 23, 0, 282, 281, 1, 0, 0, 0, 282, 283,
		1, 0, 0, 0, 283, 284, 1, 0, 0, 0, 284, 294, 3, 22, 11, 0, 285, 286, 5,
		9, 0, 0, 286, 287, 5, 52, 0, 0, 287, 288, 5, 48, 0, 0, 288, 289, 5, 52,
		0, 0, 289, 290, 5, 10, 0, 0, 290, 291, 3, 48, 24, 0, 291, 292, 3, 22, 11,
		0, 292, 294, 1, 0, 0, 0, 293, 268, 1, 0, 0, 0, 293, 272, 1, 0, 0, 0, 293,
		285, 1, 0, 0, 0, 294, 43, 1, 0, 0, 0, 295, 298, 3, 26, 13, 0, 296, 298,
		3, 28, 14, 0, 297, 295, 1, 0, 0, 0, 297, 296, 1, 0, 0, 0, 298, 45, 1, 0,
		0, 0, 299, 302, 3, 28, 14, 0, 300, 302, 3, 32, 16, 0, 301, 299, 1, 0, 0,
		0, 301, 300, 1, 0, 0, 0, 302, 47, 1, 0, 0, 0, 303, 304, 6, 24, -1, 0, 304,
		305, 5, 42, 0, 0, 305, 306, 3, 48, 24, 0, 306, 307, 5, 43, 0, 0, 307, 315,
		1, 0, 0, 0, 308, 309, 7, 3, 0, 0, 309, 315, 3, 48, 24, 11, 310, 315, 3,
		52, 26, 0, 311, 315, 3, 56, 28, 0, 312, 315, 3, 60, 30, 0, 313, 315, 5,
		52, 0, 0, 314, 303, 1, 0, 0, 0, 314, 308, 1, 0, 0, 0, 314, 310, 1, 0, 0,
		0, 314, 311, 1, 0, 0, 0, 314, 312, 1, 0, 0, 0, 314, 313, 1, 0, 0, 0, 315,
		350, 1, 0, 0, 0, 316, 317, 10, 10, 0, 0, 317, 318, 7, 4, 0, 0, 318, 349,
		3, 48, 24, 11, 319, 320, 10, 9, 0, 0, 320, 321, 7, 5, 0, 0, 321, 349, 3,
		48, 24, 10, 322, 323, 10, 8, 0, 0, 323, 324, 7, 6, 0, 0, 324, 349, 3, 48,
		24, 9, 325, 326, 10, 7, 0, 0, 326, 327, 7, 7, 0, 0, 327, 349, 3, 48, 24,
		8, 328, 329, 10, 6, 0, 0, 329, 330, 5, 31, 0, 0, 330, 349, 3, 48, 24, 7,
		331, 332, 10, 5, 0, 0, 332, 333, 5, 32, 0, 0, 333, 349, 3, 48, 24, 6, 334,
		335, 10, 15, 0, 0, 335, 337, 5, 42, 0, 0, 336, 338, 3, 50, 25, 0, 337,
		336, 1, 0, 0, 0, 337, 338, 1, 0, 0, 0, 338, 339, 1, 0, 0, 0, 339, 349,
		5, 43, 0, 0, 340, 341, 10, 14, 0, 0, 341, 342, 5, 44, 0, 0, 342, 343, 3,
		48, 24, 0, 343, 344, 5, 45, 0, 0, 344, 349, 1, 0, 0, 0, 345, 346, 10, 13,
		0, 0, 346, 347, 5, 49, 0, 0, 347, 349, 5, 52, 0, 0, 348, 316, 1, 0, 0,
		0, 348, 319, 1, 0, 0, 0, 348, 322, 1, 0, 0, 0, 348, 325, 1, 0, 0, 0, 348,
		328, 1, 0, 0, 0, 348, 331, 1, 0, 0, 0, 348, 334, 1, 0, 0, 0, 348, 340,
		1, 0, 0, 0, 348, 345, 1, 0, 0, 0, 349, 352, 1, 0, 0, 0, 350, 348, 1, 0,
		0, 0, 350, 351, 1, 0, 0, 0, 351, 49, 1, 0, 0, 0, 352, 350, 1, 0, 0, 0,
		353, 358, 3, 48, 24, 0, 354, 355, 5, 48, 0, 0, 355, 357, 3, 48, 24, 0,
		356, 354, 1, 0, 0, 0, 357, 360, 1, 0, 0, 0, 358, 356, 1, 0, 0, 0, 358,
		359, 1, 0, 0, 0, 359, 51, 1, 0, 0, 0, 360, 358, 1, 0, 0, 0, 361, 362, 5,
		44, 0, 0, 362, 363, 5, 45, 0, 0, 363, 364, 3, 16, 8, 0, 364, 366, 5, 46,
		0, 0, 365, 367, 3, 50, 25, 0, 366, 365, 1, 0, 0, 0, 366, 367, 1, 0, 0,
		0, 367, 368, 1, 0, 0, 0, 368, 369, 5, 47, 0, 0, 369, 388, 1, 0, 0, 0, 370,
		371, 5, 44, 0, 0, 371, 372, 5, 45, 0, 0, 372, 373, 3, 16, 8, 0, 373, 374,
		5, 46, 0, 0, 374, 379, 3, 54, 27, 0, 375, 376, 5, 48, 0, 0, 376, 378, 3,
		54, 27, 0, 377, 375, 1, 0, 0, 0, 378, 381, 1, 0, 0, 0, 379, 377, 1, 0,
		0, 0, 379, 380, 1, 0, 0, 0, 380, 383, 1, 0, 0, 0, 381, 379, 1, 0, 0, 0,
		382, 384, 5, 48, 0, 0, 383, 382, 1, 0, 0, 0, 383, 384, 1, 0, 0, 0, 384,
		385, 1, 0, 0, 0, 385, 386, 5, 47, 0, 0, 386, 388, 1, 0, 0, 0, 387, 361,
		1, 0, 0, 0, 387, 370, 1, 0, 0, 0, 388, 53, 1, 0, 0, 0, 389, 391, 5, 46,
		0, 0, 390, 392, 3, 50, 25, 0, 391, 390, 1, 0, 0, 0, 391, 392, 1, 0, 0,
		0, 392, 393, 1, 0, 0, 0, 393, 394, 5, 47, 0, 0, 394, 55, 1, 0, 0, 0, 395,
		396, 5, 52, 0, 0, 396, 397, 5, 46, 0, 0, 397, 402, 3, 58, 29, 0, 398, 399,
		5, 48, 0, 0, 399, 401, 3, 58, 29, 0, 400, 398, 1, 0, 0, 0, 401, 404, 1,
		0, 0, 0, 402, 400, 1, 0, 0, 0, 402, 403, 1, 0, 0, 0, 403, 406, 1, 0, 0,
		0, 404, 402, 1, 0, 0, 0, 405, 407, 5, 48, 0, 0, 406, 405, 1, 0, 0, 0, 406,
		407, 1, 0, 0, 0, 407, 408, 1, 0, 0, 0, 408, 409, 5, 47, 0, 0, 409, 57,
		1, 0, 0, 0, 410, 411, 5, 52, 0, 0, 411, 412, 5, 50, 0, 0, 412, 413, 3,
		48, 24, 0, 413, 59, 1, 0, 0, 0, 414, 415, 7, 8, 0, 0, 415, 61, 1, 0, 0,
		0, 50, 65, 73, 81, 90, 95, 99, 106, 116, 125, 137, 147, 154, 158, 162,
		165, 168, 171, 177, 180, 183, 189, 191, 197, 210, 212, 218, 230, 235, 243,
		247, 257, 265, 274, 278, 282, 293, 297, 301, 314, 337, 348, 350, 358, 366,
		379, 383, 387, 391, 402, 406,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// VLangCherryParserInit initializes any static state used to implement VLangCherryParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewVLangCherryParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func VLangCherryParserInit() {
	staticData := &VLangCherryParserStaticData
	staticData.once.Do(vlangcherryParserInit)
}

// NewVLangCherryParser produces a new parser instance for the optional input antlr.TokenStream.
func NewVLangCherryParser(input antlr.TokenStream) *VLangCherryParser {
	VLangCherryParserInit()
	this := new(VLangCherryParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &VLangCherryParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "VLangCherry.g4"

	return this
}

// VLangCherryParser tokens.
const (
	VLangCherryParserEOF               = antlr.TokenEOF
	VLangCherryParserMUT               = 1
	VLangCherryParserSTRUCT            = 2
	VLangCherryParserFUNC              = 3
	VLangCherryParserIF                = 4
	VLangCherryParserELSE              = 5
	VLangCherryParserSWITCH            = 6
	VLangCherryParserCASE              = 7
	VLangCherryParserDEFAULT           = 8
	VLangCherryParserFOR               = 9
	VLangCherryParserIN                = 10
	VLangCherryParserBREAK             = 11
	VLangCherryParserCONTINUE          = 12
	VLangCherryParserRETURN            = 13
	VLangCherryParserNIL               = 14
	VLangCherryParserTRUE              = 15
	VLangCherryParserFALSE             = 16
	VLangCherryParserTIPO_INT          = 17
	VLangCherryParserTIPO_FLOAT        = 18
	VLangCherryParserTIPO_STRING       = 19
	VLangCherryParserTIPO_BOOL         = 20
	VLangCherryParserTIPO_RUNE         = 21
	VLangCherryParserMASIGUAL          = 22
	VLangCherryParserMENOSIGUAL        = 23
	VLangCherryParserINCREMENTO        = 24
	VLangCherryParserDECREMENTO        = 25
	VLangCherryParserASIGNAINFERIDA    = 26
	VLangCherryParserIGUAL             = 27
	VLangCherryParserDIFERENTE         = 28
	VLangCherryParserMENORIGUAL        = 29
	VLangCherryParserMAYORIGUAL        = 30
	VLangCherryParserY                 = 31
	VLangCherryParserO                 = 32
	VLangCherryParserMENOR             = 33
	VLangCherryParserMAYOR             = 34
	VLangCherryParserMAS               = 35
	VLangCherryParserMENOS             = 36
	VLangCherryParserPOR               = 37
	VLangCherryParserDIV               = 38
	VLangCherryParserMODULO            = 39
	VLangCherryParserNOT               = 40
	VLangCherryParserASIGNA            = 41
	VLangCherryParserPARIZQ            = 42
	VLangCherryParserPARDER            = 43
	VLangCherryParserCORIZQ            = 44
	VLangCherryParserCORDER            = 45
	VLangCherryParserLLAIZQ            = 46
	VLangCherryParserLLADER            = 47
	VLangCherryParserCOMA              = 48
	VLangCherryParserPUNTO             = 49
	VLangCherryParserDOSPUNTOS         = 50
	VLangCherryParserPUNTOCOMA         = 51
	VLangCherryParserID                = 52
	VLangCherryParserENTERO            = 53
	VLangCherryParserDECIMAL           = 54
	VLangCherryParserCADENA            = 55
	VLangCherryParserRUNE              = 56
	VLangCherryParserCOMENTARIO_LINEA  = 57
	VLangCherryParserCOMENTARIO_BLOQUE = 58
	VLangCherryParserWS                = 59
)

// VLangCherryParser rules.
const (
	VLangCherryParserRULE_programa             = 0
	VLangCherryParserRULE_declaracionGlobal    = 1
	VLangCherryParserRULE_declaracionStruct    = 2
	VLangCherryParserRULE_campoStruct          = 3
	VLangCherryParserRULE_declaracionFuncion   = 4
	VLangCherryParserRULE_receptor             = 5
	VLangCherryParserRULE_listaParametros      = 6
	VLangCherryParserRULE_parametro            = 7
	VLangCherryParserRULE_tipo                 = 8
	VLangCherryParserRULE_tipoPrimitivo        = 9
	VLangCherryParserRULE_tipoSlice            = 10
	VLangCherryParserRULE_bloque               = 11
	VLangCherryParserRULE_sentencia            = 12
	VLangCherryParserRULE_declaracionVariable  = 13
	VLangCherryParserRULE_asignacion           = 14
	VLangCherryParserRULE_lugar                = 15
	VLangCherryParserRULE_incrementoDecremento = 16
	VLangCherryParserRULE_sentenciaIf          = 17
	VLangCherryParserRULE_sentenciaSwitch      = 18
	VLangCherryParserRULE_casoSwitch           = 19
	VLangCherryParserRULE_defaultSwitch        = 20
	VLangCherryParserRULE_sentenciaFor         = 21
	VLangCherryParserRULE_forInit              = 22
	VLangCherryParserRULE_forActualizacion     = 23
	VLangCherryParserRULE_expr                 = 24
	VLangCherryParserRULE_listaArgumentos      = 25
	VLangCherryParserRULE_literalSlice         = 26
	VLangCherryParserRULE_filaSlice            = 27
	VLangCherryParserRULE_literalStruct        = 28
	VLangCherryParserRULE_campoValor           = 29
	VLangCherryParserRULE_literal              = 30
)

// IProgramaContext is an interface to support dynamic dispatch.
type IProgramaContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EOF() antlr.TerminalNode
	AllDeclaracionGlobal() []IDeclaracionGlobalContext
	DeclaracionGlobal(i int) IDeclaracionGlobalContext

	// IsProgramaContext differentiates from other interfaces.
	IsProgramaContext()
}

type ProgramaContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgramaContext() *ProgramaContext {
	var p = new(ProgramaContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_programa
	return p
}

func InitEmptyProgramaContext(p *ProgramaContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_programa
}

func (*ProgramaContext) IsProgramaContext() {}

func NewProgramaContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgramaContext {
	var p = new(ProgramaContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_programa

	return p
}

func (s *ProgramaContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgramaContext) EOF() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserEOF, 0)
}

func (s *ProgramaContext) AllDeclaracionGlobal() []IDeclaracionGlobalContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDeclaracionGlobalContext); ok {
			len++
		}
	}

	tst := make([]IDeclaracionGlobalContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDeclaracionGlobalContext); ok {
			tst[i] = t.(IDeclaracionGlobalContext)
			i++
		}
	}

	return tst
}

func (s *ProgramaContext) DeclaracionGlobal(i int) IDeclaracionGlobalContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDeclaracionGlobalContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDeclaracionGlobalContext)
}

func (s *ProgramaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgramaContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProgramaContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitPrograma(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) Programa() (localctx IProgramaContext) {
	localctx = NewProgramaContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, VLangCherryParserRULE_programa)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(65)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&4503599627370510) != 0 {
		{
			p.SetState(62)
			p.DeclaracionGlobal()
		}

		p.SetState(67)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(68)
		p.Match(VLangCherryParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDeclaracionGlobalContext is an interface to support dynamic dispatch.
type IDeclaracionGlobalContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DeclaracionStruct() IDeclaracionStructContext
	DeclaracionFuncion() IDeclaracionFuncionContext
	DeclaracionVariable() IDeclaracionVariableContext

	// IsDeclaracionGlobalContext differentiates from other interfaces.
	IsDeclaracionGlobalContext()
}

type DeclaracionGlobalContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDeclaracionGlobalContext() *DeclaracionGlobalContext {
	var p = new(DeclaracionGlobalContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_declaracionGlobal
	return p
}

func InitEmptyDeclaracionGlobalContext(p *DeclaracionGlobalContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_declaracionGlobal
}

func (*DeclaracionGlobalContext) IsDeclaracionGlobalContext() {}

func NewDeclaracionGlobalContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclaracionGlobalContext {
	var p = new(DeclaracionGlobalContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_declaracionGlobal

	return p
}

func (s *DeclaracionGlobalContext) GetParser() antlr.Parser { return s.parser }

func (s *DeclaracionGlobalContext) DeclaracionStruct() IDeclaracionStructContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDeclaracionStructContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDeclaracionStructContext)
}

func (s *DeclaracionGlobalContext) DeclaracionFuncion() IDeclaracionFuncionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDeclaracionFuncionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDeclaracionFuncionContext)
}

func (s *DeclaracionGlobalContext) DeclaracionVariable() IDeclaracionVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDeclaracionVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDeclaracionVariableContext)
}

func (s *DeclaracionGlobalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeclaracionGlobalContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DeclaracionGlobalContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitDeclaracionGlobal(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) DeclaracionGlobal() (localctx IDeclaracionGlobalContext) {
	localctx = NewDeclaracionGlobalContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, VLangCherryParserRULE_declaracionGlobal)
	p.SetState(73)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case VLangCherryParserSTRUCT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(70)
			p.DeclaracionStruct()
		}

	case VLangCherryParserFUNC:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(71)
			p.DeclaracionFuncion()
		}

	case VLangCherryParserMUT, VLangCherryParserID:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(72)
			p.DeclaracionVariable()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDeclaracionStructContext is an interface to support dynamic dispatch.
type IDeclaracionStructContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRUCT() antlr.TerminalNode
	ID() antlr.TerminalNode
	LLAIZQ() antlr.TerminalNode
	LLADER() antlr.TerminalNode
	AllCampoStruct() []ICampoStructContext
	CampoStruct(i int) ICampoStructContext

	// IsDeclaracionStructContext differentiates from other interfaces.
	IsDeclaracionStructContext()
}

type DeclaracionStructContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDeclaracionStructContext() *DeclaracionStructContext {
	var p = new(DeclaracionStructContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_declaracionStruct
	return p
}

func InitEmptyDeclaracionStructContext(p *DeclaracionStructContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_declaracionStruct
}

func (*DeclaracionStructContext) IsDeclaracionStructContext() {}

func NewDeclaracionStructContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclaracionStructContext {
	var p = new(DeclaracionStructContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_declaracionStruct

	return p
}

func (s *DeclaracionStructContext) GetParser() antlr.Parser { return s.parser }

func (s *DeclaracionStructContext) STRUCT() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserSTRUCT, 0)
}

func (s *DeclaracionStructContext) ID() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, 0)
}

func (s *DeclaracionStructContext) LLAIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserLLAIZQ, 0)
}

func (s *DeclaracionStructContext) LLADER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserLLADER, 0)
}

func (s *DeclaracionStructContext) AllCampoStruct() []ICampoStructContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICampoStructContext); ok {
			len++
		}
	}

	tst := make([]ICampoStructContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICampoStructContext); ok {
			tst[i] = t.(ICampoStructContext)
			i++
		}
	}

	return tst
}

func (s *DeclaracionStructContext) CampoStruct(i int) ICampoStructContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICampoStructContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICampoStructContext)
}

func (s *DeclaracionStructContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeclaracionStructContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DeclaracionStructContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitDeclaracionStruct(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) DeclaracionStruct() (localctx IDeclaracionStructContext) {
	localctx = NewDeclaracionStructContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, VLangCherryParserRULE_declaracionStruct)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(75)
		p.Match(VLangCherryParserSTRUCT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(76)
		p.Match(VLangCherryParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(77)
		p.Match(VLangCherryParserLLAIZQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(79)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&4521191817478144) != 0) {
		{
			p.SetState(78)
			p.CampoStruct()
		}

		p.SetState(81)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(83)
		p.Match(VLangCherryParserLLADER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICampoStructContext is an interface to support dynamic dispatch.
type ICampoStructContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Tipo() ITipoContext
	ID() antlr.TerminalNode

	// IsCampoStructContext differentiates from other interfaces.
	IsCampoStructContext()
}

type CampoStructContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCampoStructContext() *CampoStructContext {
	var p = new(CampoStructContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_campoStruct
	return p
}

func InitEmptyCampoStructContext(p *CampoStructContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_campoStruct
}

func (*CampoStructContext) IsCampoStructContext() {}

func NewCampoStructContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CampoStructContext {
	var p = new(CampoStructContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_campoStruct

	return p
}

func (s *CampoStructContext) GetParser() antlr.Parser { return s.parser }

func (s *CampoStructContext) Tipo() ITipoContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITipoContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITipoContext)
}

func (s *CampoStructContext) ID() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, 0)
}

func (s *CampoStructContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CampoStructContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CampoStructContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitCampoStruct(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) CampoStruct() (localctx ICampoStructContext) {
	localctx = NewCampoStructContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, VLangCherryParserRULE_campoStruct)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(85)
		p.Tipo()
	}
	{
		p.SetState(86)
		p.Match(VLangCherryParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDeclaracionFuncionContext is an interface to support dynamic dispatch.
type IDeclaracionFuncionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FUNC() antlr.TerminalNode
	ID() antlr.TerminalNode
	PARIZQ() antlr.TerminalNode
	PARDER() antlr.TerminalNode
	Bloque() IBloqueContext
	Receptor() IReceptorContext
	ListaParametros() IListaParametrosContext
	Tipo() ITipoContext

	// IsDeclaracionFuncionContext differentiates from other interfaces.
	IsDeclaracionFuncionContext()
}

type DeclaracionFuncionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDeclaracionFuncionContext() *DeclaracionFuncionContext {
	var p = new(DeclaracionFuncionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_declaracionFuncion
	return p
}

func InitEmptyDeclaracionFuncionContext(p *DeclaracionFuncionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_declaracionFuncion
}

func (*DeclaracionFuncionContext) IsDeclaracionFuncionContext() {}

func NewDeclaracionFuncionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclaracionFuncionContext {
	var p = new(DeclaracionFuncionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_declaracionFuncion

	return p
}

func (s *DeclaracionFuncionContext) GetParser() antlr.Parser { return s.parser }

func (s *DeclaracionFuncionContext) FUNC() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserFUNC, 0)
}

func (s *DeclaracionFuncionContext) ID() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, 0)
}

func (s *DeclaracionFuncionContext) PARIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPARIZQ, 0)
}

func (s *DeclaracionFuncionContext) PARDER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPARDER, 0)
}

func (s *DeclaracionFuncionContext) Bloque() IBloqueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBloqueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBloqueContext)
}

func (s *DeclaracionFuncionContext) Receptor() IReceptorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IReceptorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IReceptorContext)
}

func (s *DeclaracionFuncionContext) ListaParametros() IListaParametrosContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IListaParametrosContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IListaParametrosContext)
}

func (s *DeclaracionFuncionContext) Tipo() ITipoContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITipoContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITipoContext)
}

func (s *DeclaracionFuncionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeclaracionFuncionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DeclaracionFuncionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitDeclaracionFuncion(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) DeclaracionFuncion() (localctx IDeclaracionFuncionContext) {
	localctx = NewDeclaracionFuncionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, VLangCherryParserRULE_declaracionFuncion)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(88)
		p.Match(VLangCherryParserFUNC)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(90)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == VLangCherryParserPARIZQ {
		{
			p.SetState(89)
			p.Receptor()
		}

	}
	{
		p.SetState(92)
		p.Match(VLangCherryParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(93)
		p.Match(VLangCherryParserPARIZQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(95)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == VLangCherryParserID {
		{
			p.SetState(94)
			p.ListaParametros()
		}

	}
	{
		p.SetState(97)
		p.Match(VLangCherryParserPARDER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(99)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&4521191817478144) != 0 {
		{
			p.SetState(98)
			p.Tipo()
		}

	}
	{
		p.SetState(101)
		p.Bloque()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IReceptorContext is an interface to support dynamic dispatch.
type IReceptorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PARIZQ() antlr.TerminalNode
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	PARDER() antlr.TerminalNode
	POR() antlr.TerminalNode

	// IsReceptorContext differentiates from other interfaces.
	IsReceptorContext()
}

type ReceptorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyReceptorContext() *ReceptorContext {
	var p = new(ReceptorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_receptor
	return p
}

func InitEmptyReceptorContext(p *ReceptorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_receptor
}

func (*ReceptorContext) IsReceptorContext() {}

func NewReceptorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReceptorContext {
	var p = new(ReceptorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_receptor

	return p
}

func (s *ReceptorContext) GetParser() antlr.Parser { return s.parser }

func (s *ReceptorContext) PARIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPARIZQ, 0)
}

func (s *ReceptorContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(VLangCherryParserID)
}

func (s *ReceptorContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, i)
}

func (s *ReceptorContext) PARDER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPARDER, 0)
}

func (s *ReceptorContext) POR() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPOR, 0)
}

func (s *ReceptorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReceptorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ReceptorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitReceptor(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) Receptor() (localctx IReceptorContext) {
	localctx = NewReceptorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, VLangCherryParserRULE_receptor)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(103)
		p.Match(VLangCherryParserPARIZQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(104)
		p.Match(VLangCherryParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(106)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == VLangCherryParserPOR {
		{
			p.SetState(105)
			p.Match(VLangCherryParserPOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(108)
		p.Match(VLangCherryParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(109)
		p.Match(VLangCherryParserPARDER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IListaParametrosContext is an interface to support dynamic dispatch.
type IListaParametrosContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllParametro() []IParametroContext
	Parametro(i int) IParametroContext
	AllCOMA() []antlr.TerminalNode
	COMA(i int) antlr.TerminalNode

	// IsListaParametrosContext differentiates from other interfaces.
	IsListaParametrosContext()
}

type ListaParametrosContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyListaParametrosContext() *ListaParametrosContext {
	var p = new(ListaParametrosContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_listaParametros
	return p
}

func InitEmptyListaParametrosContext(p *ListaParametrosContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_listaParametros
}

func (*ListaParametrosContext) IsListaParametrosContext() {}

func NewListaParametrosContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ListaParametrosContext {
	var p = new(ListaParametrosContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_listaParametros

	return p
}

func (s *ListaParametrosContext) GetParser() antlr.Parser { return s.parser }

func (s *ListaParametrosContext) AllParametro() []IParametroContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IParametroContext); ok {
			len++
		}
	}

	tst := make([]IParametroContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IParametroContext); ok {
			tst[i] = t.(IParametroContext)
			i++
		}
	}

	return tst
}

func (s *ListaParametrosContext) Parametro(i int) IParametroContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParametroContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParametroContext)
}

func (s *ListaParametrosContext) AllCOMA() []antlr.TerminalNode {
	return s.GetTokens(VLangCherryParserCOMA)
}

func (s *ListaParametrosContext) COMA(i int) antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCOMA, i)
}

func (s *ListaParametrosContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ListaParametrosContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ListaParametrosContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitListaParametros(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) ListaParametros() (localctx IListaParametrosContext) {
	localctx = NewListaParametrosContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, VLangCherryParserRULE_listaParametros)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(111)
		p.Parametro()
	}
	p.SetState(116)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == VLangCherryParserCOMA {
		{
			p.SetState(112)
			p.Match(VLangCherryParserCOMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(113)
			p.Parametro()
		}

		p.SetState(118)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IParametroContext is an interface to support dynamic dispatch.
type IParametroContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	Tipo() ITipoContext

	// IsParametroContext differentiates from other interfaces.
	IsParametroContext()
}

type ParametroContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParametroContext() *ParametroContext {
	var p = new(ParametroContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_parametro
	return p
}

func InitEmptyParametroContext(p *ParametroContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_parametro
}

func (*ParametroContext) IsParametroContext() {}

func NewParametroContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParametroContext {
	var p = new(ParametroContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_parametro

	return p
}

func (s *ParametroContext) GetParser() antlr.Parser { return s.parser }

func (s *ParametroContext) ID() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, 0)
}

func (s *ParametroContext) Tipo() ITipoContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITipoContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITipoContext)
}

func (s *ParametroContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParametroContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParametroContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitParametro(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) Parametro() (localctx IParametroContext) {
	localctx = NewParametroContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, VLangCherryParserRULE_parametro)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(119)
		p.Match(VLangCherryParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(120)
		p.Tipo()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITipoContext is an interface to support dynamic dispatch.
type ITipoContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TipoSlice() ITipoSliceContext
	TipoPrimitivo() ITipoPrimitivoContext
	ID() antlr.TerminalNode

	// IsTipoContext differentiates from other interfaces.
	IsTipoContext()
}

type TipoContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTipoContext() *TipoContext {
	var p = new(TipoContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_tipo
	return p
}

func InitEmptyTipoContext(p *TipoContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_tipo
}

func (*TipoContext) IsTipoContext() {}

func NewTipoContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TipoContext {
	var p = new(TipoContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_tipo

	return p
}

func (s *TipoContext) GetParser() antlr.Parser { return s.parser }

func (s *TipoContext) TipoSlice() ITipoSliceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITipoSliceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITipoSliceContext)
}

func (s *TipoContext) TipoPrimitivo() ITipoPrimitivoContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITipoPrimitivoContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITipoPrimitivoContext)
}

func (s *TipoContext) ID() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, 0)
}

func (s *TipoContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TipoContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TipoContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitTipo(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) Tipo() (localctx ITipoContext) {
	localctx = NewTipoContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, VLangCherryParserRULE_tipo)
	p.SetState(125)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case VLangCherryParserCORIZQ:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(122)
			p.TipoSlice()
		}

	case VLangCherryParserTIPO_INT, VLangCherryParserTIPO_FLOAT, VLangCherryParserTIPO_STRING, VLangCherryParserTIPO_BOOL, VLangCherryParserTIPO_RUNE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(123)
			p.TipoPrimitivo()
		}

	case VLangCherryParserID:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(124)
			p.Match(VLangCherryParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITipoPrimitivoContext is an interface to support dynamic dispatch.
type ITipoPrimitivoContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TIPO_INT() antlr.TerminalNode
	TIPO_FLOAT() antlr.TerminalNode
	TIPO_STRING() antlr.TerminalNode
	TIPO_BOOL() antlr.TerminalNode
	TIPO_RUNE() antlr.TerminalNode

	// IsTipoPrimitivoContext differentiates from other interfaces.
	IsTipoPrimitivoContext()
}

type TipoPrimitivoContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTipoPrimitivoContext() *TipoPrimitivoContext {
	var p = new(TipoPrimitivoContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_tipoPrimitivo
	return p
}

func InitEmptyTipoPrimitivoContext(p *TipoPrimitivoContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_tipoPrimitivo
}

func (*TipoPrimitivoContext) IsTipoPrimitivoContext() {}

func NewTipoPrimitivoContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TipoPrimitivoContext {
	var p = new(TipoPrimitivoContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_tipoPrimitivo

	return p
}

func (s *TipoPrimitivoContext) GetParser() antlr.Parser { return s.parser }

func (s *TipoPrimitivoContext) TIPO_INT() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserTIPO_INT, 0)
}

func (s *TipoPrimitivoContext) TIPO_FLOAT() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserTIPO_FLOAT, 0)
}

func (s *TipoPrimitivoContext) TIPO_STRING() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserTIPO_STRING, 0)
}

func (s *TipoPrimitivoContext) TIPO_BOOL() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserTIPO_BOOL, 0)
}

func (s *TipoPrimitivoContext) TIPO_RUNE() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserTIPO_RUNE, 0)
}

func (s *TipoPrimitivoContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TipoPrimitivoContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TipoPrimitivoContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitTipoPrimitivo(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) TipoPrimitivo() (localctx ITipoPrimitivoContext) {
	localctx = NewTipoPrimitivoContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, VLangCherryParserRULE_tipoPrimitivo)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(127)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&4063232) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITipoSliceContext is an interface to support dynamic dispatch.
type ITipoSliceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CORIZQ() antlr.TerminalNode
	CORDER() antlr.TerminalNode
	Tipo() ITipoContext

	// IsTipoSliceContext differentiates from other interfaces.
	IsTipoSliceContext()
}

type TipoSliceContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTipoSliceContext() *TipoSliceContext {
	var p = new(TipoSliceContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_tipoSlice
	return p
}

func InitEmptyTipoSliceContext(p *TipoSliceContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_tipoSlice
}

func (*TipoSliceContext) IsTipoSliceContext() {}

func NewTipoSliceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TipoSliceContext {
	var p = new(TipoSliceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_tipoSlice

	return p
}

func (s *TipoSliceContext) GetParser() antlr.Parser { return s.parser }

func (s *TipoSliceContext) CORIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCORIZQ, 0)
}

func (s *TipoSliceContext) CORDER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCORDER, 0)
}

func (s *TipoSliceContext) Tipo() ITipoContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITipoContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITipoContext)
}

func (s *TipoSliceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TipoSliceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TipoSliceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitTipoSlice(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) TipoSlice() (localctx ITipoSliceContext) {
	localctx = NewTipoSliceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, VLangCherryParserRULE_tipoSlice)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(129)
		p.Match(VLangCherryParserCORIZQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(130)
		p.Match(VLangCherryParserCORDER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(131)
		p.Tipo()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBloqueContext is an interface to support dynamic dispatch.
type IBloqueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LLAIZQ() antlr.TerminalNode
	LLADER() antlr.TerminalNode
	AllSentencia() []ISentenciaContext
	Sentencia(i int) ISentenciaContext

	// IsBloqueContext differentiates from other interfaces.
	IsBloqueContext()
}

type BloqueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBloqueContext() *BloqueContext {
	var p = new(BloqueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_bloque
	return p
}

func InitEmptyBloqueContext(p *BloqueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_bloque
}

func (*BloqueContext) IsBloqueContext() {}

func NewBloqueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BloqueContext {
	var p = new(BloqueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_bloque

	return p
}

func (s *BloqueContext) GetParser() antlr.Parser { return s.parser }

func (s *BloqueContext) LLAIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserLLAIZQ, 0)
}

func (s *BloqueContext) LLADER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserLLADER, 0)
}

func (s *BloqueContext) AllSentencia() []ISentenciaContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISentenciaContext); ok {
			len++
		}
	}

	tst := make([]ISentenciaContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISentenciaContext); ok {
			tst[i] = t.(ISentenciaContext)
			i++
		}
	}

	return tst
}

func (s *BloqueContext) Sentencia(i int) ISentenciaContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISentenciaContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISentenciaContext)
}

func (s *BloqueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BloqueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BloqueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitBloque(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) Bloque() (localctx IBloqueContext) {
	localctx = NewBloqueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, VLangCherryParserRULE_bloque)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(133)
		p.Match(VLangCherryParserLLAIZQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(137)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&139705115656452690) != 0 {
		{
			p.SetState(134)
			p.Sentencia()
		}

		p.SetState(139)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(140)
		p.Match(VLangCherryParserLLADER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISentenciaContext is an interface to support dynamic dispatch.
type ISentenciaContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsSentenciaContext differentiates from other interfaces.
	IsSentenciaContext()
}

type SentenciaContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySentenciaContext() *SentenciaContext {
	var p = new(SentenciaContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_sentencia
	return p
}

func InitEmptySentenciaContext(p *SentenciaContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_sentencia
}

func (*SentenciaContext) IsSentenciaContext() {}

func NewSentenciaContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SentenciaContext {
	var p = new(SentenciaContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_sentencia

	return p
}

func (s *SentenciaContext) GetParser() antlr.Parser { return s.parser }

func (s *SentenciaContext) CopyAll(ctx *SentenciaContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *SentenciaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentenciaContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SentIncDecContext struct {
	SentenciaContext
}

func NewSentIncDecContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SentIncDecContext {
	var p = new(SentIncDecContext)

	InitEmptySentenciaContext(&p.SentenciaContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaContext))

	return p
}

func (s *SentIncDecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentIncDecContext) IncrementoDecremento() IIncrementoDecrementoContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIncrementoDecrementoContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIncrementoDecrementoContext)
}

func (s *SentIncDecContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitSentIncDec(s)

	default:
		return t.VisitChildren(s)
	}
}

type SentBloqueContext struct {
	SentenciaContext
}

func NewSentBloqueContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SentBloqueContext {
	var p = new(SentBloqueContext)

	InitEmptySentenciaContext(&p.SentenciaContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaContext))

	return p
}

func (s *SentBloqueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentBloqueContext) Bloque() IBloqueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBloqueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBloqueContext)
}

func (s *SentBloqueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitSentBloque(s)

	default:
		return t.VisitChildren(s)
	}
}

type SentAsignacionContext struct {
	SentenciaContext
}

func NewSentAsignacionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SentAsignacionContext {
	var p = new(SentAsignacionContext)

	InitEmptySentenciaContext(&p.SentenciaContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaContext))

	return p
}

func (s *SentAsignacionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentAsignacionContext) Asignacion() IAsignacionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAsignacionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAsignacionContext)
}

func (s *SentAsignacionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitSentAsignacion(s)

	default:
		return t.VisitChildren(s)
	}
}

type SentDeclaracionContext struct {
	SentenciaContext
}

func NewSentDeclaracionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SentDeclaracionContext {
	var p = new(SentDeclaracionContext)

	InitEmptySentenciaContext(&p.SentenciaContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaContext))

	return p
}

func (s *SentDeclaracionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentDeclaracionContext) DeclaracionVariable() IDeclaracionVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDeclaracionVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDeclaracionVariableContext)
}

func (s *SentDeclaracionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitSentDeclaracion(s)

	default:
		return t.VisitChildren(s)
	}
}

type SentIfContext struct {
	SentenciaContext
}

func NewSentIfContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SentIfContext {
	var p = new(SentIfContext)

	InitEmptySentenciaContext(&p.SentenciaContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaContext))

	return p
}

func (s *SentIfContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentIfContext) SentenciaIf() ISentenciaIfContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISentenciaIfContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISentenciaIfContext)
}

func (s *SentIfContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitSentIf(s)

	default:
		return t.VisitChildren(s)
	}
}

type SentForContext struct {
	SentenciaContext
}

func NewSentForContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SentForContext {
	var p = new(SentForContext)

	InitEmptySentenciaContext(&p.SentenciaContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaContext))

	return p
}

func (s *SentForContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentForContext) SentenciaFor() ISentenciaForContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISentenciaForContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISentenciaForContext)
}

func (s *SentForContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitSentFor(s)

	default:
		return t.VisitChildren(s)
	}
}

type SentExpresionContext struct {
	SentenciaContext
}

func NewSentExpresionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SentExpresionContext {
	var p = new(SentExpresionContext)

	InitEmptySentenciaContext(&p.SentenciaContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaContext))

	return p
}

func (s *SentExpresionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentExpresionContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *SentExpresionContext) PUNTOCOMA() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPUNTOCOMA, 0)
}

func (s *SentExpresionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitSentExpresion(s)

	default:
		return t.VisitChildren(s)
	}
}

type SentContinueContext struct {
	SentenciaContext
}

func NewSentContinueContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SentContinueContext {
	var p = new(SentContinueContext)

	InitEmptySentenciaContext(&p.SentenciaContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaContext))

	return p
}

func (s *SentContinueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentContinueContext) CONTINUE() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCONTINUE, 0)
}

func (s *SentContinueContext) PUNTOCOMA() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPUNTOCOMA, 0)
}

func (s *SentContinueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitSentContinue(s)

	default:
		return t.VisitChildren(s)
	}
}

type SentBreakContext struct {
	SentenciaContext
}

func NewSentBreakContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SentBreakContext {
	var p = new(SentBreakContext)

	InitEmptySentenciaContext(&p.SentenciaContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaContext))

	return p
}

func (s *SentBreakContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentBreakContext) BREAK() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserBREAK, 0)
}

func (s *SentBreakContext) PUNTOCOMA() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPUNTOCOMA, 0)
}

func (s *SentBreakContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitSentBreak(s)

	default:
		return t.VisitChildren(s)
	}
}

type SentReturnContext struct {
	SentenciaContext
}

func NewSentReturnContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SentReturnContext {
	var p = new(SentReturnContext)

	InitEmptySentenciaContext(&p.SentenciaContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaContext))

	return p
}

func (s *SentReturnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentReturnContext) RETURN() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserRETURN, 0)
}

func (s *SentReturnContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *SentReturnContext) PUNTOCOMA() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPUNTOCOMA, 0)
}

func (s *SentReturnContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitSentReturn(s)

	default:
		return t.VisitChildren(s)
	}
}

type SentSwitchContext struct {
	SentenciaContext
}

func NewSentSwitchContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SentSwitchContext {
	var p = new(SentSwitchContext)

	InitEmptySentenciaContext(&p.SentenciaContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaContext))

	return p
}

func (s *SentSwitchContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentSwitchContext) SentenciaSwitch() ISentenciaSwitchContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISentenciaSwitchContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISentenciaSwitchContext)
}

func (s *SentSwitchContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitSentSwitch(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) Sentencia() (localctx ISentenciaContext) {
	localctx = NewSentenciaContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, VLangCherryParserRULE_sentencia)
	var _la int

	p.SetState(168)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 15, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSentDeclaracionContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(142)
			p.DeclaracionVariable()
		}

	case 2:
		localctx = NewSentAsignacionContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(143)
			p.Asignacion()
		}

	case 3:
		localctx = NewSentIncDecContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(144)
			p.IncrementoDecremento()
		}

	case 4:
		localctx = NewSentExpresionContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(145)
			p.expr(0)
		}
		p.SetState(147)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == VLangCherryParserPUNTOCOMA {
			{
				p.SetState(146)
				p.Match(VLangCherryParserPUNTOCOMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case 5:
		localctx = NewSentIfContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(149)
			p.SentenciaIf()
		}

	case 6:
		localctx = NewSentSwitchContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(150)
			p.SentenciaSwitch()
		}

	case 7:
		localctx = NewSentForContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(151)
			p.SentenciaFor()
		}

	case 8:
		localctx = NewSentBreakContext(p, localctx)
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(152)
			p.Match(VLangCherryParserBREAK)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(154)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == VLangCherryParserPUNTOCOMA {
			{
				p.SetState(153)
				p.Match(VLangCherryParserPUNTOCOMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case 9:
		localctx = NewSentContinueContext(p, localctx)
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(156)
			p.Match(VLangCherryParserCONTINUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(158)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == VLangCherryParserPUNTOCOMA {
			{
				p.SetState(157)
				p.Match(VLangCherryParserPUNTOCOMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case 10:
		localctx = NewSentReturnContext(p, localctx)
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(160)
			p.Match(VLangCherryParserRETURN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(162)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 13, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(161)
				p.expr(0)
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}
		p.SetState(165)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == VLangCherryParserPUNTOCOMA {
			{
				p.SetState(164)
				p.Match(VLangCherryParserPUNTOCOMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case 11:
		localctx = NewSentBloqueContext(p, localctx)
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(167)
			p.Bloque()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDeclaracionVariableContext is an interface to support dynamic dispatch.
type IDeclaracionVariableContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsDeclaracionVariableContext differentiates from other interfaces.
	IsDeclaracionVariableContext()
}

type DeclaracionVariableContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDeclaracionVariableContext() *DeclaracionVariableContext {
	var p = new(DeclaracionVariableContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_declaracionVariable
	return p
}

func InitEmptyDeclaracionVariableContext(p *DeclaracionVariableContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_declaracionVariable
}

func (*DeclaracionVariableContext) IsDeclaracionVariableContext() {}

func NewDeclaracionVariableContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclaracionVariableContext {
	var p = new(DeclaracionVariableContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_declaracionVariable

	return p
}

func (s *DeclaracionVariableContext) GetParser() antlr.Parser { return s.parser }

func (s *DeclaracionVariableContext) CopyAll(ctx *DeclaracionVariableContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *DeclaracionVariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeclaracionVariableContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type DeclTipadaContext struct {
	DeclaracionVariableContext
}

func NewDeclTipadaContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DeclTipadaContext {
	var p = new(DeclTipadaContext)

	InitEmptyDeclaracionVariableContext(&p.DeclaracionVariableContext)
	p.parser = parser
	p.CopyAll(ctx.(*DeclaracionVariableContext))

	return p
}

func (s *DeclTipadaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeclTipadaContext) ID() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, 0)
}

func (s *DeclTipadaContext) Tipo() ITipoContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITipoContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITipoContext)
}

func (s *DeclTipadaContext) MUT() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserMUT, 0)
}

func (s *DeclTipadaContext) ASIGNA() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserASIGNA, 0)
}

func (s *DeclTipadaContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *DeclTipadaContext) PUNTOCOMA() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPUNTOCOMA, 0)
}

func (s *DeclTipadaContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitDeclTipada(s)

	default:
		return t.VisitChildren(s)
	}
}

type DeclInferidaContext struct {
	DeclaracionVariableContext
}

func NewDeclInferidaContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DeclInferidaContext {
	var p = new(DeclInferidaContext)

	InitEmptyDeclaracionVariableContext(&p.DeclaracionVariableContext)
	p.parser = parser
	p.CopyAll(ctx.(*DeclaracionVariableContext))

	return p
}

func (s *DeclInferidaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeclInferidaContext) ID() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, 0)
}

func (s *DeclInferidaContext) ASIGNAINFERIDA() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserASIGNAINFERIDA, 0)
}

func (s *DeclInferidaContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *DeclInferidaContext) MUT() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserMUT, 0)
}

func (s *DeclInferidaContext) PUNTOCOMA() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPUNTOCOMA, 0)
}

func (s *DeclInferidaContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitDeclInferida(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) DeclaracionVariable() (localctx IDeclaracionVariableContext) {
	localctx = NewDeclaracionVariableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, VLangCherryParserRULE_declaracionVariable)
	var _la int

	p.SetState(191)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		localctx = NewDeclTipadaContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		p.SetState(171)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == VLangCherryParserMUT {
			{
				p.SetState(170)
				p.Match(VLangCherryParserMUT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(173)
			p.Match(VLangCherryParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(174)
			p.Tipo()
		}
		p.SetState(177)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == VLangCherryParserASIGNA {
			{
				p.SetState(175)
				p.Match(VLangCherryParserASIGNA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(176)
				p.expr(0)
			}

		}
		p.SetState(180)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 18, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(179)
				p.Match(VLangCherryParserPUNTOCOMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case 2:
		localctx = NewDeclInferidaContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		p.SetState(183)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == VLangCherryParserMUT {
			{
				p.SetState(182)
				p.Match(VLangCherryParserMUT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(185)
			p.Match(VLangCherryParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(186)
			p.Match(VLangCherryParserASIGNAINFERIDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(187)
			p.expr(0)
		}
		p.SetState(189)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 20, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(188)
				p.Match(VLangCherryParserPUNTOCOMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAsignacionContext is an interface to support dynamic dispatch.
type IAsignacionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOp returns the op token.
	GetOp() antlr.Token

	// SetOp sets the op token.
	SetOp(antlr.Token)

	// Getter signatures
	Lugar() ILugarContext
	Expr() IExprContext
	ASIGNA() antlr.TerminalNode
	MASIGUAL() antlr.TerminalNode
	MENOSIGUAL() antlr.TerminalNode
	PUNTOCOMA() antlr.TerminalNode

	// IsAsignacionContext differentiates from other interfaces.
	IsAsignacionContext()
}

type AsignacionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	op     antlr.Token
}

func NewEmptyAsignacionContext() *AsignacionContext {
	var p = new(AsignacionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_asignacion
	return p
}

func InitEmptyAsignacionContext(p *AsignacionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_asignacion
}

func (*AsignacionContext) IsAsignacionContext() {}

func NewAsignacionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AsignacionContext {
	var p = new(AsignacionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_asignacion

	return p
}

func (s *AsignacionContext) GetParser() antlr.Parser { return s.parser }

func (s *AsignacionContext) GetOp() antlr.Token { return s.op }

func (s *AsignacionContext) SetOp(v antlr.Token) { s.op = v }

func (s *AsignacionContext) Lugar() ILugarContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILugarContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILugarContext)
}

func (s *AsignacionContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *AsignacionContext) ASIGNA() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserASIGNA, 0)
}

func (s *AsignacionContext) MASIGUAL() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserMASIGUAL, 0)
}

func (s *AsignacionContext) MENOSIGUAL() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserMENOSIGUAL, 0)
}

func (s *AsignacionContext) PUNTOCOMA() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPUNTOCOMA, 0)
}

func (s *AsignacionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AsignacionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AsignacionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitAsignacion(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) Asignacion() (localctx IAsignacionContext) {
	localctx = NewAsignacionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, VLangCherryParserRULE_asignacion)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(193)
		p.lugar(0)
	}
	{
		p.SetState(194)

		var _lt = p.GetTokenStream().LT(1)

		localctx.(*AsignacionContext).op = _lt

		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&2199035838464) != 0) {
			var _ri = p.GetErrorHandler().RecoverInline(p)

			localctx.(*AsignacionContext).op = _ri
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(195)
		p.expr(0)
	}
	p.SetState(197)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 22, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(196)
			p.Match(VLangCherryParserPUNTOCOMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILugarContext is an interface to support dynamic dispatch.
type ILugarContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsLugarContext differentiates from other interfaces.
	IsLugarContext()
}

type LugarContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLugarContext() *LugarContext {
	var p = new(LugarContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_lugar
	return p
}

func InitEmptyLugarContext(p *LugarContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_lugar
}

func (*LugarContext) IsLugarContext() {}

func NewLugarContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LugarContext {
	var p = new(LugarContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_lugar

	return p
}

func (s *LugarContext) GetParser() antlr.Parser { return s.parser }

func (s *LugarContext) CopyAll(ctx *LugarContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *LugarContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LugarContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type LugarCampoContext struct {
	LugarContext
}

func NewLugarCampoContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LugarCampoContext {
	var p = new(LugarCampoContext)

	InitEmptyLugarContext(&p.LugarContext)
	p.parser = parser
	p.CopyAll(ctx.(*LugarContext))

	return p
}

func (s *LugarCampoContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LugarCampoContext) Lugar() ILugarContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILugarContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILugarContext)
}

func (s *LugarCampoContext) PUNTO() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPUNTO, 0)
}

func (s *LugarCampoContext) ID() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, 0)
}

func (s *LugarCampoContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitLugarCampo(s)

	default:
		return t.VisitChildren(s)
	}
}

type LugarIndiceContext struct {
	LugarContext
}

func NewLugarIndiceContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LugarIndiceContext {
	var p = new(LugarIndiceContext)

	InitEmptyLugarContext(&p.LugarContext)
	p.parser = parser
	p.CopyAll(ctx.(*LugarContext))

	return p
}

func (s *LugarIndiceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LugarIndiceContext) Lugar() ILugarContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILugarContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILugarContext)
}

func (s *LugarIndiceContext) CORIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCORIZQ, 0)
}

func (s *LugarIndiceContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *LugarIndiceContext) CORDER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCORDER, 0)
}

func (s *LugarIndiceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitLugarIndice(s)

	default:
		return t.VisitChildren(s)
	}
}

type LugarIdContext struct {
	LugarContext
}

func NewLugarIdContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LugarIdContext {
	var p = new(LugarIdContext)

	InitEmptyLugarContext(&p.LugarContext)
	p.parser = parser
	p.CopyAll(ctx.(*LugarContext))

	return p
}

func (s *LugarIdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LugarIdContext) ID() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, 0)
}

func (s *LugarIdContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitLugarId(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) Lugar() (localctx ILugarContext) {
	return p.lugar(0)
}

func (p *VLangCherryParser) lugar(_p int) (localctx ILugarContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewLugarContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx ILugarContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 30
	p.EnterRecursionRule(localctx, 30, VLangCherryParserRULE_lugar, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	localctx = NewLugarIdContext(p, localctx)
	p.SetParserRuleContext(localctx)
	_prevctx = localctx

	{
		p.SetState(200)
		p.Match(VLangCherryParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(212)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 24, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(210)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 23, p.GetParserRuleContext()) {
			case 1:
				localctx = NewLugarIndiceContext(p, NewLugarContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, VLangCherryParserRULE_lugar)
				p.SetState(202)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(203)
					p.Match(VLangCherryParserCORIZQ)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(204)
					p.expr(0)
				}
				{
					p.SetState(205)
					p.Match(VLangCherryParserCORDER)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case 2:
				localctx = NewLugarCampoContext(p, NewLugarContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, VLangCherryParserRULE_lugar)
				p.SetState(207)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
					goto errorExit
				}
				{
					p.SetState(208)
					p.Match(VLangCherryParserPUNTO)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(209)
					p.Match(VLangCherryParserID)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(214)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 24, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIncrementoDecrementoContext is an interface to support dynamic dispatch.
type IIncrementoDecrementoContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOp returns the op token.
	GetOp() antlr.Token

	// SetOp sets the op token.
	SetOp(antlr.Token)

	// Getter signatures
	Lugar() ILugarContext
	INCREMENTO() antlr.TerminalNode
	DECREMENTO() antlr.TerminalNode
	PUNTOCOMA() antlr.TerminalNode

	// IsIncrementoDecrementoContext differentiates from other interfaces.
	IsIncrementoDecrementoContext()
}

type IncrementoDecrementoContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	op     antlr.Token
}

func NewEmptyIncrementoDecrementoContext() *IncrementoDecrementoContext {
	var p = new(IncrementoDecrementoContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_incrementoDecremento
	return p
}

func InitEmptyIncrementoDecrementoContext(p *IncrementoDecrementoContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_incrementoDecremento
}

func (*IncrementoDecrementoContext) IsIncrementoDecrementoContext() {}

func NewIncrementoDecrementoContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IncrementoDecrementoContext {
	var p = new(IncrementoDecrementoContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_incrementoDecremento

	return p
}

func (s *IncrementoDecrementoContext) GetParser() antlr.Parser { return s.parser }

func (s *IncrementoDecrementoContext) GetOp() antlr.Token { return s.op }

func (s *IncrementoDecrementoContext) SetOp(v antlr.Token) { s.op = v }

func (s *IncrementoDecrementoContext) Lugar() ILugarContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILugarContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILugarContext)
}

func (s *IncrementoDecrementoContext) INCREMENTO() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserINCREMENTO, 0)
}

func (s *IncrementoDecrementoContext) DECREMENTO() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserDECREMENTO, 0)
}

func (s *IncrementoDecrementoContext) PUNTOCOMA() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPUNTOCOMA, 0)
}

func (s *IncrementoDecrementoContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IncrementoDecrementoContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IncrementoDecrementoContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitIncrementoDecremento(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) IncrementoDecremento() (localctx IIncrementoDecrementoContext) {
	localctx = NewIncrementoDecrementoContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, VLangCherryParserRULE_incrementoDecremento)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(215)
		p.lugar(0)
	}
	{
		p.SetState(216)

		var _lt = p.GetTokenStream().LT(1)

		localctx.(*IncrementoDecrementoContext).op = _lt

		_la = p.GetTokenStream().LA(1)

		if !(_la == VLangCherryParserINCREMENTO || _la == VLangCherryParserDECREMENTO) {
			var _ri = p.GetErrorHandler().RecoverInline(p)

			localctx.(*IncrementoDecrementoContext).op = _ri
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	p.SetState(218)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == VLangCherryParserPUNTOCOMA {
		{
			p.SetState(217)
			p.Match(VLangCherryParserPUNTOCOMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISentenciaIfContext is an interface to support dynamic dispatch.
type ISentenciaIfContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIF() []antlr.TerminalNode
	IF(i int) antlr.TerminalNode
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	AllBloque() []IBloqueContext
	Bloque(i int) IBloqueContext
	AllELSE() []antlr.TerminalNode
	ELSE(i int) antlr.TerminalNode

	// IsSentenciaIfContext differentiates from other interfaces.
	IsSentenciaIfContext()
}

type SentenciaIfContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySentenciaIfContext() *SentenciaIfContext {
	var p = new(SentenciaIfContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_sentenciaIf
	return p
}

func InitEmptySentenciaIfContext(p *SentenciaIfContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_sentenciaIf
}

func (*SentenciaIfContext) IsSentenciaIfContext() {}

func NewSentenciaIfContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SentenciaIfContext {
	var p = new(SentenciaIfContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_sentenciaIf

	return p
}

func (s *SentenciaIfContext) GetParser() antlr.Parser { return s.parser }

func (s *SentenciaIfContext) AllIF() []antlr.TerminalNode {
	return s.GetTokens(VLangCherryParserIF)
}

func (s *SentenciaIfContext) IF(i int) antlr.TerminalNode {
	return s.GetToken(VLangCherryParserIF, i)
}

func (s *SentenciaIfContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *SentenciaIfContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *SentenciaIfContext) AllBloque() []IBloqueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IBloqueContext); ok {
			len++
		}
	}

	tst := make([]IBloqueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IBloqueContext); ok {
			tst[i] = t.(IBloqueContext)
			i++
		}
	}

	return tst
}

func (s *SentenciaIfContext) Bloque(i int) IBloqueContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBloqueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBloqueContext)
}

func (s *SentenciaIfContext) AllELSE() []antlr.TerminalNode {
	return s.GetTokens(VLangCherryParserELSE)
}

func (s *SentenciaIfContext) ELSE(i int) antlr.TerminalNode {
	return s.GetToken(VLangCherryParserELSE, i)
}

func (s *SentenciaIfContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentenciaIfContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SentenciaIfContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitSentenciaIf(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) SentenciaIf() (localctx ISentenciaIfContext) {
	localctx = NewSentenciaIfContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, VLangCherryParserRULE_sentenciaIf)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(220)
		p.Match(VLangCherryParserIF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(221)
		p.expr(0)
	}
	{
		p.SetState(222)
		p.Bloque()
	}
	p.SetState(230)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 26, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(223)
				p.Match(VLangCherryParserELSE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(224)
				p.Match(VLangCherryParserIF)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(225)
				p.expr(0)
			}
			{
				p.SetState(226)
				p.Bloque()
			}

		}
		p.SetState(232)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 26, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(235)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == VLangCherryParserELSE {
		{
			p.SetState(233)
			p.Match(VLangCherryParserELSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(234)
			p.Bloque()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISentenciaSwitchContext is an interface to support dynamic dispatch.
type ISentenciaSwitchContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SWITCH() antlr.TerminalNode
	Expr() IExprContext
	LLAIZQ() antlr.TerminalNode
	LLADER() antlr.TerminalNode
	AllCasoSwitch() []ICasoSwitchContext
	CasoSwitch(i int) ICasoSwitchContext
	DefaultSwitch() IDefaultSwitchContext

	// IsSentenciaSwitchContext differentiates from other interfaces.
	IsSentenciaSwitchContext()
}

type SentenciaSwitchContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySentenciaSwitchContext() *SentenciaSwitchContext {
	var p = new(SentenciaSwitchContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_sentenciaSwitch
	return p
}

func InitEmptySentenciaSwitchContext(p *SentenciaSwitchContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_sentenciaSwitch
}

func (*SentenciaSwitchContext) IsSentenciaSwitchContext() {}

func NewSentenciaSwitchContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SentenciaSwitchContext {
	var p = new(SentenciaSwitchContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_sentenciaSwitch

	return p
}

func (s *SentenciaSwitchContext) GetParser() antlr.Parser { return s.parser }

func (s *SentenciaSwitchContext) SWITCH() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserSWITCH, 0)
}

func (s *SentenciaSwitchContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *SentenciaSwitchContext) LLAIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserLLAIZQ, 0)
}

func (s *SentenciaSwitchContext) LLADER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserLLADER, 0)
}

func (s *SentenciaSwitchContext) AllCasoSwitch() []ICasoSwitchContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICasoSwitchContext); ok {
			len++
		}
	}

	tst := make([]ICasoSwitchContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICasoSwitchContext); ok {
			tst[i] = t.(ICasoSwitchContext)
			i++
		}
	}

	return tst
}

func (s *SentenciaSwitchContext) CasoSwitch(i int) ICasoSwitchContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICasoSwitchContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICasoSwitchContext)
}

func (s *SentenciaSwitchContext) DefaultSwitch() IDefaultSwitchContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDefaultSwitchContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDefaultSwitchContext)
}

func (s *SentenciaSwitchContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentenciaSwitchContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SentenciaSwitchContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitSentenciaSwitch(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) SentenciaSwitch() (localctx ISentenciaSwitchContext) {
	localctx = NewSentenciaSwitchContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, VLangCherryParserRULE_sentenciaSwitch)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(237)
		p.Match(VLangCherryParserSWITCH)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(238)
		p.expr(0)
	}
	{
		p.SetState(239)
		p.Match(VLangCherryParserLLAIZQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(243)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == VLangCherryParserCASE {
		{
			p.SetState(240)
			p.CasoSwitch()
		}

		p.SetState(245)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(247)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == VLangCherryParserDEFAULT {
		{
			p.SetState(246)
			p.DefaultSwitch()
		}

	}
	{
		p.SetState(249)
		p.Match(VLangCherryParserLLADER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICasoSwitchContext is an interface to support dynamic dispatch.
type ICasoSwitchContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CASE() antlr.TerminalNode
	Expr() IExprContext
	DOSPUNTOS() antlr.TerminalNode
	AllSentencia() []ISentenciaContext
	Sentencia(i int) ISentenciaContext

	// IsCasoSwitchContext differentiates from other interfaces.
	IsCasoSwitchContext()
}

type CasoSwitchContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCasoSwitchContext() *CasoSwitchContext {
	var p = new(CasoSwitchContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_casoSwitch
	return p
}

func InitEmptyCasoSwitchContext(p *CasoSwitchContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_casoSwitch
}

func (*CasoSwitchContext) IsCasoSwitchContext() {}

func NewCasoSwitchContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CasoSwitchContext {
	var p = new(CasoSwitchContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_casoSwitch

	return p
}

func (s *CasoSwitchContext) GetParser() antlr.Parser { return s.parser }

func (s *CasoSwitchContext) CASE() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCASE, 0)
}

func (s *CasoSwitchContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *CasoSwitchContext) DOSPUNTOS() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserDOSPUNTOS, 0)
}

func (s *CasoSwitchContext) AllSentencia() []ISentenciaContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISentenciaContext); ok {
			len++
		}
	}

	tst := make([]ISentenciaContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISentenciaContext); ok {
			tst[i] = t.(ISentenciaContext)
			i++
		}
	}

	return tst
}

func (s *CasoSwitchContext) Sentencia(i int) ISentenciaContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISentenciaContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISentenciaContext)
}

func (s *CasoSwitchContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CasoSwitchContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CasoSwitchContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitCasoSwitch(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) CasoSwitch() (localctx ICasoSwitchContext) {
	localctx = NewCasoSwitchContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, VLangCherryParserRULE_casoSwitch)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(251)
		p.Match(VLangCherryParserCASE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(252)
		p.expr(0)
	}
	{
		p.SetState(253)
		p.Match(VLangCherryParserDOSPUNTOS)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(257)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&139705115656452690) != 0 {
		{
			p.SetState(254)
			p.Sentencia()
		}

		p.SetState(259)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDefaultSwitchContext is an interface to support dynamic dispatch.
type IDefaultSwitchContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEFAULT() antlr.TerminalNode
	DOSPUNTOS() antlr.TerminalNode
	AllSentencia() []ISentenciaContext
	Sentencia(i int) ISentenciaContext

	// IsDefaultSwitchContext differentiates from other interfaces.
	IsDefaultSwitchContext()
}

type DefaultSwitchContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDefaultSwitchContext() *DefaultSwitchContext {
	var p = new(DefaultSwitchContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_defaultSwitch
	return p
}

func InitEmptyDefaultSwitchContext(p *DefaultSwitchContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_defaultSwitch
}

func (*DefaultSwitchContext) IsDefaultSwitchContext() {}

func NewDefaultSwitchContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DefaultSwitchContext {
	var p = new(DefaultSwitchContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_defaultSwitch

	return p
}

func (s *DefaultSwitchContext) GetParser() antlr.Parser { return s.parser }

func (s *DefaultSwitchContext) DEFAULT() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserDEFAULT, 0)
}

func (s *DefaultSwitchContext) DOSPUNTOS() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserDOSPUNTOS, 0)
}

func (s *DefaultSwitchContext) AllSentencia() []ISentenciaContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISentenciaContext); ok {
			len++
		}
	}

	tst := make([]ISentenciaContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISentenciaContext); ok {
			tst[i] = t.(ISentenciaContext)
			i++
		}
	}

	return tst
}

func (s *DefaultSwitchContext) Sentencia(i int) ISentenciaContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISentenciaContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISentenciaContext)
}

func (s *DefaultSwitchContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DefaultSwitchContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DefaultSwitchContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitDefaultSwitch(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) DefaultSwitch() (localctx IDefaultSwitchContext) {
	localctx = NewDefaultSwitchContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, VLangCherryParserRULE_defaultSwitch)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(260)
		p.Match(VLangCherryParserDEFAULT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(261)
		p.Match(VLangCherryParserDOSPUNTOS)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(265)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&139705115656452690) != 0 {
		{
			p.SetState(262)
			p.Sentencia()
		}

		p.SetState(267)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISentenciaForContext is an interface to support dynamic dispatch.
type ISentenciaForContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsSentenciaForContext differentiates from other interfaces.
	IsSentenciaForContext()
}

type SentenciaForContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySentenciaForContext() *SentenciaForContext {
	var p = new(SentenciaForContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_sentenciaFor
	return p
}

func InitEmptySentenciaForContext(p *SentenciaForContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_sentenciaFor
}

func (*SentenciaForContext) IsSentenciaForContext() {}

func NewSentenciaForContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SentenciaForContext {
	var p = new(SentenciaForContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_sentenciaFor

	return p
}

func (s *SentenciaForContext) GetParser() antlr.Parser { return s.parser }

func (s *SentenciaForContext) CopyAll(ctx *SentenciaForContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *SentenciaForContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SentenciaForContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ForCondicionContext struct {
	SentenciaForContext
}

func NewForCondicionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ForCondicionContext {
	var p = new(ForCondicionContext)

	InitEmptySentenciaForContext(&p.SentenciaForContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaForContext))

	return p
}

func (s *ForCondicionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForCondicionContext) FOR() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserFOR, 0)
}

func (s *ForCondicionContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ForCondicionContext) Bloque() IBloqueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBloqueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBloqueContext)
}

func (s *ForCondicionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitForCondicion(s)

	default:
		return t.VisitChildren(s)
	}
}

type ForRangoContext struct {
	SentenciaForContext
}

func NewForRangoContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ForRangoContext {
	var p = new(ForRangoContext)

	InitEmptySentenciaForContext(&p.SentenciaForContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaForContext))

	return p
}

func (s *ForRangoContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForRangoContext) FOR() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserFOR, 0)
}

func (s *ForRangoContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(VLangCherryParserID)
}

func (s *ForRangoContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, i)
}

func (s *ForRangoContext) COMA() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCOMA, 0)
}

func (s *ForRangoContext) IN() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserIN, 0)
}

func (s *ForRangoContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ForRangoContext) Bloque() IBloqueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBloqueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBloqueContext)
}

func (s *ForRangoContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitForRango(s)

	default:
		return t.VisitChildren(s)
	}
}

type ForClasicoContext struct {
	SentenciaForContext
}

func NewForClasicoContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ForClasicoContext {
	var p = new(ForClasicoContext)

	InitEmptySentenciaForContext(&p.SentenciaForContext)
	p.parser = parser
	p.CopyAll(ctx.(*SentenciaForContext))

	return p
}

func (s *ForClasicoContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForClasicoContext) FOR() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserFOR, 0)
}

func (s *ForClasicoContext) AllPUNTOCOMA() []antlr.TerminalNode {
	return s.GetTokens(VLangCherryParserPUNTOCOMA)
}

func (s *ForClasicoContext) PUNTOCOMA(i int) antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPUNTOCOMA, i)
}

func (s *ForClasicoContext) Bloque() IBloqueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBloqueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBloqueContext)
}

func (s *ForClasicoContext) ForInit() IForInitContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForInitContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForInitContext)
}

func (s *ForClasicoContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ForClasicoContext) ForActualizacion() IForActualizacionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForActualizacionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForActualizacionContext)
}

func (s *ForClasicoContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitForClasico(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) SentenciaFor() (localctx ISentenciaForContext) {
	localctx = NewSentenciaForContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, VLangCherryParserRULE_sentenciaFor)
	var _la int

	p.SetState(293)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 35, p.GetParserRuleContext()) {
	case 1:
		localctx = NewForCondicionContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(268)
			p.Match(VLangCherryParserFOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(269)
			p.expr(0)
		}
		{
			p.SetState(270)
			p.Bloque()
		}

	case 2:
		localctx = NewForClasicoContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(272)
			p.Match(VLangCherryParserFOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(274)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == VLangCherryParserMUT || _la == VLangCherryParserID {
			{
				p.SetState(273)
				p.ForInit()
			}

		}
		{
			p.SetState(276)
			p.Match(VLangCherryParserPUNTOCOMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(278)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&139634746912260096) != 0 {
			{
				p.SetState(277)
				p.expr(0)
			}

		}
		{
			p.SetState(280)
			p.Match(VLangCherryParserPUNTOCOMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(282)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == VLangCherryParserID {
			{
				p.SetState(281)
				p.ForActualizacion()
			}

		}
		{
			p.SetState(284)
			p.Bloque()
		}

	case 3:
		localctx = NewForRangoContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(285)
			p.Match(VLangCherryParserFOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(286)
			p.Match(VLangCherryParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(287)
			p.Match(VLangCherryParserCOMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(288)
			p.Match(VLangCherryParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(289)
			p.Match(VLangCherryParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(290)
			p.expr(0)
		}
		{
			p.SetState(291)
			p.Bloque()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IForInitContext is an interface to support dynamic dispatch.
type IForInitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DeclaracionVariable() IDeclaracionVariableContext
	Asignacion() IAsignacionContext

	// IsForInitContext differentiates from other interfaces.
	IsForInitContext()
}

type ForInitContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForInitContext() *ForInitContext {
	var p = new(ForInitContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_forInit
	return p
}

func InitEmptyForInitContext(p *ForInitContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_forInit
}

func (*ForInitContext) IsForInitContext() {}

func NewForInitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForInitContext {
	var p = new(ForInitContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_forInit

	return p
}

func (s *ForInitContext) GetParser() antlr.Parser { return s.parser }

func (s *ForInitContext) DeclaracionVariable() IDeclaracionVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDeclaracionVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDeclaracionVariableContext)
}

func (s *ForInitContext) Asignacion() IAsignacionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAsignacionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAsignacionContext)
}

func (s *ForInitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForInitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForInitContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitForInit(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) ForInit() (localctx IForInitContext) {
	localctx = NewForInitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, VLangCherryParserRULE_forInit)
	p.SetState(297)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 36, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(295)
			p.DeclaracionVariable()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(296)
			p.Asignacion()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IForActualizacionContext is an interface to support dynamic dispatch.
type IForActualizacionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Asignacion() IAsignacionContext
	IncrementoDecremento() IIncrementoDecrementoContext

	// IsForActualizacionContext differentiates from other interfaces.
	IsForActualizacionContext()
}

type ForActualizacionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForActualizacionContext() *ForActualizacionContext {
	var p = new(ForActualizacionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_forActualizacion
	return p
}

func InitEmptyForActualizacionContext(p *ForActualizacionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_forActualizacion
}

func (*ForActualizacionContext) IsForActualizacionContext() {}

func NewForActualizacionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForActualizacionContext {
	var p = new(ForActualizacionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_forActualizacion

	return p
}

func (s *ForActualizacionContext) GetParser() antlr.Parser { return s.parser }

func (s *ForActualizacionContext) Asignacion() IAsignacionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAsignacionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAsignacionContext)
}

func (s *ForActualizacionContext) IncrementoDecremento() IIncrementoDecrementoContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIncrementoDecrementoContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIncrementoDecrementoContext)
}

func (s *ForActualizacionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForActualizacionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForActualizacionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitForActualizacion(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) ForActualizacion() (localctx IForActualizacionContext) {
	localctx = NewForActualizacionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, VLangCherryParserRULE_forActualizacion)
	p.SetState(301)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(299)
			p.Asignacion()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(300)
			p.IncrementoDecremento()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExprContext is an interface to support dynamic dispatch.
type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}

type ExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprContext() *ExprContext {
	var p = new(ExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_expr
	return p
}

func InitEmptyExprContext(p *ExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_expr
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_expr

	return p
}

func (s *ExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprContext) CopyAll(ctx *ExprContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ExprIdentificadorContext struct {
	ExprContext
}

func NewExprIdentificadorContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprIdentificadorContext {
	var p = new(ExprIdentificadorContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprIdentificadorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprIdentificadorContext) ID() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, 0)
}

func (s *ExprIdentificadorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprIdentificador(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprAditivaContext struct {
	ExprContext
	op antlr.Token
}

func NewExprAditivaContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprAditivaContext {
	var p = new(ExprAditivaContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprAditivaContext) GetOp() antlr.Token { return s.op }

func (s *ExprAditivaContext) SetOp(v antlr.Token) { s.op = v }

func (s *ExprAditivaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprAditivaContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *ExprAditivaContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprAditivaContext) MAS() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserMAS, 0)
}

func (s *ExprAditivaContext) MENOS() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserMENOS, 0)
}

func (s *ExprAditivaContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprAditiva(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprStructLitContext struct {
	ExprContext
}

func NewExprStructLitContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprStructLitContext {
	var p = new(ExprStructLitContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprStructLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprStructLitContext) LiteralStruct() ILiteralStructContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralStructContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILiteralStructContext)
}

func (s *ExprStructLitContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprStructLit(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprRelacionalContext struct {
	ExprContext
	op antlr.Token
}

func NewExprRelacionalContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprRelacionalContext {
	var p = new(ExprRelacionalContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprRelacionalContext) GetOp() antlr.Token { return s.op }

func (s *ExprRelacionalContext) SetOp(v antlr.Token) { s.op = v }

func (s *ExprRelacionalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprRelacionalContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *ExprRelacionalContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprRelacionalContext) MENOR() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserMENOR, 0)
}

func (s *ExprRelacionalContext) MENORIGUAL() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserMENORIGUAL, 0)
}

func (s *ExprRelacionalContext) MAYORIGUAL() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserMAYORIGUAL, 0)
}

func (s *ExprRelacionalContext) MAYOR() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserMAYOR, 0)
}

func (s *ExprRelacionalContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprRelacional(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprIndiceContext struct {
	ExprContext
}

func NewExprIndiceContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprIndiceContext {
	var p = new(ExprIndiceContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprIndiceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprIndiceContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *ExprIndiceContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprIndiceContext) CORIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCORIZQ, 0)
}

func (s *ExprIndiceContext) CORDER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCORDER, 0)
}

func (s *ExprIndiceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprIndice(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprParentesisContext struct {
	ExprContext
}

func NewExprParentesisContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprParentesisContext {
	var p = new(ExprParentesisContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprParentesisContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprParentesisContext) PARIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPARIZQ, 0)
}

func (s *ExprParentesisContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprParentesisContext) PARDER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPARDER, 0)
}

func (s *ExprParentesisContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprParentesis(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprSliceLitContext struct {
	ExprContext
}

func NewExprSliceLitContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprSliceLitContext {
	var p = new(ExprSliceLitContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprSliceLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprSliceLitContext) LiteralSlice() ILiteralSliceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralSliceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILiteralSliceContext)
}

func (s *ExprSliceLitContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprSliceLit(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprLiteralContext struct {
	ExprContext
}

func NewExprLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprLiteralContext {
	var p = new(ExprLiteralContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprLiteralContext) Literal() ILiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *ExprLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprOrContext struct {
	ExprContext
}

func NewExprOrContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprOrContext {
	var p = new(ExprOrContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprOrContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprOrContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *ExprOrContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprOrContext) O() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserO, 0)
}

func (s *ExprOrContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprOr(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprLlamadaContext struct {
	ExprContext
}

func NewExprLlamadaContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprLlamadaContext {
	var p = new(ExprLlamadaContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprLlamadaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprLlamadaContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprLlamadaContext) PARIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPARIZQ, 0)
}

func (s *ExprLlamadaContext) PARDER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPARDER, 0)
}

func (s *ExprLlamadaContext) ListaArgumentos() IListaArgumentosContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IListaArgumentosContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IListaArgumentosContext)
}

func (s *ExprLlamadaContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprLlamada(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprCampoContext struct {
	ExprContext
}

func NewExprCampoContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprCampoContext {
	var p = new(ExprCampoContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprCampoContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprCampoContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprCampoContext) PUNTO() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPUNTO, 0)
}

func (s *ExprCampoContext) ID() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, 0)
}

func (s *ExprCampoContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprCampo(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprIgualdadContext struct {
	ExprContext
	op antlr.Token
}

func NewExprIgualdadContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprIgualdadContext {
	var p = new(ExprIgualdadContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprIgualdadContext) GetOp() antlr.Token { return s.op }

func (s *ExprIgualdadContext) SetOp(v antlr.Token) { s.op = v }

func (s *ExprIgualdadContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprIgualdadContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *ExprIgualdadContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprIgualdadContext) IGUAL() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserIGUAL, 0)
}

func (s *ExprIgualdadContext) DIFERENTE() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserDIFERENTE, 0)
}

func (s *ExprIgualdadContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprIgualdad(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprAndContext struct {
	ExprContext
}

func NewExprAndContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprAndContext {
	var p = new(ExprAndContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprAndContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprAndContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *ExprAndContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprAndContext) Y() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserY, 0)
}

func (s *ExprAndContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprAnd(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprUnarioContext struct {
	ExprContext
	op antlr.Token
}

func NewExprUnarioContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprUnarioContext {
	var p = new(ExprUnarioContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprUnarioContext) GetOp() antlr.Token { return s.op }

func (s *ExprUnarioContext) SetOp(v antlr.Token) { s.op = v }

func (s *ExprUnarioContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprUnarioContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprUnarioContext) NOT() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserNOT, 0)
}

func (s *ExprUnarioContext) MENOS() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserMENOS, 0)
}

func (s *ExprUnarioContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprUnario(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExprMultiplicativaContext struct {
	ExprContext
	op antlr.Token
}

func NewExprMultiplicativaContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprMultiplicativaContext {
	var p = new(ExprMultiplicativaContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *ExprMultiplicativaContext) GetOp() antlr.Token { return s.op }

func (s *ExprMultiplicativaContext) SetOp(v antlr.Token) { s.op = v }

func (s *ExprMultiplicativaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprMultiplicativaContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *ExprMultiplicativaContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprMultiplicativaContext) POR() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserPOR, 0)
}

func (s *ExprMultiplicativaContext) DIV() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserDIV, 0)
}

func (s *ExprMultiplicativaContext) MODULO() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserMODULO, 0)
}

func (s *ExprMultiplicativaContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitExprMultiplicativa(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) Expr() (localctx IExprContext) {
	return p.expr(0)
}

func (p *VLangCherryParser) expr(_p int) (localctx IExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 48
	p.EnterRecursionRule(localctx, 48, VLangCherryParserRULE_expr, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(314)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 38, p.GetParserRuleContext()) {
	case 1:
		localctx = NewExprParentesisContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(304)
			p.Match(VLangCherryParserPARIZQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(305)
			p.expr(0)
		}
		{
			p.SetState(306)
			p.Match(VLangCherryParserPARDER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		localctx = NewExprUnarioContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(308)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*ExprUnarioContext).op = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == VLangCherryParserMENOS || _la == VLangCherryParserNOT) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*ExprUnarioContext).op = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(309)
			p.expr(11)
		}

	case 3:
		localctx = NewExprSliceLitContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(310)
			p.LiteralSlice()
		}

	case 4:
		localctx = NewExprStructLitContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(311)
			p.LiteralStruct()
		}

	case 5:
		localctx = NewExprLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(312)
			p.Literal()
		}

	case 6:
		localctx = NewExprIdentificadorContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(313)
			p.Match(VLangCherryParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(350)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 41, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(348)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 40, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExprMultiplicativaContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, VLangCherryParserRULE_expr)
				p.SetState(316)

				if !(p.Precpred(p.GetParserRuleContext(), 10)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 10)", ""))
					goto errorExit
				}
				{
					p.SetState(317)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*ExprMultiplicativaContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&962072674304) != 0) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*ExprMultiplicativaContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(318)
					p.expr(11)
				}

			case 2:
				localctx = NewExprAditivaContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, VLangCherryParserRULE_expr)
				p.SetState(319)

				if !(p.Precpred(p.GetParserRuleContext(), 9)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 9)", ""))
					goto errorExit
				}
				{
					p.SetState(320)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*ExprAditivaContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == VLangCherryParserMAS || _la == VLangCherryParserMENOS) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*ExprAditivaContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(321)
					p.expr(10)
				}

			case 3:
				localctx = NewExprRelacionalContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, VLangCherryParserRULE_expr)
				p.SetState(322)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
					goto errorExit
				}
				{
					p.SetState(323)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*ExprRelacionalContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&27380416512) != 0) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*ExprRelacionalContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(324)
					p.expr(9)
				}

			case 4:
				localctx = NewExprIgualdadContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, VLangCherryParserRULE_expr)
				p.SetState(325)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
					goto errorExit
				}
				{
					p.SetState(326)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*ExprIgualdadContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == VLangCherryParserIGUAL || _la == VLangCherryParserDIFERENTE) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*ExprIgualdadContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(327)
					p.expr(8)
				}

			case 5:
				localctx = NewExprAndContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, VLangCherryParserRULE_expr)
				p.SetState(328)

				if !(p.Precpred(p.GetParserRuleContext(), 6)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
					goto errorExit
				}
				{
					p.SetState(329)
					p.Match(VLangCherryParserY)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(330)
					p.expr(7)
				}

			case 6:
				localctx = NewExprOrContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, VLangCherryParserRULE_expr)
				p.SetState(331)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
					goto errorExit
				}
				{
					p.SetState(332)
					p.Match(VLangCherryParserO)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(333)
					p.expr(6)
				}

			case 7:
				localctx = NewExprLlamadaContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, VLangCherryParserRULE_expr)
				p.SetState(334)

				if !(p.Precpred(p.GetParserRuleContext(), 15)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 15)", ""))
					goto errorExit
				}
				{
					p.SetState(335)
					p.Match(VLangCherryParserPARIZQ)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				p.SetState(337)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)

				if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&139634746912260096) != 0 {
					{
						p.SetState(336)
						p.ListaArgumentos()
					}

				}
				{
					p.SetState(339)
					p.Match(VLangCherryParserPARDER)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case 8:
				localctx = NewExprIndiceContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, VLangCherryParserRULE_expr)
				p.SetState(340)

				if !(p.Precpred(p.GetParserRuleContext(), 14)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 14)", ""))
					goto errorExit
				}
				{
					p.SetState(341)
					p.Match(VLangCherryParserCORIZQ)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(342)
					p.expr(0)
				}
				{
					p.SetState(343)
					p.Match(VLangCherryParserCORDER)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case 9:
				localctx = NewExprCampoContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, VLangCherryParserRULE_expr)
				p.SetState(345)

				if !(p.Precpred(p.GetParserRuleContext(), 13)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 13)", ""))
					goto errorExit
				}
				{
					p.SetState(346)
					p.Match(VLangCherryParserPUNTO)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(347)
					p.Match(VLangCherryParserID)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(352)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 41, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IListaArgumentosContext is an interface to support dynamic dispatch.
type IListaArgumentosContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	AllCOMA() []antlr.TerminalNode
	COMA(i int) antlr.TerminalNode

	// IsListaArgumentosContext differentiates from other interfaces.
	IsListaArgumentosContext()
}

type ListaArgumentosContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyListaArgumentosContext() *ListaArgumentosContext {
	var p = new(ListaArgumentosContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_listaArgumentos
	return p
}

func InitEmptyListaArgumentosContext(p *ListaArgumentosContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_listaArgumentos
}

func (*ListaArgumentosContext) IsListaArgumentosContext() {}

func NewListaArgumentosContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ListaArgumentosContext {
	var p = new(ListaArgumentosContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_listaArgumentos

	return p
}

func (s *ListaArgumentosContext) GetParser() antlr.Parser { return s.parser }

func (s *ListaArgumentosContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *ListaArgumentosContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ListaArgumentosContext) AllCOMA() []antlr.TerminalNode {
	return s.GetTokens(VLangCherryParserCOMA)
}

func (s *ListaArgumentosContext) COMA(i int) antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCOMA, i)
}

func (s *ListaArgumentosContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ListaArgumentosContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ListaArgumentosContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitListaArgumentos(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) ListaArgumentos() (localctx IListaArgumentosContext) {
	localctx = NewListaArgumentosContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, VLangCherryParserRULE_listaArgumentos)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(353)
		p.expr(0)
	}
	p.SetState(358)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == VLangCherryParserCOMA {
		{
			p.SetState(354)
			p.Match(VLangCherryParserCOMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(355)
			p.expr(0)
		}

		p.SetState(360)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILiteralSliceContext is an interface to support dynamic dispatch.
type ILiteralSliceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CORIZQ() antlr.TerminalNode
	CORDER() antlr.TerminalNode
	Tipo() ITipoContext
	LLAIZQ() antlr.TerminalNode
	LLADER() antlr.TerminalNode
	ListaArgumentos() IListaArgumentosContext
	AllFilaSlice() []IFilaSliceContext
	FilaSlice(i int) IFilaSliceContext
	AllCOMA() []antlr.TerminalNode
	COMA(i int) antlr.TerminalNode

	// IsLiteralSliceContext differentiates from other interfaces.
	IsLiteralSliceContext()
}

type LiteralSliceContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralSliceContext() *LiteralSliceContext {
	var p = new(LiteralSliceContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_literalSlice
	return p
}

func InitEmptyLiteralSliceContext(p *LiteralSliceContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_literalSlice
}

func (*LiteralSliceContext) IsLiteralSliceContext() {}

func NewLiteralSliceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralSliceContext {
	var p = new(LiteralSliceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_literalSlice

	return p
}

func (s *LiteralSliceContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralSliceContext) CORIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCORIZQ, 0)
}

func (s *LiteralSliceContext) CORDER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCORDER, 0)
}

func (s *LiteralSliceContext) Tipo() ITipoContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITipoContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITipoContext)
}

func (s *LiteralSliceContext) LLAIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserLLAIZQ, 0)
}

func (s *LiteralSliceContext) LLADER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserLLADER, 0)
}

func (s *LiteralSliceContext) ListaArgumentos() IListaArgumentosContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IListaArgumentosContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IListaArgumentosContext)
}

func (s *LiteralSliceContext) AllFilaSlice() []IFilaSliceContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFilaSliceContext); ok {
			len++
		}
	}

	tst := make([]IFilaSliceContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFilaSliceContext); ok {
			tst[i] = t.(IFilaSliceContext)
			i++
		}
	}

	return tst
}

func (s *LiteralSliceContext) FilaSlice(i int) IFilaSliceContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilaSliceContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilaSliceContext)
}

func (s *LiteralSliceContext) AllCOMA() []antlr.TerminalNode {
	return s.GetTokens(VLangCherryParserCOMA)
}

func (s *LiteralSliceContext) COMA(i int) antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCOMA, i)
}

func (s *LiteralSliceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralSliceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LiteralSliceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitLiteralSlice(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) LiteralSlice() (localctx ILiteralSliceContext) {
	localctx = NewLiteralSliceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, VLangCherryParserRULE_literalSlice)
	var _la int

	var _alt int

	p.SetState(387)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 46, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(361)
			p.Match(VLangCherryParserCORIZQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(362)
			p.Match(VLangCherryParserCORDER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(363)
			p.Tipo()
		}
		{
			p.SetState(364)
			p.Match(VLangCherryParserLLAIZQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(366)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&139634746912260096) != 0 {
			{
				p.SetState(365)
				p.ListaArgumentos()
			}

		}
		{
			p.SetState(368)
			p.Match(VLangCherryParserLLADER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(370)
			p.Match(VLangCherryParserCORIZQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(371)
			p.Match(VLangCherryParserCORDER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(372)
			p.Tipo()
		}
		{
			p.SetState(373)
			p.Match(VLangCherryParserLLAIZQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(374)
			p.FilaSlice()
		}
		p.SetState(379)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(375)
					p.Match(VLangCherryParserCOMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(376)
					p.FilaSlice()
				}

			}
			p.SetState(381)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}
		p.SetState(383)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == VLangCherryParserCOMA {
			{
				p.SetState(382)
				p.Match(VLangCherryParserCOMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(385)
			p.Match(VLangCherryParserLLADER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilaSliceContext is an interface to support dynamic dispatch.
type IFilaSliceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LLAIZQ() antlr.TerminalNode
	LLADER() antlr.TerminalNode
	ListaArgumentos() IListaArgumentosContext

	// IsFilaSliceContext differentiates from other interfaces.
	IsFilaSliceContext()
}

type FilaSliceContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilaSliceContext() *FilaSliceContext {
	var p = new(FilaSliceContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_filaSlice
	return p
}

func InitEmptyFilaSliceContext(p *FilaSliceContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_filaSlice
}

func (*FilaSliceContext) IsFilaSliceContext() {}

func NewFilaSliceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilaSliceContext {
	var p = new(FilaSliceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_filaSlice

	return p
}

func (s *FilaSliceContext) GetParser() antlr.Parser { return s.parser }

func (s *FilaSliceContext) LLAIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserLLAIZQ, 0)
}

func (s *FilaSliceContext) LLADER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserLLADER, 0)
}

func (s *FilaSliceContext) ListaArgumentos() IListaArgumentosContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IListaArgumentosContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IListaArgumentosContext)
}

func (s *FilaSliceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilaSliceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilaSliceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitFilaSlice(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) FilaSlice() (localctx IFilaSliceContext) {
	localctx = NewFilaSliceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, VLangCherryParserRULE_filaSlice)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(389)
		p.Match(VLangCherryParserLLAIZQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(391)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&139634746912260096) != 0 {
		{
			p.SetState(390)
			p.ListaArgumentos()
		}

	}
	{
		p.SetState(393)
		p.Match(VLangCherryParserLLADER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILiteralStructContext is an interface to support dynamic dispatch.
type ILiteralStructContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	LLAIZQ() antlr.TerminalNode
	AllCampoValor() []ICampoValorContext
	CampoValor(i int) ICampoValorContext
	LLADER() antlr.TerminalNode
	AllCOMA() []antlr.TerminalNode
	COMA(i int) antlr.TerminalNode

	// IsLiteralStructContext differentiates from other interfaces.
	IsLiteralStructContext()
}

type LiteralStructContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralStructContext() *LiteralStructContext {
	var p = new(LiteralStructContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_literalStruct
	return p
}

func InitEmptyLiteralStructContext(p *LiteralStructContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_literalStruct
}

func (*LiteralStructContext) IsLiteralStructContext() {}

func NewLiteralStructContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralStructContext {
	var p = new(LiteralStructContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_literalStruct

	return p
}

func (s *LiteralStructContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralStructContext) ID() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, 0)
}

func (s *LiteralStructContext) LLAIZQ() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserLLAIZQ, 0)
}

func (s *LiteralStructContext) AllCampoValor() []ICampoValorContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICampoValorContext); ok {
			len++
		}
	}

	tst := make([]ICampoValorContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICampoValorContext); ok {
			tst[i] = t.(ICampoValorContext)
			i++
		}
	}

	return tst
}

func (s *LiteralStructContext) CampoValor(i int) ICampoValorContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICampoValorContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICampoValorContext)
}

func (s *LiteralStructContext) LLADER() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserLLADER, 0)
}

func (s *LiteralStructContext) AllCOMA() []antlr.TerminalNode {
	return s.GetTokens(VLangCherryParserCOMA)
}

func (s *LiteralStructContext) COMA(i int) antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCOMA, i)
}

func (s *LiteralStructContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralStructContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LiteralStructContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitLiteralStruct(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) LiteralStruct() (localctx ILiteralStructContext) {
	localctx = NewLiteralStructContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, VLangCherryParserRULE_literalStruct)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(395)
		p.Match(VLangCherryParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(396)
		p.Match(VLangCherryParserLLAIZQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(397)
		p.CampoValor()
	}
	p.SetState(402)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(398)
				p.Match(VLangCherryParserCOMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(399)
				p.CampoValor()
			}

		}
		p.SetState(404)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(406)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == VLangCherryParserCOMA {
		{
			p.SetState(405)
			p.Match(VLangCherryParserCOMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(408)
		p.Match(VLangCherryParserLLADER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICampoValorContext is an interface to support dynamic dispatch.
type ICampoValorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	DOSPUNTOS() antlr.TerminalNode
	Expr() IExprContext

	// IsCampoValorContext differentiates from other interfaces.
	IsCampoValorContext()
}

type CampoValorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCampoValorContext() *CampoValorContext {
	var p = new(CampoValorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_campoValor
	return p
}

func InitEmptyCampoValorContext(p *CampoValorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_campoValor
}

func (*CampoValorContext) IsCampoValorContext() {}

func NewCampoValorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CampoValorContext {
	var p = new(CampoValorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_campoValor

	return p
}

func (s *CampoValorContext) GetParser() antlr.Parser { return s.parser }

func (s *CampoValorContext) ID() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserID, 0)
}

func (s *CampoValorContext) DOSPUNTOS() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserDOSPUNTOS, 0)
}

func (s *CampoValorContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *CampoValorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CampoValorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CampoValorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitCampoValor(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) CampoValor() (localctx ICampoValorContext) {
	localctx = NewCampoValorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, VLangCherryParserRULE_campoValor)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(410)
		p.Match(VLangCherryParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(411)
		p.Match(VLangCherryParserDOSPUNTOS)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(412)
		p.expr(0)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILiteralContext is an interface to support dynamic dispatch.
type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ENTERO() antlr.TerminalNode
	DECIMAL() antlr.TerminalNode
	CADENA() antlr.TerminalNode
	RUNE() antlr.TerminalNode
	TRUE() antlr.TerminalNode
	FALSE() antlr.TerminalNode
	NIL() antlr.TerminalNode

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}

type LiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralContext() *LiteralContext {
	var p = new(LiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_literal
	return p
}

func InitEmptyLiteralContext(p *LiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = VLangCherryParserRULE_literal
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = VLangCherryParserRULE_literal

	return p
}

func (s *LiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralContext) ENTERO() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserENTERO, 0)
}

func (s *LiteralContext) DECIMAL() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserDECIMAL, 0)
}

func (s *LiteralContext) CADENA() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserCADENA, 0)
}

func (s *LiteralContext) RUNE() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserRUNE, 0)
}

func (s *LiteralContext) TRUE() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserTRUE, 0)
}

func (s *LiteralContext) FALSE() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserFALSE, 0)
}

func (s *LiteralContext) NIL() antlr.TerminalNode {
	return s.GetToken(VLangCherryParserNIL, 0)
}

func (s *LiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case VLangCherryVisitor:
		return t.VisitLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *VLangCherryParser) Literal() (localctx ILiteralContext) {
	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, VLangCherryParserRULE_literal)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(414)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&135107988821229568) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

func (p *VLangCherryParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 15:
		var t *LugarContext = nil
		if localctx != nil {
			t = localctx.(*LugarContext)
		}
		return p.Lugar_Sempred(t, predIndex)

	case 24:
		var t *ExprContext = nil
		if localctx != nil {
			t = localctx.(*ExprContext)
		}
		return p.Expr_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *VLangCherryParser) Lugar_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 2)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *VLangCherryParser) Expr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 2:
		return p.Precpred(p.GetParserRuleContext(), 10)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 9)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 8)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 7)

	case 6:
		return p.Precpred(p.GetParserRuleContext(), 6)

	case 7:
		return p.Precpred(p.GetParserRuleContext(), 5)

	case 8:
		return p.Precpred(p.GetParserRuleContext(), 15)

	case 9:
		return p.Precpred(p.GetParserRuleContext(), 14)

	case 10:
		return p.Precpred(p.GetParserRuleContext(), 13)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
