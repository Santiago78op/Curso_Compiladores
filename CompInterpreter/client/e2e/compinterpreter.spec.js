import { test, expect } from '@playwright/test';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const EJEMPLO_ERRORES = fs.readFileSync(
  path.resolve(__dirname, '../../entradas/ejemplo_errores.ci'),
  'utf-8'
);

test.describe('CompInterpreter — cliente web', () => {
  test('ejecuta el ejemplo del anexo y muestra consola/símbolos/AST', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('CompInterpreter')).toBeVisible();

    await page.getByRole('button', { name: '▶ Ejecutar' }).click();

    // Sin errores: el panel debe quedar en Consola automáticamente.
    await expect(page.locator('.pestana.es-activa')).toHaveText('Consola');
    const consola = page.locator('.consola');
    await expect(consola).toContainText('X es mayor que Y y no es cero');
    await expect(consola).toContainText('El valor de i es: 4');
    await expect(consola).toContainText('La suma de 3 y 7 es: 10');

    await page.getByRole('button', { name: /Errores/ }).click();
    await expect(page.locator('.panel-vacio')).toContainText('Sin errores');

    await page.getByRole('button', { name: 'Símbolos' }).click();
    const filaX = page.locator('.tabla-reporte tbody tr', { hasText: 'x' }).first();
    await expect(filaX).toBeVisible();

    await page.getByRole('button', { name: 'AST' }).click();
    await expect(page.locator('.ast-grafo canvas')).toBeVisible();
  });

  test('reporta errores léxicos/sintácticos/semánticos y salta a la línea', async ({ page }) => {
    await page.goto('/');

    const textarea = page.locator('.editor-textarea');
    await textarea.click();
    await textarea.fill(EJEMPLO_ERRORES);

    await page.getByRole('button', { name: '▶ Ejecutar' }).click();

    // Con errores: el panel debe cambiar solo a Errores.
    await expect(page.locator('.pestana.es-activa')).toHaveText(/Errores/);

    const filas = page.locator('.tabla-reporte tbody tr');
    await expect(filas).toHaveCount(8);
    await expect(page.locator('.fila-error-léxico')).toHaveCount(1);
    await expect(page.locator('.fila-error-sintáctico')).toHaveCount(2);
    await expect(page.locator('.fila-error-semántico')).toHaveCount(5);

    // Clic en la primera fila: el gutter debe marcar la línea 5 como actual.
    await filas.first().click();
    await expect(page.locator('.editor-gutter-linea.es-actual')).toHaveText('5');
  });

  test('nuevo archivo agrega una pestaña editable', async ({ page }) => {
    await page.goto('/');
    const pestanasIniciales = await page.locator('.file-tab').count();

    await page.getByRole('button', { name: 'Nuevo' }).click();
    await expect(page.locator('.file-tab')).toHaveCount(pestanasIniciales + 1);
    await expect(page.locator('.file-tab.es-activa')).toContainText('sin-titulo-');
  });
});
