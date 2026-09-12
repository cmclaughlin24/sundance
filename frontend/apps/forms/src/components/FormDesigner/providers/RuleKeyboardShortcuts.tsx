import {
  useFormDesignerDispatch,
  useFormDesignerHistory,
} from "@/store/formDesigner";
import { KeyboardShortcutProvider } from "@/store/keyboardShortcut/KeyboardShortcutProvider";
import { useKeyboardShortcut } from "@/store/keyboardShortcut/useKeyboardShortcut";
import type { RuleClipboardData } from "@/types/clipboard";

export const RuleKeyboardShortcuts: React.FC<React.PropsWithChildren<{}>> =
  function ({ children }) {
    return (
      <KeyboardShortcutProvider>
        <Component>{children}</Component>
      </KeyboardShortcutProvider>
    );
  };

const Component: React.FC<React.PropsWithChildren<{}>> = function ({
  children,
}) {
  const { dispatch } = useFormDesignerDispatch();
  const { undo, redo } = useFormDesignerHistory();

  useKeyboardShortcut(
    {
      name: "Paste",
      combination: { key: "v", ctrlOrMeta: true },
      action: async () => {
        try {
          const text = await navigator.clipboard.readText();
          const data: RuleClipboardData = JSON.parse(text);
          dispatch({ type: "PasteRule", rule: data.rule });
        } catch {
          return;
        }
      },
    },
    [dispatch],
  );

  useKeyboardShortcut(
    { name: "Undo", combination: { key: "z", ctrlOrMeta: true }, action: undo },
    [undo],
  );

  useKeyboardShortcut(
    {
      name: "Redo",
      combination: { key: "z", ctrlOrMeta: true, shift: true },
      action: redo,
    },
    [redo],
  );

  return children;
};
