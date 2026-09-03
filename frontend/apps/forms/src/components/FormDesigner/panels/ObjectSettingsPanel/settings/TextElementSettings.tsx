import { checkElementType } from "@/utils/error";
import type { ElementSettingsComponent } from "../ObjectSettingsPanel";
import { settingsStyle } from "./Settings.style";
import Box from "@mui/material/Box";
import type { TextElementAttributes } from "@/types/elementAttributes";
import TextField from "@mui/material/TextField";
import { useDebouncedCallback } from "@/hooks/useDebouncedCallback";
import { useEffect, useState, type ChangeEvent } from "react";
import { stringToNumber } from "@/utils/utils";

type Fields = {
  minLength: string;
  maxLength: string;
  placeholder: string;
  defaultValue: string;
};

export const TextElementSettings: ElementSettingsComponent = function ({
  element,
  onChange,
}) {
  checkElementType(element.type, "text");

  const attr = element.attributes as TextElementAttributes;

  const [minLength, setMinLength] = useState(attr.minLength?.toString() ?? "");
  const [maxLength, setMaxLength] = useState(attr.maxLength?.toString() ?? "");
  const [placeholder, setPlaceholder] = useState(attr.placeholder ?? "");
  const [defaultValue, setDefaultValue] = useState(attr.defaultValue ?? "");

  const {
    debounced: debounceEmit,
    cancel,
    flush,
  } = useDebouncedCallback((next: Fields) => {
    onChange({
      minLength: stringToNumber(next.minLength),
      maxLength: stringToNumber(next.maxLength),
      placeholder: next.placeholder || undefined,
      defaultValue: next.defaultValue || undefined,
    });
  });

  useEffect(() => {
    cancel();
    setMinLength(attr.minLength?.toString() ?? "");
    setMaxLength(attr.maxLength?.toString() ?? "");
    setPlaceholder(attr.placeholder ?? "");
    setDefaultValue(attr.defaultValue ?? "");
  }, [element]);

  const handleChange = (field: keyof Fields) => {
    return (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      const value = event.target.value;
      const next: Fields = { minLength, maxLength, placeholder, defaultValue };

      switch (field) {
        case "minLength":
          next.minLength = value;
          setMinLength(value);
          break;
        case "maxLength":
          next.maxLength = value;
          setMaxLength(value);
          break;
        case "placeholder":
          next.placeholder = value;
          setPlaceholder(value);
          break;
        case "defaultValue":
          next.defaultValue = value;
          setDefaultValue(value);
          break;
      }

      debounceEmit(next);
    };
  };

  const handleBlur = () =>
    flush({ minLength, maxLength, placeholder, defaultValue });

  return (
    <Box sx={settingsStyle.container}>
      <TextField
        type="number"
        label="Min. Length"
        value={minLength}
        onChange={handleChange("minLength")}
        onBlur={handleBlur}
      />
      <TextField
        type="number"
        label="Max Length"
        value={maxLength}
        onChange={handleChange("maxLength")}
        onBlur={handleBlur}
      />
      <TextField
        label="Placeholder"
        value={placeholder}
        onChange={handleChange("placeholder")}
        onBlur={handleBlur}
      />
      <TextField
        label="Default Value"
        value={defaultValue}
        onChange={handleChange("defaultValue")}
        onBlur={handleBlur}
      />
    </Box>
  );
};
