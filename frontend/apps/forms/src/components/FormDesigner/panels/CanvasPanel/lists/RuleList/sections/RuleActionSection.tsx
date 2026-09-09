import type { IPage } from "@/types/page";
import type { IFlatRule, RuleParentType, RuleType } from "@/types/rule";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Select, { type SelectChangeEvent } from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import FormControl from "@mui/material/FormControl";
import { ruleActionSectionStyles } from "./RuleActionSection.style";
import React from "react";
import { groupFormObjects } from "@/utils/form";
import type { IFormObject } from "@/types/formObject";
import { RuleSectionCard } from "./RuleSectionCard";

export interface RuleActionSectionProps {
  rule: IFlatRule;
  pages: IPage[];
  onChange: (
    changes: Partial<Pick<IFlatRule, "parentType" | "parentId">>,
  ) => void;
}

const ruleTypeLabels: Record<RuleType, string> = {
  required: "Required",
  visible: "Visible",
  readonly: "Read Only",
};

export const RuleActionSection: React.FC<RuleActionSectionProps> = function ({
  rule,
  pages,
  onChange,
}) {
  const scope: RuleParentType = rule.parentType || "element";
  const objects = groupFormObjects(pages);

  const handleScopeChange = (event: SelectChangeEvent<RuleParentType>) => {
    onChange({
      parentType: event.target.value as RuleParentType,
      parentId: undefined,
    });
  };

  const handleTargetChange = (event: SelectChangeEvent<string>) => {
    onChange({ parentId: event.target.value, parentType: scope });
  };

  const renderTargetOptions = () => {
    let options: IFormObject[] = [];

    switch (scope) {
      case "page":
        options = objects.pages;
        break;
      case "section":
        options = objects.sections;
        break;
      case "element":
        options = objects.elements;
        break;
    }

    return options.map((item) => (
      <MenuItem key={item.id} value={item.id}>
        {item.name || item.key}
      </MenuItem>
    ));
  };

  return (
    <RuleSectionCard title="Action">
      <Box sx={ruleActionSectionStyles.actionBox}>
        <FormControl size="small">
          <Select
            value={scope}
            onChange={handleScopeChange}
            sx={{ minWidth: 120 }}
          >
            <MenuItem value="page">Page</MenuItem>
            <MenuItem value="section">Section</MenuItem>
            <MenuItem value="element">Element</MenuItem>
          </Select>
        </FormControl>
        <FormControl size="small" sx={{ minWidth: 260, flex: 1 }}>
          <Select
            value={rule.parentId || ""}
            onChange={handleTargetChange}
            displayEmpty
          >
            <MenuItem value="" disabled>
              <em>
                Select {scope.charAt(0).toUpperCase() + scope.slice(1)}...
              </em>
            </MenuItem>
            {renderTargetOptions()}
          </Select>
        </FormControl>

        <Typography sx={ruleActionSectionStyles.actionText}>
          is{" "}
          <Typography
            component="span"
            sx={{ fontWeight: 600, textTransform: "uppercase" }}
          >
            {ruleTypeLabels[rule.type]}
          </Typography>{" "}
          when ...
        </Typography>
      </Box>
    </RuleSectionCard>
  );
};
