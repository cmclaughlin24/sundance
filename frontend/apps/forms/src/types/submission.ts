export interface ISubmission {
  formId: string;
  versionId: string;
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
   * The index of the field in a collectino, if applicable. This is used for fields that are part of a collection
   * or array of fields. If the field is not part of a collection, this property can be omitted or set to undefined.
   */
  collectionIndex?: number;
}

export type NormalizedSubmission = Record<string, any>;
