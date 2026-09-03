import { checkElementType } from "@/utils/error";
import type { ElementSettingsComponent } from "../ObjectSettingsPanel";
import { settingsStyle } from "./Settings.style";
import Box from "@mui/material/Box";
import type { NumberElementAttributes } from "@/types/elementAttributes";
import TextField from "@mui/material/TextField";
import { useDebouncedCallback } from "@/hooks/useDebouncedCallback";
import { useEffect, useState, type ChangeEvent } from "react";
import { stringToNumber } from "@/utils/utils";

type Fields = { min: string; max: string; step: string; defaultValue: string };

export const NumberElementSettings: ElementSettingsComponent = function ({
  element,
  onChange,
}) {
  checkElementType(element.type, "number");

  const attr = element.attributes as NumberElementAttributes;

  const [min, setMin] = useState(attr.min?.toString() ?? "");
  const [max, setMax] = useState(attr.max?.toString() ?? "");
  const [step, setStep] = useState(attr.step?.toString() ?? "");
  const [defaultValue, setDefaultValue] = useState(
    attr.defaultValue?.toString() ?? "",
  );

  const {
    debounced: debounceEmit,
    cancel,
    flush,
  } = useDebouncedCallback((next: Fields) => {
    onChange({
      min: stringToNumber(next.min),
      max: stringToNumber(next.max),
      step: stringToNumber(next.step),
      defaultValue: stringToNumber(next.defaultValue),
    });
  });

  useEffect(() => {
    cancel();
    setMin(attr.min?.toString() ?? "");
    setMax(attr.max?.toString() ?? "");
    setStep(attr.step?.toString() ?? "");
    setDefaultValue(attr.defaultValue?.toString() ?? "");
  }, [element]);

  const handleChange = (field: keyof Fields) => {
    return (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      const value = event.target.value;
      const next: Fields = { min, max, step, defaultValue };

      switch (field) {
        case "min":
          next.min = value;
          setMin(value);
          break;
        case "max":
          next.max = value;
          setMax(value);
          break;
        case "step":
          next.step = value;
          setStep(value);
          break;
        case "defaultValue":
          next.defaultValue = value;
          setDefaultValue(value);
          break;
      }

      debounceEmit(next);
    };
  };

  const handleBlur = () => flush({ min, max, step, defaultValue });

  return (
    <Box sx={settingsStyle.container}>
      <Box sx={{ display: "flex", gap: 1 }}>
        <TextField
          type="number"
          label="Min"
          value={min}
          onChange={handleChange("min")}
          onBlur={handleBlur}
        />
        <TextField
          type="number"
          label="Max"
          value={max}
          onChange={handleChange("max")}
          onBlur={handleBlur}
        />
      </Box>
      <TextField
        type="number"
        label="Step"
        value={step}
        onChange={handleChange("step")}
        onBlur={handleBlur}
      />
      <TextField
        type="number"
        label="Default Value"
        value={defaultValue}
        onChange={handleChange("defaultValue")}
        onBlur={handleBlur}
      />
    </Box>
  );
};
