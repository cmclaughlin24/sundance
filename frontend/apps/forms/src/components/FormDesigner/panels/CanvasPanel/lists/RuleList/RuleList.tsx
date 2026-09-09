import { DropZoneIndicator } from "@/components/DragDrop/DropZoneIndicator";
import { useFormRuleDragData } from "@/components/FormDesigner/providers/FormRuleDragProvider";
import { RuleItemDragType } from "@/components/FormDesigner/types/formDragEvent";
import type { IFlatRule, IRule } from "@/types/rule";
import type { Styles } from "@/types/styles";
import { useDroppable } from "@dnd-kit/react";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import Typography from "@mui/material/Typography";
import { AnimatePresence, motion, type Variants } from "motion/react";
import * as ArrayUtils from "@/utils/array";
import { RuleItem } from "./RuleItem";

const styles: Styles = {
  list: {
    m: 0,
    p: 0,
    display: "flex",
    flexDirection: "column",
    width: "100%",
    "> *": {
      mb: 1.5,
    },
  },
  instructionCard: {
    borderRadius: "10px",
    p: 2.5,
  },
};

const variants: Variants = {
  initial: { opacity: 0, height: 0, marginBottom: 0 },
  animate: { opacity: 1, height: "auto", marginBottom: "0.75rem" },
  exit: { opacity: 0, height: 0, marginBottom: 0 },
};

export interface RuleListProps {
  rules: (IFlatRule | IRule)[];
}

export const RuleList: React.FC<RuleListProps> = function ({ rules }) {
  const dragData = useFormRuleDragData();

  const { ref: dropRef, isDropTarget } = useDroppable({
    id: "rule-list",
    accept: RuleItemDragType.Rule,
    data: {},
  });

  const content = !ArrayUtils.hasLengthGreaterThan(rules, 0) ? (
    <RuleInstructionCard key="instruction-card" />
  ) : (
    rules.map((rule) => (
      <Box
        component={motion.li}
        variants={variants}
        initial="initial"
        animate="animate"
        exit="exit"
        transition={{ type: "spring", bounce: 0, duration: 0.3 }}
        sx={{ listStyle: "none" }}
        key={rule.id}
      >
        <RuleItem rule={rule} key={rule.id} />
      </Box>
    ))
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
