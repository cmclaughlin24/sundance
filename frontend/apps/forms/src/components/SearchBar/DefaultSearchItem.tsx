import type { ILookup } from "@/types/data";

export const DefaultSearchItem: React.FC<{ option: ILookup }> = function ({
  option,
}) {
  return <>{option.label}</>;
};
