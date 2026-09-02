import { checkElementType } from "@/utils/error";
import type { ElementSettingsComponent } from "../ObjectSettingsPanel";
import { settingsStyle } from "./Settings.style";
import Box from "@mui/material/Box";
import type { NumberElementAttributes } from "@/types/elementAttributes";
import TextField from "@mui/material/TextField";
import { useDebounce } from "@/hooks/useDebounce";
import { useEffect, type ChangeEvent } from "react";
import { stringToNumber } from "@/utils/utils";

export const NumberElementSettings: ElementSettingsComponent = function ({
  element,
  onChange,
}) {
  checkElementType(element.type, "number");

  const attr = element.attributes as NumberElementAttributes;

  const {
    value: min,
    debounceValue: debounceMin,
    setValue: setMin,
  } = useDebounce(attr.min?.toString() ?? "", [attr.min]);

  const {
    value: max,
    debounceValue: debounceMax,
    setValue: setMax,
  } = useDebounce(attr.max?.toString() ?? "", [attr.max]);

  const {
    value: step,
    debounceValue: debounceStep,
    setValue: setStep,
  } = useDebounce(attr.step?.toString() ?? "", [attr.step]);

  const {
    value: defaultValue,
    debounceValue: debounceDefaultValue,
    setValue: setDefaultValue,
  } = useDebounce(attr.defaultValue?.toString() ?? "", [attr.defaultValue]);

  useEffect(() => {
    const change: Partial<NumberElementAttributes> = {
      min: stringToNumber(debounceMin),
      max: stringToNumber(debounceMax),
      step: stringToNumber(debounceStep),
      defaultValue: stringToNumber(debounceDefaultValue),
    };

    onChange(change);
  }, [debounceMin, debounceMax, debounceStep, debounceDefaultValue]);

  const handleChange = (field: keyof NumberElementAttributes) => {
    return (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      const value = event.target.value;

      switch (field) {
        case "min":
          setMin(value);
          break;
        case "max":
          setMax(value);
          break;
        case "step":
          setStep(value);
          break;
        case "defaultValue":
          setDefaultValue(value);
          break;
      }
    };
  };

  return (
    <Box sx={settingsStyle.container}>
      <Box sx={{ display: "flex", gap: 1 }}>
        <TextField
          type="number"
          label="Min"
          value={min}
          onChange={handleChange("min")}
        />
        <TextField
          type="number"
          label="Max"
          value={max}
          onChange={handleChange("max")}
        />
      </Box>
      <TextField
        type="number"
        label="Step"
        value={step}
        onChange={handleChange("step")}
      />
      <TextField
        type="number"
        label="Default Value"
        value={defaultValue}
        onChange={handleChange("defaultValue")}
      />
    </Box>
  );
};
