const EMPTY_DATE = "0001-01-01T00:00:00Z";

/**
 * Formats a date in american format (e.g., "MM DD, YYYY")
 *
 * @param date The date to format, either as a string or a Date object.
 * @param locales The locales to use for formatting. defaults to `en-US`
 * @returns The formatted date string, or an empty string if the date is the default "0001-01-01T00:00:00Z"
 */
export function formatAmerican(
  date: string | Date,
  locales: Intl.LocalesArgument = "en-US",
): string {
  if (typeof date === "string" && date === EMPTY_DATE) {
    return "";
  }

  const formatter = new Intl.DateTimeFormat(locales, {
    month: "short",
    day: "2-digit",
    year: "numeric",
  });

  return formatter.format(new Date(date));
}
