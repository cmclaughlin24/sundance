import type {
  ISubmissionValue,
  NormalizedSubmission,
} from "@/types/submission";
import { BaseHttpService, type DefaultRequestOptions } from "./baseHttpService";

interface SubmitRequest {
  formId: string;
  versionId: string;
  values: ISubmissionValue[];
}

export class SubmissionsService extends BaseHttpService {
  static readonly serviceKey = "SubmissionsService";

  constructor(baseURL: string) {
    super(baseURL);
  }

  /**
   * Normalizes the submission values for a specific form and version.
   * @param formId The ID of the form.
   * @param versionId The ID of the form version.
   * @param values The submission field values to normalize.
   * @param options The default request options.
   * @param A promise that resolves to the normaliized submission.
   */
  async normalize(
    formId: string,
    versionId: string,
    values: ISubmissionValue[],
    options: DefaultRequestOptions,
  ): Promise<NormalizedSubmission> {
    const payload: SubmitRequest = { formId, versionId, values };
    const resp = await this._post<SubmitRequest, NormalizedSubmission>(
      "/api/v1/submissions/normalize",
      payload,
      options,
    );

    return resp.data;
  }

  /**
   * Submits the form values for a specific form and version.
   * @param formId The ID of the form.
   * @param versionId The ID of the form version.
   * @param values The submission field values to submit.
   * @param options The default request options.
   * @returns A promise that resolves when the submission is complete.
   */
  async submit(
    formId: string,
    versionId: string,
    values: ISubmissionValue[],
    options: DefaultRequestOptions,
  ): Promise<void> {
    const payload: SubmitRequest = { formId, versionId, values };
    await this._post<SubmitRequest, void>("api/v1/submissions", payload, options);
  }
}
