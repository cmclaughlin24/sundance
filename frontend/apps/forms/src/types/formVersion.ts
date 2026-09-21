import type { HasVersion } from "@/utils/version";
import type { IPage } from "./page";

export type FormVersionStatus = "draft" | "active" | "retired";

export interface IFormVersion extends HasVersion {
  id: string;
  formId: string;
  status: FormVersionStatus;
  publishedBy: string;
  publishedAt: Date;
  retiredBy: string;
  retiredAt: Date;
  createdAt: Date;
  updatedAt: Date;
  pages: IPage[];
}
