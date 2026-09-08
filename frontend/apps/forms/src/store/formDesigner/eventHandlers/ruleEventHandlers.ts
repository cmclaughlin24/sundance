import type { IFlatRule } from "@/types/rule";
import type {
  AddRuleEvent,
  UpdateRuleEvent,
  RemoveRuleEvent,
  PasteRuleEvent,
} from "../events";
import type { IFormAggregate } from "./eventHandler";
import { createRule } from "@/factories/ruleFactory";
import { generatedID } from "@/utils/id";
import { removeById } from "./utils";

export function onAddRule(
  aggregate: IFormAggregate,
  event: AddRuleEvent,
): IFormAggregate {
  const rule: IFlatRule = createRule(event.id, event.ruleType);

  return { ...aggregate, rules: [...aggregate.rules, rule] };
}

export function onUpdateRule(
  aggregate: IFormAggregate,
  event: UpdateRuleEvent,
): IFormAggregate {
  const rules = aggregate.rules.map((rule) => {
    if (rule.id !== event.id) {
      return rule;
    }

    return { ...rule, ...event.changes };
  });

  return { ...aggregate, rules };
}

export function onRemoveRule(
  aggregate: IFormAggregate,
  event: RemoveRuleEvent,
): IFormAggregate {
  return { ...aggregate, rules: removeById(aggregate.rules, event.id) };
}

export function onPasteRule(
  aggregate: IFormAggregate,
  event: PasteRuleEvent,
): IFormAggregate {
  const rule: IFlatRule = {
    ...event.rule,
    id: generatedID(),
  };

  return { ...aggregate, rules: [...aggregate.rules, rule] };
}
