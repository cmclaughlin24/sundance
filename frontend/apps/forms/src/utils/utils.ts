export function stringToNumber(value: string | undefined): number | undefined {
  return !isNaN(Number(value)) ? Number(value) : undefined;
}
