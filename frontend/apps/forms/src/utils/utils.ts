export function numberToString(value: number | undefined): string {
  return value != null ? value.toString() : "";
}

export function stringToNumber(value: string | undefined): number | undefined {
  if (!value || value.trim() === "") {
    return undefined;
  }

  return !isNaN(Number(value)) ? Number(value) : undefined;
}
