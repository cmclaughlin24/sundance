import { type IForm } from "@/types/form";
import { BaseHttpService, type DefaultRequestOptions } from "./baseHttpService";
import { type IFormVersion } from "@/types/formVersion";
import { CONFIG } from "@/constants/config";
import type { FormVersionRequest } from "./formService.type";

export class FormsService extends BaseHttpService {
  static readonly serviceKey = "FormsService";

  constructor() {
    super(CONFIG.formsUrl);
  }

  /**
   * Gets a list of forms.
   * @param options The default request options.
   * @returns A promise that resolves to a list of forms.
   */
  async getForms(options: DefaultRequestOptions): Promise<IForm[]> {
    return this._get<IForm[]>(`/api/v1/forms`, options);
  }

  /**
   * Gets a form.
   * @param formId The ID of the form.
   * @param options The default request options.
   * @returns A promise that resolves to a form.
   */
  async getForm(
    formId: string,
    options: DefaultRequestOptions,
  ): Promise<IForm> {
    return await this._get<IForm>(`/api/v1/forms/${formId}`, options);
  }

  /**
   * Gets a form version.
   * @param formId The ID of the form.
   * @param versionId The ID of the form version.
   * @param options The default request options.
   * @returns A promise that resolves to a form version.
   */
  async getFormVersion(
    formId: string,
    versionId: string,
    options: DefaultRequestOptions,
  ): Promise<IFormVersion> {
    return await this._get<IFormVersion>(
      `/api/v1/forms/${formId}/versions/${versionId}`,
      options,
    );
  }

  /**
   * Gets a form and its version.
   * @param formId The ID of the form.
   * @param versionId The ID of the form version.
   * @param options The default request options.
   * @returns A promise that resolves to a tuple containing the form and its version.
   */
  async getFormAndVersion(
    formId: string,
    versionId: string,
    options: DefaultRequestOptions,
  ): Promise<[IForm, IFormVersion]> {
    const [form, version] = await Promise.all([
      this.getForm(formId, options),
      this.getFormVersion(formId, versionId, options),
    ]);

    return [form, version];
  }

  /**
   * Creates a new form version.
   * @param formId The ID of the form.
   * @param version The form version to create.
   * @param options The default request options.
   * @returns A promise that resolves to the created form version.
   */
  async createFormVersion(
    formId: string,
    version: FormVersionRequest,
    options: DefaultRequestOptions,
  ): Promise<IFormVersion> {
    const resp = await this._post<FormVersionRequest, IFormVersion>(
      `/api/v1/forms/${formId}/versions`,
      version,
      options,
    );

    return resp.data;
  }
}
