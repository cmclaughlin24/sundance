import { type IForm } from "@/types/form";
import { BaseHttpService, type DefaultRequestOptions } from "./baseHttpService";
import { type IFormVersion } from "@/types/formVersion";

export class FormsService extends BaseHttpService {
  static readonly serviceKey = "FormsService";

  constructor(baseURL: string) {
    super(baseURL);
  }

  /**
   * Gets a form and its version.
   * @param formId The ID of the form.
   * @param versionId The ID of the form version.
   * @param options The default request options.
   * @returns A promise that resolves to a tuple containing the form and its version.
   */
  async getForm(
    formId: string,
    versionId: string,
    options: DefaultRequestOptions,
  ) {
    const [form, version] = await Promise.all([
      this._get<IForm>(`/forms/${formId}`, options),
      this._get<IFormVersion>(
        `/forms/${formId}/versions/${versionId}`,
        options,
      ),
    ]);

    return [form, version];
  }
}
