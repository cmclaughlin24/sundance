import {
  useFormDesignerDispatch,
  useFormPagesSnapshot,
} from "@/store/formDesigner";
import { ClipboardEventType, type RuleClipboardData } from "@/types/clipboard";
import type { IFlatRule, IRule } from "@/types/rule";
import Card from "@mui/material/Card";
import { useId, useState } from "react";
import { RuleItemHeader } from "./RuleItemHeader";
import { AnimatePresence, motion } from "motion/react";
import Box from "@mui/material/Box";
import { ruleItemStyles } from "./RuleItem.style";
import { RuleActionSection } from "./sections/RuleActionSection";
import { RuleConditionsSection } from "./sections/RuleConditionsSection";
import { findSelectedById } from "@/utils/form";
import type { IPage } from "@/types/page";

export interface RuleItemProps {
  rule: IFlatRule | IRule;
}

export const RuleItem: React.FC<RuleItemProps> = function ({ rule: rawRule }) {
  const rule = rawRule as IFlatRule;
  const pages = useFormPagesSnapshot();
  const { dispatch } = useFormDesignerDispatch();
  const [isCollapsed, setIsCollapsed] = useState(false);
  const contentId = useId();
  const headerTitle = createRuleTitle(rule, pages);

  const handleCopy = () => {
    const data: RuleClipboardData = { type: ClipboardEventType.CopyRule, rule };
    navigator.clipboard.writeText(JSON.stringify(data));
  };

  const handleDelete = () => dispatch({ type: "RemoveRule", id: rule.id });

  const handleActionChange = (
    changes: Partial<Pick<IFlatRule, "parentType" | "parentId">>,
  ) => {
    dispatch({
      type: "UpdateRule",
      id: rule.id,
      changes,
    });
  };

  const handleConditionsChange = (expressions: IFlatRule["expressions"]) => {
    dispatch({
      type: "UpdateRule",
      id: rule.id,
      changes: { expressions },
    });
  };

  return (
    <Card sx={ruleItemStyles.ruleItem}>
      <RuleItemHeader
        id={contentId}
        isCollapsed={isCollapsed}
        ruleType={rule.type}
        title={headerTitle}
        conditionCount={rule.expressions?.length ?? 0}
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
            sx={ruleItemStyles.content}
          >
            <RuleActionSection
              rule={rule}
              pages={pages}
              onChange={handleActionChange}
            />
            <RuleConditionsSection
              expressions={rule.expressions}
              pages={pages}
              onChange={handleConditionsChange}
            />
          </Box>
        )}
      </AnimatePresence>
    </Card>
  );
};

function createRuleTitle(rule: IFlatRule, pages: IPage[]): string | undefined {
  const targetItem = rule.parentId
    ? findSelectedById(pages, rule.parentId)
    : null;

  if (!targetItem) {
    return undefined;
  }

  const targetName = targetItem.item.name || targetItem.item.key;
  const targetTypeLabel =
    targetItem.type.charAt(0).toUpperCase() + targetItem.type.slice(1);

  return targetName ? `${targetName} (${targetTypeLabel})` : undefined;
}
