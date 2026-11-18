import type { ISubmissionValue } from "@/types/submission";
import * as ArrayUtils from "@/utils/array";

export class SectionInstanceState {
  private constructor(
    private readonly _id: string,
    private readonly _values: Record<string, any>,
    private readonly _errors: Record<string, string[]>,
  ) {}

  static create(values: Record<string, any> = {}): SectionInstanceState {
    return new SectionInstanceState(crypto.randomUUID(), { ...values }, {});
  }

  static fromSubmission(raw: ISubmissionValue[]): SectionInstanceState {
    let values: Record<string, any> = {};

    if (raw) {
      values = Object.fromEntries(
        raw.map(({ elementId, value }) => [elementId, value]),
      );
    }

    return new SectionInstanceState(crypto.randomUUID(), values, {});
  }

  get id(): string {
    return this._id;
  }

  getValue(elementId: string): any {
    return this._values[elementId];
  }

  getErrors(elementId: string): string[] | undefined {
    return this._errors[elementId];
  }

  withValue(elementId: string, value: any): SectionInstanceState {
    if (Object.is(this._values[elementId], value)) {
      return this;
    }

    return new SectionInstanceState(
      this._id,
      { ...this._values, [elementId]: value },
      this._errors,
    );
  }

  withErrors(elementId: string, errors: string[] | null) {
    const next = { ...this._errors };

    if (!ArrayUtils.hasLengthGreaterThan(errors, 0)) {
      delete next[elementId];
    } else {
      next[elementId] = errors!;
    }

    return new SectionInstanceState(this._id, this._values, next);
  }

  toSubmissionValues(): ISubmissionValue[] {
    const values: ISubmissionValue[] = [];

    for (const elementId in this._values) {
      const answer = this._values[elementId];

      if (answer == null) {
        continue;
      }

      if (Array.isArray(answer)) {
        values.push(
          ...answer.map((v, idx) => ({
            value: v,
            elementId,
            itemIndex: idx,
          })),
        );
        continue;
      }

      values.push({
        elementId,
        value: answer,
      });
    }

    return values;
  }
}
