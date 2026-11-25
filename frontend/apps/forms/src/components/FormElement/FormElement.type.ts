import type {
  ISubmissionSection,
  NormalizedSubmission,
} from "@/types/submission";

export interface ISyncSubmitEvent {
  /**
   * An array of `ISubmissionSection` objects, one per form section, each holding one or more
   * filled-in instances of that section.
   */
  raw: ISubmissionSection[];

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
 * Customization options for the parts of the `FormRenderer` UI that consumers can override.
 * All properties are optional. Anything left out uses the component's default.
 */
export interface FormRendererSlotProps {
  /**
   * Overrides for the action button in the form footer.
   */
  buttons?: {
    /**
     * Options for the `cancel` button, which submits the `SubmitAction.AddToCart`
     */
    cancel?: {
      /**
       * Whether to show the `add-to-cart` button. Defaults to `true`.
       */
      show?: boolean;

      /**
       * Text for the `cancel` button. Defaults to "Cancel".
       */
      label?: string;
    };
    submit?: {
      /**
       * Whether to show the `submit` button. Defaults to `true`.
       */
      show?: boolean;

      /**
       * Text for the `submit` button. Defaults to "Submit Request".
       */
      label?: string;
    };
  };
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
   * It should be an array of `ISubmissionSection` objects, each representing a section in the form and its corresponding values. If
   * this prop is not provided the form will be rendered with empty fields.
   */
  rawSubmission?: ISubmissionSection[];

  /**
   * Optional overrides for parts of teh form UI, such as the footer button labels and
   * whether the `cancel` button is shown. Anything left out uses the default.
   */
  slotProps?: FormRendererSlotProps;

  /**
   * Callback function that is called when the user cancels the form.
   */
  onCancel: () => void;
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
  SyncFormElementProps | AsyncFormElementProps | DefaultFormElementProps;
