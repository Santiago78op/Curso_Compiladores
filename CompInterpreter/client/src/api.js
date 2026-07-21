const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:4000';

export async function interpretar(codigo) {
  const res = await fetch(`${API_URL}/interpretar`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ codigo }),
  });
  if (!res.ok) {
    throw new Error('El servidor respondió ' + res.status + '. ¿Está corriendo en ' + API_URL + '?');
  }
  return res.json();
}
