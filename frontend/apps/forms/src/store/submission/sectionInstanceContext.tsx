import { createContext, useContext, useMemo } from "react";

const SectionInstanceContext = createContext<{
  sectionId: string;
  instanceId: string;
} | null>(null);

export function useSectionInstance() {
  const ctx = useContext(SectionInstanceContext);

  if (!ctx) {
    throw new Error(
      "useSectionInstance must be used within SectionInstanceProvider",
    );
  }

  return ctx;
}

export const SectionInstanceProvider: React.FC<
  React.PropsWithChildren<{
    sectionId: string;
    instanceId: string;
  }>
> = function ({ sectionId, instanceId, children }) {
  const value = useMemo(
    () => ({ sectionId, instanceId }),
    [sectionId, instanceId],
  );

  return (
    <SectionInstanceContext value={value}>{children}</SectionInstanceContext>
  );
};
