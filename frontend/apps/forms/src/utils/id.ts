const ID_PREFIX: string = "TEMP_";

export function generatedID(): string {
  return `${ID_PREFIX}_${crypto.randomUUID()}`;
}

export function isTemporaryID(id: string): boolean {
  return id.startsWith(ID_PREFIX);
}
