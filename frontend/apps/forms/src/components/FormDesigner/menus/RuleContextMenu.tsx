import {
  useFormDesignerDispatch,
  useFormDesignerHistory,
  useFormSnapshot,
} from "@/store/formDesigner";
import { ContextMenu, useContextMenuDispatch } from "../../ContextMenu";
import Typography from "@mui/material/Typography";
import Divider from "@mui/material/Divider";
import type { IFlatRule } from "@/types/rule";
import { useFormsService } from "@/hooks/useHttpService";
import { ClipboardEventType, type RuleClipboardData } from "@/types/clipboard";
import { useEffect, useState } from "react";
import { isDraftVersion, versionToRequest } from "@/utils/form";
import { TENANT_ID } from "@/constants/tenant";
import { contextMenuStyles as styles } from "./ContextMenu.style";


export const RuleContextMenu: React.FC<{ target: IFlatRule | undefined }> =
  function ({ target }) {
    const { undo, canUndo, redo, canRedo, commit } = useFormDesignerHistory();
    const { dispatch } = useFormDesignerDispatch();
    const { form, version, rules } = useFormSnapshot();
    const formsService = useFormsService();
    const { close } = useContextMenuDispatch();
    const [clipboardData, setClipboardData] =
      useState<RuleClipboardData | null>(null);

    useEffect(() => {
      navigator.clipboard
        .readText()
        .then((text) => {
          try {
            setClipboardData(JSON.parse(text) as RuleClipboardData);
          } catch {
            setClipboardData(null);
          }
        })
        .catch(() => {
          setClipboardData(null);
        });
    }, [target]);

    const handleCopy = () => {
      if (!target) {
        return;
      }

      const data: RuleClipboardData = {
        type: ClipboardEventType.CopyRule,
        rule: target,
      };
      navigator.clipboard.writeText(JSON.stringify(data));
      close();
    };

    const handlePaste = async () => {
      if (!clipboardData) {
        return;
      }

      dispatch({ type: "PasteRule", rule: clipboardData.rule });
      close();
    };

    const handleSaveDraft = async () => {
      if (!isDraftVersion(version.status)) {
        throw new Error("cannot update a non-draft-version");
      }

      const request = versionToRequest(version, rules);
      const updated = await formsService.updateFormVersion(
        form.id,
        version.id,
        request,
        {
          tenantId: TENANT_ID,
          token: "placeholder",
        },
      );
      commit(updated);
      close();
    };

    const handleDelete = () => {
      if (!target) {
        return;
      }

      dispatch({ type: "RemoveRule", id: target.id });
      close();
    };

    return (
      <>
        <ContextMenu.Button onClick={handleCopy}>Copy</ContextMenu.Button>
        <ContextMenu.Button
          sx={styles.btnWithShortcut}
          onClick={handlePaste}
          disabled={true}
        >
          <Typography>Paste</Typography>
          <Typography sx={styles.shortcutText}>Ctrl+v</Typography>
        </ContextMenu.Button>
        <Divider sx={{ my: 1 }} />
        <ContextMenu.Button
          sx={styles.btnWithShortcut}
          onClick={undo}
          disabled={!canUndo}
        >
          <Typography>Undo</Typography>
          <Typography sx={styles.shortcutText}>Ctrl+z</Typography>
        </ContextMenu.Button>
        <ContextMenu.Button
          sx={styles.btnWithShortcut}
          onClick={redo}
          disabled={!canRedo}
        >
          <Typography>Redo</Typography>
          <Typography sx={styles.shortcutText}>Ctrl+Shift+z</Typography>
        </ContextMenu.Button>
        <Divider sx={{ my: 1 }} />
        <ContextMenu.Button
          sx={styles.btnWithShortcut}
          onClick={handleSaveDraft}
        >
          <Typography>Save Draft</Typography>
        </ContextMenu.Button>
        <ContextMenu.Button sx={styles.btnWithShortcut} onClick={handleDelete}>
          <Typography>Delete</Typography>
          <Typography sx={styles.shortcutText}>Del</Typography>
        </ContextMenu.Button>
      </>
    );
  };
