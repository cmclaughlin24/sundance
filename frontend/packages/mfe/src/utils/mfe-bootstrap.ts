/**
 * Parses the given path by removing the base path if it exists.
 *
 * @param path - The path to be parsed.
 * @param basePath - The base path to be removed from the given path.
 * @returns The parsed path with the base path removed, or the original path if the base path is not found.
 */
export function parsePath(path: string, basePath: string): string {
  if (!path) {
    return "";
  }

  if (!basePath) {
    return path || "";
  }

  if (!path.startsWith(basePath)) {
    return path;
  }

  return path.replace(basePath, "");
}
