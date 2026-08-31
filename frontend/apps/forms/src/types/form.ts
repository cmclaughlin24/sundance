import type { HasName } from "./formObject";

export interface IForm extends HasName {
  id: string;
  tenantId: string;
  description: string;
  createdAt: Date;
  updatedAt: Date;
}
