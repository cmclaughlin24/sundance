/**
 * A user's submission for a specific version of a form, grouped by section.
 */
export interface ISubmission {
  /**
   * Unique identifier of the form being submitted.
   */
  formId: string;

  /**
   * Unique identifier of the form version the submission was made against.
   */
  versionId: string;

  /**
   * The submitted values, grouped by the form section they belong to.
   */
  sections: ISubmissionSection[];
}

/**
 * The submitted values for a single form section.
 */
export interface ISubmissionSection {
  /**
   * Unique identifier of the section in the form version.
   */
  sectionId: string;

  /**
   * One entry per filled-in copy of the section. A non-repeatable section has a
   * single instance. A repeatable section has one instance per repetition, limited
   * by the section's `minInstances` and `maxInstances`.
   */
  instances: ISubmissionSectionInstance[];
}

/**
 * A single filled-in copy of a section, holding the field values entered for it.
 */
export interface ISubmissionSectionInstance {
  /**
   * The field values for this instance. Each element may appear only once,
   * unless it is a list field, in which case each entry must have a unique `itemIndex`.
   */
  values: ISubmissionValue[];
}

export interface ISubmissionValue {
  /**
   * Unique identifier for the field.
   */
  elementId: string;

  /**
   * The value of the field.
   */
  value: any;

  /**
   * The index of the value within a list field, if applicable. This is used for fields that hold
   * multiple values. If the field holds a single value, this property can be omitted or set to undefined.
   * An element's values must either all have an `itemIndex` or all omit it.
   */
  itemIndex?: number;
}

export type NormalizedSubmission = Record<string, any>;
