import * as ArrayUtils from "./array";

export interface HasVersion {
  version: number;
}

export function sortVersions<T extends HasVersion>(
  versions: T[],
  order: "asc" | "dsc" = "dsc",
): T[] {
  if (!ArrayUtils.hasLengthGreaterThan(versions, 0)) {
    return [];
  }

  const cmp = (a: T, b: T) =>
    order === "dsc" ? b.version - a.version : a.version - b.version;

  return [...versions].sort(cmp);
}
