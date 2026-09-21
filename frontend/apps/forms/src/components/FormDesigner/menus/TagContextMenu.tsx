import {
  useFormDesignerDispatch,
  useFormDesignerHistory,
  useFormDesignerSelect,
  useFormSnapshot,
} from "@/store/formDesigner";
import type { ITagVersion } from "@/types/tag";
import Divider from "@mui/material/Divider";
import { contextMenuStyles as styles } from "./ContextMenu.style";
import { ContextMenu } from "@/components/ContextMenu";
import Typography from "@mui/material/Typography";
import { TENANT_ID } from "@/constants/tenant";
import { useFormsService } from "@/hooks/useHttpService";
import { isDraftVersion, versionToRequest } from "@/utils/form";

export type TagContextMenuData = TagVersionContextData;

export interface TagVersionContextData {
  type: "version";
  version: ITagVersion;
}

export const TagContextMenu: React.FC<{ target: TagContextMenuData }> =
  function ({ target }) {
    const { undo, canUndo, redo, canRedo, commit } = useFormDesignerHistory();
    const { selected } = useFormDesignerSelect();
    const { dispatch } = useFormDesignerDispatch();
    const { form, version, rules } = useFormSnapshot();
    const formsService = useFormsService();

    const handleLink = () => {
      if (target.type !== "version") {
        return;
      }

      if (!selected || selected.type !== "element") {
        return;
      }

      dispatch({
        type: "AddElementTag",
        id: selected.item.id,
        mapping: {
          tagVersionId: target.version.id,
          priority: 0,
          hasStaticValue: false,
          staticValue: null,
        },
      });
    };

    const handleUnlink = () => {
      if (target.type !== "version") {
        return;
      }

      if (!selected || selected.type !== "element") {
        return;
      }

      dispatch({
        type: "RemoveElementTag",
        id: selected.item.id,
        tagVersionId: target.version.id,
      });
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

    const hasLink =
      target.type === "version" &&
      selected &&
      selected.type === "element" &&
      selected.item.tags?.some((t) => t.tagVersionId === target.version.id);

    return (
      <>
        {!hasLink && (
          <ContextMenu.Button onClick={handleLink}>Link</ContextMenu.Button>
        )}
        {hasLink && (
          <ContextMenu.Button onClick={handleUnlink}>Unlink</ContextMenu.Button>
        )}
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
      </>
    );
  };
