import { checkElementType } from "@/utils/error";
import type { ElementSettingsComponent } from "../ObjectSettingsPanel";
import { settingsStyle } from "./Settings.style";
import Box from "@mui/material/Box";
import type { TextElementAttributes } from "@/types/elementAttributes";
import TextField from "@mui/material/TextField";
import { useDebounce } from "@/hooks/useDebounce";
import { useEffect, type ChangeEvent } from "react";
import { stringToNumber } from "@/utils/utils";

export const TextElementSettings: ElementSettingsComponent = function ({
  element,
  onChange,
}) {
  checkElementType(element.type, "text");

  const attr = element.attributes as TextElementAttributes;

  const {
    value: minLength,
    debounceValue: debounceMinLength,
    setValue: setMinLength,
  } = useDebounce(attr.minLength?.toString() ?? "", [attr.minLength]);

  const {
    value: maxLength,
    debounceValue: debounceMaxLength,
    setValue: setMaxLength,
  } = useDebounce(attr.maxLength?.toString() ?? "", [attr.maxLength]);

  const {
    value: placeholder,
    debounceValue: debouncePlaceholder,
    setValue: setPlaceholder,
  } = useDebounce(attr.placeholder, [attr.placeholder]);

  const {
    value: defaultValue,
    debounceValue: debounceDefaultValue,
    setValue: setDefaultValue,
  } = useDebounce(attr.defaultValue, [attr.defaultValue]);

  useEffect(() => {
    const change: Partial<TextElementAttributes> = {
      minLength: stringToNumber(debounceMinLength),
      maxLength: stringToNumber(debounceMaxLength),
      placeholder: debouncePlaceholder ?? undefined,
      defaultValue: debounceDefaultValue ?? undefined,
    };

    onChange(change);
  }, [
    debounceMinLength,
    debounceMaxLength,
    debouncePlaceholder,
    debounceDefaultValue,
  ]);

  const handleChange = (field: keyof TextElementAttributes) => {
    return (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      const value = event.target.value;

      switch (field) {
        case "minLength":
          setMinLength(value);
          break;
        case "maxLength":
          setMaxLength(value);
          break;
        case "placeholder":
          setPlaceholder(value);
          break;
        case "defaultValue":
          setDefaultValue(value);
          break;
      }
    };
  };

  return (
    <Box sx={settingsStyle.container}>
      <TextField
        type="number"
        label="Min. Length"
        value={minLength}
        onChange={handleChange("minLength")}
      />
      <TextField
        type="number"
        label="Max Length"
        value={maxLength}
        onChange={handleChange("maxLength")}
      />
      <TextField
        label="Placeholder"
        value={placeholder}
        onChange={handleChange("placeholder")}
      />
      <TextField
        label="Default Value"
        value={defaultValue}
        onChange={handleChange("defaultValue")}
      />
    </Box>
  );
};
