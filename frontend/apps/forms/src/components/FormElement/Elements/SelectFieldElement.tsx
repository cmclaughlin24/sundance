import MuiSelectField, { type SelectChangeEvent } from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import FormControl from "@mui/material/FormControl";
import FormHelperText from "@mui/material/FormHelperText";
import type { ElementComponent } from "../Renderer/ElementRenderer";
import type { ILookup, LookupValue } from "@/types/data";
import type { SelectElementAttributes } from "@/types/elementAttributes";
import { checkElementType } from "@/utils/error";
import { FieldElementContainer } from "../Layout/FieldElementContainer";
import { useDataSource } from "@/hooks/useDataSource";
import {
  useElementErrors,
  useElementValue,
  useSubmissionDispatch,
} from "@/store/submission/useSubmissionContext";
import { z } from "zod";

interface BaseSelectFieldElementProps {
  id: string;
  elementId: string;
  data: ILookup[];
  required?: boolean;
  readonly?: boolean;
  multiple?: boolean;
  minSelected?: number;
  maxSelected?: number;
  onChange: (event: any) => void;
}

const BaseSelectFieldElement: React.FC<BaseSelectFieldElementProps> =
  function ({
    data,
    required,
    readonly,
    multiple,
    minSelected,
    maxSelected,
    id,
    elementId,
    onChange,
  }) {
    const { setError } = useSubmissionDispatch();
    const value = useElementValue<LookupValue | LookupValue[]>(
      elementId,
      !multiple ? "" : [],
    );
    const errors = useElementErrors(elementId);

    const validationSchema = buildSelectValidationSchema({
      required: required ?? false,
      multiple: multiple ?? false,
      minSelected,
      maxSelected,
    });

    const handleChange = (
      event: SelectChangeEvent<LookupValue | LookupValue[]>,
    ) => {
      onChange(event.target.value);
    };

    const handleBlur = () => {
      const result = validationSchema.safeParse(value);
      setError(
        elementId,
        result.success ? [] : result.error.issues.map((e) => e.message),
      );
    };

    const content = data.map((lookup) => (
      <MenuItem value={lookup.value} key={`${lookup.value}=${lookup.label}`}>
        {lookup.label}
      </MenuItem>
    ));

    return (
      <FormControl error={errors?.length > 0}>
        <MuiSelectField
          id={id}
          value={value}
          required={required}
          disabled={readonly}
          multiple={multiple}
          onChange={handleChange}
          onBlur={handleBlur}
        >
          {content}
        </MuiSelectField>
        {errors && errors[0] && <FormHelperText>{errors[0]}</FormHelperText>}
      </FormControl>
    );
  };

const StaticSelectFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  const attr = element.attributes as SelectElementAttributes;
  return (
    <BaseSelectFieldElement
      id={element.id}
      elementId={element.id}
      multiple={attr.multiple}
      readonly={ruleState.readonly}
      required={ruleState.required}
      minSelected={attr.minSelected}
      maxSelected={attr.maxSelected}
      data={attr.data}
      onChange={onChange}
    />
  );
};

const DynamicSelectFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  const attr = element.attributes as SelectElementAttributes;
  const { data, isLoading } = useDataSource(attr.dataSourceRef!);

  return (
    <BaseSelectFieldElement
      id={element.id}
      elementId={element.id}
      multiple={attr.multiple}
      readonly={ruleState.readonly || isLoading}
      required={ruleState.required}
      minSelected={attr.minSelected}
      maxSelected={attr.maxSelected}
      data={data || []}
      onChange={onChange}
    />
  );
};

export const SelectFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "select");

  const attr = element.attributes as SelectElementAttributes;
  let Component = StaticSelectFieldElement;

  if (attr.dataSourceRef) {
    Component = DynamicSelectFieldElement;
  }

  return (
    <FieldElementContainer element={element}>
      <Component element={element} ruleState={ruleState} onChange={onChange} />
    </FieldElementContainer>
  );
};

function buildSelectValidationSchema(options: {
  required: boolean;
  multiple: boolean;
  minSelected?: number;
  maxSelected?: number;
}): z.ZodTypeAny {
  const item = z.union([z.string(), z.number()]);

  if (!options.multiple) {
    if (!options.required) {
      return item.optional();
    }

    return item.refine((v) => v !== "" && v != null, "This field is required");
  }

  let schema = z.array(item);
  const min = options.minSelected ?? (options.required ? 1 : null);

  if (min) {
    schema = schema.min(min, `Select at least ${min} option(s)`);
  }

  if (options.maxSelected != null) {
    schema = schema.max(
      options.maxSelected,
      `Select at most ${options.maxSelected} option(s)`,
    );
  }

  return schema;
}
