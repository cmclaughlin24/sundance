import type {
  ISubmissionSection,
  NormalizedSubmission,
} from "@/types/submission";
import {
  BaseHttpService,
  type ApiCUDResponse,
  type DefaultRequestOptions,
} from "./baseHttpService";
import { CONFIG } from "@/constants/config";
import type {
  AsyncSubmitResponse,
  SubmitRequest,
} from "./submissionService.type";

export class SubmissionsService extends BaseHttpService {
  static readonly serviceKey = "SubmissionsService";

  constructor() {
    super(CONFIG.formsUrl);
  }

  /**
   * Normalizes the submission values for a specific form and version.
   * @param formId The ID of the form.
   * @param versionId The ID of the form version.
   * @param sections The submitted values, grouped by section, to normalize.
   * @param options The default request options.
   * @returns A promise that resolves to the normalized submission.
   */
  async normalize(
    formId: string,
    versionId: string,
    sections: ISubmissionSection[],
    options: DefaultRequestOptions,
  ): Promise<NormalizedSubmission> {
    const payload: SubmitRequest = { formId, versionId, sections };
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
   * @param sections The submitted values, grouped by section.
   * @param options The default request options.
   * @returns A promise that resolves to the accepted submission. Processing continues in the background; check `status` for progress.
   */
  async submit(
    formId: string,
    versionId: string,
    sections: ISubmissionSection[],
    options: DefaultRequestOptions,
  ): Promise<AsyncSubmitResponse> {
    const payload: SubmitRequest = { formId, versionId, sections };
    const resp = await this._post<
      SubmitRequest,
      ApiCUDResponse<AsyncSubmitResponse>
    >("api/v1/submissions", payload, options);
    return resp.data.data;
  }
}
