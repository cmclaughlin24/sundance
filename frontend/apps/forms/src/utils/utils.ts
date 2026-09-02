export function stringToNumber(value: string): number | undefined {
  return !isNaN(Number(value)) ? Number(value) : undefined;
}
