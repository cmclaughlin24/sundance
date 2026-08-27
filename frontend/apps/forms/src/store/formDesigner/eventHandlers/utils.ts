export function removeById<T extends { id: string }>(
  items: T[],
  id: string,
): T[] {
  return items.filter((item) => item.id !== id);
}

export function insertAtPosition<T extends { position: number }>(
  items: T[],
  item: T,
  position: number,
): T[] {
  const positioned = { ...item, position };
  return [...items, positioned].sort((a, b) => a.position - b.position);
}
