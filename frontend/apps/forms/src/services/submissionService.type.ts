import type { ISubmissionSection } from "@/types/submission";

export interface SubmitRequest {
  formId: string;
  versionId: string;
  sections: ISubmissionSection[];
}

export interface AsyncSubmitResponse {
  id: string;
  tenantId: string;
  formId: string;
  versionId: string;
  referenceId: string;
  status: string;
  sections: ISubmissionSection[];
  createdAt: Date;
  updatedAt: Date;
}
