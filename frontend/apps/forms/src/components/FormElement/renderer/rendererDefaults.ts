import type { FormRendererSlotProps } from "../FormElement.type";

type SlotPropPath<T> = {
  [K in keyof T & string]-?: NonNullable<T[K]> extends object
    ? `${K}.${SlotPropPath<NonNullable<T[K]>>}`
    : K;
}[keyof T & string];

type SlotPropValue<T, P extends string> = P extends `${infer K}.${infer Rest}`
  ? K extends keyof T
    ? SlotPropValue<NonNullable<T[K]>, Rest>
    : never
  : P extends keyof T
    ? NonNullable<T[P]>
    : never;

export type FormRendererSlotPropPath = SlotPropPath<FormRendererSlotProps>;

const formRendererSlotPropsDefaults: Readonly<FormRendererSlotProps> = {
  buttons: {
    cancel: {
      show: true,
      label: "Cancel",
    },
    submit: {
      show: true,
      label: "Submit",
    },
  },
};

function getByPath(source: unknown, path: string): unknown {
  return path
    .split(".")
    .reduce<unknown>(
      (node, key) =>
        node == null ? undefined : (node as Record<string, unknown>)[key],
      source,
    );
}

/**
 * Returns the slot prop at `key` from `props`. If the consumer didn't set it, returns the default.
 */
export function getFormRendererSlotProp<P extends FormRendererSlotPropPath>(
  key: P,
  props: FormRendererSlotProps | undefined,
): SlotPropValue<FormRendererSlotProps, P> {
  return (getByPath(props, key) ??
    getByPath(formRendererSlotPropsDefaults, key)) as SlotPropValue<
    FormRendererSlotProps,
    P
  >;
}
