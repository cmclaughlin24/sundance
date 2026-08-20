/**
 * Checks if the given array has length greater than the specified number.
 * @param array - The array to check.
 * @param n - The number to compare the array's length against.
 * @returns `true` if the array's length is greater than the specified number, `false` otherwise.
 */
export function hasLengthGreaterThan<T>(array: T[] | null, n: number): boolean {
  return !!array && array.length > n;
}

/**
 * Checks if the given array has length equal to the specified number.
 * @param array - The array to check.
 * @param n - The number to compare the array's length against.
 * @returns `true` if the array's length is equal to the specified number, `false` otherwise.
 */
export function hasExactLength<T>(array: T[] | null, n: number): boolean {
  return !!array && array.length === n;
}
