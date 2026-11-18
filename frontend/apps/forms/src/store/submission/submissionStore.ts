import { createStore } from "zustand";
import type { ISubmissionSection } from "@/types/submission";
import { SectionInstanceState } from "./SectionInstanceState";
import * as ArrayUtils from "@/utils/array";
import { hydrateSectionState, replaceInstance } from "./utils";

export interface ISubmissionStore {
  sections: Record<string, SectionInstanceState[]>;
  addInstance: (sectionId: string) => void;
  removeInstance: (sectionId: string, instanceId: string) => void;
  ensureInstance: (sectionId: string) => void;
  setValue: (
    sectionId: string,
    instanceId: string,
    elementId: string,
    value: any,
  ) => void;
  setError: (
    sectionId: string,
    instanceId: string,
    elementId: string,
    errors: string[],
  ) => void;
}

export type SubmissionStoreApi = ReturnType<typeof createSubmissionStore>;

export function createSubmissionStore(raw: ISubmissionSection[] = []) {
  return createStore<ISubmissionStore>((set) => ({
    sections: hydrateSectionState(raw),
    addInstance(sectionId) {
      return set((s) => {
        const instances = s.sections[sectionId] ?? [];

        return {
          sections: {
            ...s.sections,
            [sectionId]: [...instances, SectionInstanceState.create()],
          },
        };
      });
    },
    removeInstance(sectionId, instanceId) {
      return set((s) => {
        let section = s.sections[sectionId];

        if (ArrayUtils.hasLengthGreaterThan(section, 0)) {
          section = section.filter((i) => i.id !== instanceId);
        }

        return { sections: { ...s.sections, [sectionId]: section } };
      });
    },
    ensureInstance(sectionId) {
      return set((s) => {
        let section = s.sections[sectionId];

        if (ArrayUtils.hasLengthGreaterThan(section, 0)) {
          return s;
        }

        return {
          sections: {
            ...s.sections,
            [sectionId]: [SectionInstanceState.create()],
          },
        };
      });
    },
    setValue(sectionId, instanceId, elementId, value) {
      set((s) =>
        replaceInstance(s, sectionId, instanceId, (i) =>
          i.withValue(elementId, value),
        ),
      );
    },
    setError(sectionId, instanceId, elementId, errors) {
      set((s) =>
        replaceInstance(s, sectionId, instanceId, (i) =>
          i.withErrors(elementId, errors),
        ),
      );
    },
  }));
}
