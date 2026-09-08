import { DropZoneIndicator } from "@/components/DragDrop/DropZoneIndicator";
import { useFormRuleDragData } from "@/components/FormDesigner/providers/FormRuleDragProvider";
import { RuleItemDragType } from "@/components/FormDesigner/types/formDragEvent";
import type { IRule } from "@/types/rule";
import type { Styles } from "@/types/styles";
import { useDroppable } from "@dnd-kit/react";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import Typography from "@mui/material/Typography";
import { AnimatePresence } from "motion/react";
import * as ArrayUtils from "@/utils/array";
import { RuleItem } from "./RuleItem";

const styles: Styles = {
  list: {
    m: 0,
    p: 0,
    display: "flex",
    flexDirection: "column",
    gap: 1.5,
    width: "100%",
  },
  instructionCard: {
    borderRadius: "10px",
    p: 2.5,
  },
};

export interface RuleListProps {
  rules: IRule[];
}

export const RuleList: React.FC<RuleListProps> = function ({ rules }) {
  const dragData = useFormRuleDragData();

  const { ref: dropRef, isDropTarget } = useDroppable({
    id: "rule-list",
    accept: RuleItemDragType.Rule,
    data: {},
  });

  let content = !ArrayUtils.hasLengthGreaterThan(rules, 0) ? (
    <RuleInstructionCard key="instruction-card" />
  ) : (
    rules.map((rule) => <RuleItem rule={rule} key={rule.id} />)
  );

  return (
    <Box component="ul" sx={styles.list} ref={dropRef}>
      <AnimatePresence initial={false}>
        {content}
        <DropZoneIndicator
          text="Drop Rule here"
          isVisible={!!dragData}
          isDropTarget={isDropTarget}
          key="drop-zone-indicator"
        />
      </AnimatePresence>
    </Box>
  );
};

function RuleInstructionCard() {
  return (
    <Card sx={styles.instructionCard}>
      <Typography variant="h5" sx={{ mb: 0.5 }}>
        How Rules Run
      </Typography>
      <Typography variant="body2">
        Each field starts from its attribute defaults (required / read-only,
        visible by default), then its rules run in order and override those
        flags, reacting to other fields by key and re-running on every change.
        The three rule types - visible, required, and read-only - are
        independent: for each flag the last matching rule wins, and a false
        result turns it back off. Visibility takes priority - a hidden field
        submits no value and skips its required rule. The portal evaluates rules
        live in the browser, but the forms service re-checks them on submit, so
        the backend is the source of truth.
      </Typography>
    </Card>
  );
}
