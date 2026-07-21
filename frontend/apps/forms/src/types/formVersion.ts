import type { IPage } from "./page";

export type FormVersionStatus = "draft" | "active" | "retired";

export interface IFormVersion {
  id: string;
  formId: string;
  version: number;
  status: FormVersionStatus;
  publishedBy: string;
  publishedAt: Date;
  retiredBy: string;
  retiredAt: Date;
  createdAt: Date;
  updatedAt: Date;
  pages: IPage[];
}
