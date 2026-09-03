import type { HasID } from "@/types/formObject";

const ID_PREFIX: string = "TEMP_";

export function generatedID(): string {
  return `${ID_PREFIX}_${crypto.randomUUID()}`;
}

export function isTemporaryID(id: string): boolean {
  return id.startsWith(ID_PREFIX);
}

export function stripTemporaryID<T extends HasID>(item: T): T {
  return { ...item, id: isTemporaryID(item.id) ? undefined : item.id };
}
