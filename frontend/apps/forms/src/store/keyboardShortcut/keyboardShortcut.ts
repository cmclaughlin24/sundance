export interface KeyCombination {
  key: string;
  ctrlOrMeta?: boolean;
  ctrl?: boolean;
  meta?: boolean;
  shift?: boolean;
  alt?: boolean;
}

export interface KeyboardShortcut {
  name: string;
  combination: KeyCombination;
  allowInEditable?: boolean;
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

  if (
    !shortcut.combination.ctrlOrMeta &&
    !!shortcut.combination.ctrl !== event.ctrlKey
  ) {
    return false;
  }

  if (
    !shortcut.combination.ctrlOrMeta &&
    !!shortcut.combination.meta !== event.metaKey
  ) {
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

export function isValidTarget(
  shortcut: KeyboardShortcut,
  target: EventTarget | null,
): boolean {
  if (shortcut.allowInEditable) {
    return true;
  }

  return !isEditableTarget(target) && !hasTextSelection();
}

function isEditableTarget(target: EventTarget | null): boolean {
  if (!(target instanceof HTMLElement)) {
    return false;
  }

  if (target.isContentEditable) {
    return true;
  }

  const role = target.getAttribute("role");

  if (role && ["combobox", "listbox", "option", "textbox"].includes(role)) {
    return true;
  }
  return false;
}

function hasTextSelection(): boolean {
  const sel = window.getSelection();
  return !!sel && !sel.isCollapsed && sel.toString().length > 0;
}
