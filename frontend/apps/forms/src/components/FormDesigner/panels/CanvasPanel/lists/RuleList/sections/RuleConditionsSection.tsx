import type { IPage } from "@/types/page";
import {
  RuleExpressionJoinOp,
  RuleExpressionOp,
  type IRuleExpression,
} from "@/types/rule";
import { getFlattenedElements } from "@/utils/form";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Add from "@mui/icons-material/Add";
import { ruleConditionsSectionStyles } from "./RuleConditionsSection.style";
import { RuleConditionRow } from "./RuleConditionRow";
import { sortPositioned } from "@/utils/position";
import React from "react";
import { RuleSectionCard } from "./RuleSectionCard";

export interface RuleConditionsSectionProps {
  expressions: IRuleExpression[];
  pages: IPage[];
  onChange: (expressions: IRuleExpression[]) => void;
}

export const RuleConditionsSection: React.FC<RuleConditionsSectionProps> =
  function ({ expressions, pages, onChange }) {
    const sorted = sortPositioned(expressions ?? []);
    const elements = getFlattenedElements(pages);

    const handleAddCondition = () => {
      const defaultFieldKey = elements[0]?.key ?? "";
      const newExpression: IRuleExpression = {
        source: { type: "field", key: defaultFieldKey },
        operator: RuleExpressionOp.Equal,
        value: "",
        joinWithPrevious: RuleExpressionJoinOp.And,
        position: sorted.length,
      };

      onChange([...sorted, newExpression]);
    };

    const handleExpressionChange = (
      index: number,
      updated: IRuleExpression,
    ) => {
      const next = [...sorted];
      next[index] = { ...updated, position: index };
      onChange(next);
    };

    const handleDeleteExpression = (index: number) => {
      const next = sorted
        .filter((_, i) => i !== index)
        .map((exp, i) => ({ ...exp, position: i }));
      onChange(next);
    };

    return (
      <RuleSectionCard title="Conditions">
        <Box sx={ruleConditionsSectionStyles.content}>
          {sorted.length === 0 ? (
            <Typography sx={ruleConditionsSectionStyles.emptyState}>
              No conditions defined. This rule will always apply.
            </Typography>
          ) : (
            <Box sx={ruleConditionsSectionStyles.conditionsList}>
              {sorted.map((expression, index) => (
                <RuleConditionRow
                  key={`expr-${index}`}
                  expression={expression}
                  index={index}
                  pages={pages}
                  onChange={(updated) => handleExpressionChange(index, updated)}
                  onDelete={() => handleDeleteExpression(index)}
                />
              ))}
            </Box>
          )}
          <Button
            startIcon={<Add />}
            variant="outlined"
            onClick={handleAddCondition}
            sx={ruleConditionsSectionStyles.addButton}
          >
            Add Condition
          </Button>
        </Box>
      </RuleSectionCard>
    );
  };
