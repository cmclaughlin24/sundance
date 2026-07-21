import type { HasName } from "./hasName";

export interface IForm extends HasName {
  id: string;
  tenantId: string;
  description: string;
  createdAt: Date;
  updatedAt: Date;
}
