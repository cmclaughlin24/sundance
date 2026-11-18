import type { ISubmissionSection } from "@/types/submission";
import { SectionInstanceState } from "./SectionInstanceState";
import * as ArrayUtils from "@/utils/array";
import type { ISubmissionStore } from "./submissionStore";

export function findSectionInstance(
  sections: Record<string, SectionInstanceState[]>,
  sectionId: string,
  instanceId: string,
): SectionInstanceState | undefined {
  const instances = sections[sectionId];

  if (!ArrayUtils.hasLengthGreaterThan(instances, 0)) {
    return undefined;
  }

  return instances.find((i) => i.id === instanceId);
}

export function hydrateSectionState(
  sections: ISubmissionSection[],
): Record<string, SectionInstanceState[]> {
  const state: Record<string, SectionInstanceState[]> = {};

  for (const section of sections) {
    state[section.sectionId] = section.instances.map((i) => {
      return SectionInstanceState.fromSubmission(i.values);
    });
  }

  return state;
}

export function replaceInstance(
  s: ISubmissionStore,
  sectionId: string,
  instanceId: string,
  fn: (i: SectionInstanceState) => SectionInstanceState,
): Partial<ISubmissionStore> | ISubmissionStore {
  const instances = s.sections[sectionId];

  if (!instances?.some((i) => i.id === instanceId)) {
    return s;
  }

  return {
    sections: {
      ...s.sections,
      [sectionId]: instances.map((i) => (i.id === instanceId ? fn(i) : i)),
    },
  };
}
