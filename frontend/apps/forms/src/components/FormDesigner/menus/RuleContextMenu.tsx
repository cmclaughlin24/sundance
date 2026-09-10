import { useFormDesignerHistory } from "@/store/formDesigner";
import { ContextMenu } from "../../ContextMenu";
import Typography from "@mui/material/Typography";
import Divider from "@mui/material/Divider";
import type { Styles } from "@/types/styles";

const styles: Styles = {
  btnWithShortcut: {
    display: "flex",
    justifyContent: "space-between",
  },
  shortcutText: {
    fontSize: "0.75rem",
    color: "#4B4444",
  },
};

export const RuleContextMenu: React.FC<{}> = function () {
  const { undo, redo } = useFormDesignerHistory();
  // const { close } = useContextMenuDispatch();

  const handleCopy = () => {};

  const handlePaste = async () => {};

  const handleDelete = () => {};

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
      <ContextMenu.Button sx={styles.btnWithShortcut} onClick={undo}>
        <Typography>Undo</Typography>
        <Typography sx={styles.shortcutText}>Ctrl+z</Typography>
      </ContextMenu.Button>
      <ContextMenu.Button sx={styles.btnWithShortcut} onClick={redo}>
        <Typography>Redo</Typography>
        <Typography sx={styles.shortcutText}>Ctrl+Shift+z</Typography>
      </ContextMenu.Button>
      <Divider sx={{ my: 1 }} />
      <ContextMenu.Button sx={styles.btnWithShortcut} onClick={handleDelete}>
        <Typography>Delete</Typography>
        <Typography sx={styles.shortcutText}>Del</Typography>
      </ContextMenu.Button>
    </>
  );
};
