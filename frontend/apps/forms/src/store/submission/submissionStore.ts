import { createStore } from "zustand";
import type { ISubmissionSection, ISubmissionValue } from "@/types/submission";
import * as ArrayUtils from "@/utils/array";

interface ISectionInstanceState {
  id: string;
  values: Record<string, any>;
  errors: Record<string, string[]>;
}

export interface ISubmissionStore {
  sections: Record<string, ISectionInstanceState[]>;
  addInstance: (sectionId: string) => void;
  removeInstance: (sectionId: string, instanceId: string) => void;
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
    sections: hydrateSubmissionValues(raw),
    addInstance(sectioonId) {
      return set((s) => s);
    },
    removeInstance(sectionId, instanceId) {
      return set((s) => s);
    },
    setValue(sectionId, instanceId, elementId, value) {
      return set((s) => s);
    },
    setError(sectionId, instanceId, elementId, errors) {
      return set((s) => s);
    },
  }));
}

export function hydrateSubmissionValues(
  raw: ISubmissionSection[] | undefined,
): Record<string, ISectionInstanceState[]> {
  if (!raw) {
    return {};
  }

  const sections: Record<string, ISectionInstanceState[]> = {};

  for (const section of raw) {
    if (!ArrayUtils.hasLengthGreaterThan(sections[section.sectionId], 0)) {
      sections[section.sectionId] = [];
    }

    for (const instance of section.instances) {
      sections[section.sectionId].push({
        id: crypto.randomUUID(),
        values: createElementValues(instance.values),
        errors: {},
      });
    }
  }

  return sections;
}

export function createElementValues(
  raw: ISubmissionValue[] | undefined,
): Record<string, any> {
  if (!raw) {
    return {};
  }

  return Object.fromEntries(
    raw.map(({ elementId, value }) => [elementId, value]),
  );
}
