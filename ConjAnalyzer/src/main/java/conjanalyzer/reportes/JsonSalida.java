package conjanalyzer.reportes;

import java.io.File;
import java.nio.file.Files;
import java.util.LinkedHashMap;
import java.util.Map;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;

import conjanalyzer.interprete.Entorno;
import conjanalyzer.interprete.Operacion;

/**
 * Genera el JSON de salida de la seccion 5.4: por cada operacion, las leyes
 * aplicadas y el conjunto simplificado (en la misma notacion prefija del
 * lenguaje), o el string literal "No se puede simplificar la operacion" si no
 * se aplico ninguna ley.
 *
 * Formato exacto del ejemplo del enunciado:
 * {
 *   "operacion1": { "leyes": [...], "conjunto simplificado": "U & {A} {B} {C}" },
 *   "operacion2": "No se puede simplificar la operacion"
 * }
 */
public class JsonSalida {

    private static final Gson GSON = new GsonBuilder().setPrettyPrinting().disableHtmlEscaping().create();

    /** Devuelve el JSON como String (util para tests y para mostrar en la GUI). */
    public static String construir(Entorno ent) {
        LinkedHashMap<String, Object> raiz = new LinkedHashMap<>();
        for (Operacion op : ent.getOperaciones().values()) {
            if (op.simplificacion.seSimplifico) {
                Map<String, Object> detalle = new LinkedHashMap<>();
                detalle.put("leyes", op.simplificacion.leyes);
                detalle.put("conjunto simplificado", op.simplificacion.simplificado.toPrefijo());
                raiz.put(op.nombre, detalle);
            } else {
                raiz.put(op.nombre, "No se puede simplificar la operacion");
            }
        }
        return GSON.toJson(raiz);
    }

    /** Escribe simplificacion.json en la carpeta y devuelve el archivo. */
    public static File generar(Entorno ent, File carpeta) throws Exception {
        carpeta.mkdirs();
        File f = new File(carpeta, "simplificacion.json");
        Files.writeString(f.toPath(), construir(ent));
        return f;
    }
}
