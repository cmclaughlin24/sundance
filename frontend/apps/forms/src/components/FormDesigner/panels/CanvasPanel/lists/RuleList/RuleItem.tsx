import { useFormDesignerDispatch } from "@/store/formDesigner";
import { ClipboardEventType, type RuleClipboardData } from "@/types/clipboard";
import type { IRule } from "@/types/rule";
import Card from "@mui/material/Card";
import { useId, useState } from "react";
import { RuleItemHeader } from "./RuleItemHeader";
import { AnimatePresence, motion } from "motion/react";
import Box from "@mui/material/Box";
import type { Styles } from "@/types/styles";

const styles: Styles = {
  ruleItem: {
    borderRadius: "10px",
    p: 2.5,
  },
  content: {
    overflow: "hidden",
    display: "flex",
    flexDirection: "column",
    gap: 2.5,
  },
};

export const RuleItem: React.FC<{ rule: IRule }> = function ({ rule }) {
  const { dispatch } = useFormDesignerDispatch();
  const [isCollapsed, setIsCollapsed] = useState(false);
  const contentId = useId();

  const handleCopy = () => {
    const data: RuleClipboardData = { type: ClipboardEventType.CopyRule, rule };
    navigator.clipboard.writeText(JSON.stringify(data));
  };

  const handleDelete = () => dispatch({ type: "RemoveRule", id: rule.id });

  return (
    <Card sx={styles.ruleItem}>
      <RuleItemHeader
        id={contentId}
        isCollapsed={isCollapsed}
        onCollapse={(value) => setIsCollapsed(value)}
        onCopy={handleCopy}
        onDelete={handleDelete}
      />
      <AnimatePresence initial={false}>
        {!isCollapsed && (
          <Box
            component={motion.div}
            id={contentId}
            key="content"
            initial={{ height: 0, opacity: 0, marginTop: 0 }}
            animate={{ height: "auto", opacity: 1, marginTop: "1.25rem" }}
            exit={{ height: 0, opacity: 0, marginTop: 0 }}
            transition={{ type: "spring", bounce: 0, duration: 0.4 }}
            sx={styles.content}
          ></Box>
        )}
      </AnimatePresence>
    </Card>
  );
};
