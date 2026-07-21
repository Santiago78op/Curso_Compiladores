import { test, expect } from '@playwright/test';

test.describe('VLangCherry — cliente web', () => {
  test('ejecuta el ejemplo inicial y muestra consola/símbolos/AST', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('VLangCherry')).toBeVisible();

    await page.getByRole('button', { name: '▶ Ejecutar' }).click();

    await expect(page.locator('.pestana.es-activa')).toHaveText('Consola');
    const consola = page.locator('.consola');
    await expect(consola).toContainText('Hola, soy Alice');
    await expect(consola).toContainText('impar: 1');
    await expect(consola).toContainText('suma(3,7) = 10');

    await page.getByRole('button', { name: /Errores/ }).click();
    await expect(page.locator('.panel-vacio')).toContainText('Sin errores');

    await page.getByRole('button', { name: 'Símbolos' }).click();
    const filaP = page.locator('.tabla-reporte tbody tr', { hasText: 'Persona' }).first();
    await expect(filaP).toBeVisible();

    await page.getByRole('button', { name: 'AST' }).click();
    await expect(page.locator('.ast-grafo canvas')).toBeVisible();
  });

  test('reporta errores semánticos y salta a la línea', async ({ page }) => {
    await page.goto('/');

    const codigoConErrores = `struct Persona {
    string Nombre
}

func main() {
    mut edad int = "no es un numero"
    print(persona.Apellido)
    mut x int = 10
    mut x int = 20
}
`;

    const textarea = page.locator('.editor-textarea');
    await textarea.click();
    await textarea.fill(codigoConErrores);

    await page.getByRole('button', { name: '▶ Ejecutar' }).click();

    await expect(page.locator('.pestana.es-activa')).toHaveText(/Errores/);

    const filas = page.locator('.tabla-reporte tbody tr');
    await expect(filas).toHaveCount(3);

    await filas.first().click();
    await expect(page.locator('.editor-gutter-linea.es-actual')).toHaveText('6');
  });

  test('nuevo archivo agrega una pestaña editable con extensión .vch', async ({ page }) => {
    await page.goto('/');
    const pestanasIniciales = await page.locator('.file-tab').count();

    await page.getByRole('button', { name: 'Nuevo' }).click();
    await expect(page.locator('.file-tab')).toHaveCount(pestanasIniciales + 1);
    await expect(page.locator('.file-tab.es-activa')).toContainText('sin-titulo-');
    await expect(page.locator('.file-tab.es-activa')).toContainText('.vch');
  });
});
