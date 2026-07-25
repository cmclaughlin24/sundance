import type {
  ISubmissionValue,
  NormalizedSubmission,
} from "@/types/submission";

export interface ISyncSubmitEvent {
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

export interface IAsyncSubmitEvent {
  /**
   * A server-generated reference ID used to track the status of the async submission.
   */
  referenceId: string;
}

/**
 * Base props shared across all `FormElement` variants. These props are common to all submission modes.
 */
interface BaseFormElementProps {
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
}

/**
 * Props for the `FormElement` component when operating in synchronous submission mode.
 * The `onSubmit` callback receives the raw and normalized submission data directly.
 */
interface SyncFormElementProps extends BaseFormElementProps {
  submitType: "sync";

  /**
   * Callback function that is called when the form is submitted.
   * @param event Receives an `ISyncSubmitEvent` containing the raw and normalized submission data.
   */
  onSubmit: (event: ISyncSubmitEvent) => void;
}

/**
 * Props for the `FormElement` component when operating in asynchronous submission mode.
 * The form submission is processed server-side and a server-generated reference ID is returned via `onSubmit`.
 */
interface AsyncFormElementProps extends BaseFormElementProps {
  submitType: "async";

  /**
   * Callback function that is called when the form is submitted.
   * @param event Receives an `IAsyncSubmitEvent` containing a server-generated reference ID for tracking the submission.
   */
  onSubmit: (event: IAsyncSubmitEvent) => void;
}

/**
 * Props for the `FormElement` component when no submission mode is specified.
 * Defaults to synchronous submission behavior.
 */
interface DefaultFormElementProps extends BaseFormElementProps {
  submitType?: never;

  /**
   * Callback function that is called when the form is submitted.
   * @param event Receives an `ISyncSubmitEvent` containing the raw and normalized submission data.
   */
  onSubmit: (event: ISyncSubmitEvent) => void;
}

export type FormElementProps =
  | SyncFormElementProps
  | AsyncFormElementProps
  | DefaultFormElementProps;
