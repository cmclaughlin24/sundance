export interface KeyCombination {
  key: string;
  ctrlOrMeta?: boolean;
  shift?: boolean;
  alt?: boolean;
}

export interface KeyboardShortcut {
  name: string;
  combination: KeyCombination;
  action: () => void;
}

export function isMatch(
  event: KeyboardEvent,
  shortcut: KeyboardShortcut,
): boolean {
  if (event.key !== shortcut.combination.key) {
    return false;
  }

  if (shortcut.combination.ctrlOrMeta && !(event.ctrlKey || event.metaKey)) {
    return false;
  }

  if (!!shortcut.combination.shift !== event.shiftKey) {
    return false;
  }

  if (!!shortcut.combination.alt !== event.altKey) {
    return false;
  }

  return true;
}
