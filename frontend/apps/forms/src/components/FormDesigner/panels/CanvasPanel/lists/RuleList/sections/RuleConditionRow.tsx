import type { IPage } from "@/types/page";
import {
  RuleExpressionJoinOp,
  RuleExpressionOp,
  type IRuleExpression,
} from "@/types/rule";
import Box from "@mui/material/Box";
import Select, { type SelectChangeEvent } from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import TextField from "@mui/material/TextField";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import Delete from "@mui/icons-material/Delete";
import FormControl from "@mui/material/FormControl";
import { ruleConditionRowStyles } from "./RuleConditionRow.style";
import React from "react";

export interface RuleConditionRowProps {
  expression: IRuleExpression;
  index: number;
  pages: IPage[];
  onChange: (updated: IRuleExpression) => void;
  onDelete: () => void;
}

const operatorLabels: Record<RuleExpressionOp, string> = {
  [RuleExpressionOp.Equal]: "equals (=)",
  [RuleExpressionOp.NEqual]: "does not equal (≠)",
  [RuleExpressionOp.GreaterThan]: "greater than (>)",
  [RuleExpressionOp.GreaterThanEqualTo]: "greater than or equal (≥)",
  [RuleExpressionOp.LessThan]: "less than (<)",
  [RuleExpressionOp.LessThanEqualTo]: "less than or equal (≤)",
};

export const RuleConditionRow: React.FC<RuleConditionRowProps> = function ({
  expression,
  index,
  pages,
  onChange,
  onDelete,
}) {
  const elements = pages.flatMap((p) =>
    (p.sections ?? []).flatMap((s) => s.elements ?? []),
  );

  const handleJoinOpChange = (
    event: SelectChangeEvent<RuleExpressionJoinOp>,
  ) => {
    onChange({
      ...expression,
      joinWithPrevious: event.target.value as RuleExpressionJoinOp,
    });
  };

  const handleFieldChange = (event: SelectChangeEvent<string>) => {
    onChange({
      ...expression,
      source: { type: "field", key: event.target.value },
    });
  };

  const handleOperatorChange = (event: SelectChangeEvent<RuleExpressionOp>) => {
    onChange({
      ...expression,
      operator: event.target.value as RuleExpressionOp,
    });
  };

  const handleValueChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange({
      ...expression,
      value: event.target.value,
    });
  };

  return (
    <Box sx={ruleConditionRowStyles.conditionRow}>
      {index > 0 ? (
        <FormControl size="small" sx={ruleConditionRowStyles.joinOpSelect}>
          <Select
            value={expression.joinWithPrevious || RuleExpressionJoinOp.And}
            onChange={handleJoinOpChange}
          >
            <MenuItem value={RuleExpressionJoinOp.And}>AND</MenuItem>
            <MenuItem value={RuleExpressionJoinOp.Or}>OR</MenuItem>
          </Select>
        </FormControl>
      ) : null}

      <FormControl size="small" sx={ruleConditionRowStyles.fieldSelect}>
        <Select
          value={expression.source.key || ""}
          onChange={handleFieldChange}
          displayEmpty
        >
          <MenuItem value="" disabled>
            <em>Select Field...</em>
          </MenuItem>
          {elements.map((el) => (
            <MenuItem key={el.id} value={el.key}>
              {el.name || el.key} ({el.key})
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      <FormControl size="small" sx={ruleConditionRowStyles.operatorSelect}>
        <Select
          value={expression.operator || RuleExpressionOp.Equal}
          onChange={handleOperatorChange}
        >
          {Object.entries(operatorLabels).map(([op, label]) => (
            <MenuItem key={op} value={op}>
              {label}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      <TextField
        size="small"
        placeholder="Value"
        value={expression.value ?? ""}
        onChange={handleValueChange}
        sx={ruleConditionRowStyles.valueInput}
      />

      <Tooltip title="Delete Condition">
        <IconButton
          size="small"
          onClick={onDelete}
          sx={ruleConditionRowStyles.deleteButton}
          aria-label="Delete condition"
        >
          <Delete fontSize="small" />
        </IconButton>
      </Tooltip>
    </Box>
  );
};
