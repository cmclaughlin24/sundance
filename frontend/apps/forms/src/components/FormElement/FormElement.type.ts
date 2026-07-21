import type {
  ISubmissionValue,
  NormalizedSubmission,
} from "@/types/submission";

export interface IFormElementSubmitEvent {
  /**
   * The raw submission data. This is an array of `ISubmissionVlaue` objects, each representing a field in the form and its
   * corresponding value.
   */
  raw: ISubmissionValue[];

  /**
   * The normalized submission data. This is an object representing the submission in a normalized format.
   */
  normalized: NormalizedSubmission;
}

export interface FormElementProps {
  /**
   * The tenant ID for the form. Used to identify the owner of the form.
   */
  tenantId: string;

  /**
   * The unique identifier for the form.
   */
  formId: string;

  /**
   * The unique identifier for the version of the form. This is used to specify which version of the form
   * should be rendered and submitted.
   */
  versionId: string;

  /**
   * The raw submission data. This is an optional prop that can be used to pre-fill the form with existing submission data.
   * It should be an array of `ISubmissionValue` objects, each representing a field in the form and its corresponding value. If
   * this prop is not provided the form will be rendered with empty fields.
   */
  rawSubmission?: ISubmissionValue[];

  /**
   * Callback function that is called when the form is submitted.
   * @param event The event object containing the raw and normalized submission data.
   */
  onSubmit: (event: IFormElementSubmitEvent) => void;
}
