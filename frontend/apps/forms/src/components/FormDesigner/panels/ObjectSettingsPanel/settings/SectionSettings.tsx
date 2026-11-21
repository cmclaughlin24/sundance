import { useDebouncedCallback } from "@/hooks/useDebouncedCallback";
import type { ISection } from "@/types/section";
import { stringToNumber } from "@/utils/utils";
import Box from "@mui/material/Box";
import { useEffect, useState, type ChangeEvent } from "react";
import { settingsStyle } from "./Settings.style";
import TextField from "@mui/material/TextField";

export type SectionSettingsChangeEvent =
  | {
      minInstances: number;
      maxInstances?: number;
      addLabel?: string;
    }
  | undefined;

interface Fields {
  minInstances: string;
  maxInstances?: string;
  addLabel?: string;
}

export const SectionSettings: React.FC<{
  section: ISection;
  onChange: (event: SectionSettingsChangeEvent) => void;
}> = function ({ section, onChange }) {
  const [minInstances, setMinInstances] = useState(
    section.repeat?.minInstances.toString() ?? "",
  );
  const [maxInstances, setMaxInstances] = useState(
    section.repeat?.maxInstances?.toString() ?? "",
  );
  const [addLabel, setAddLabel] = useState(section.repeat?.addInstanceLabel ?? "");

  const {
    debounced: debounceEmit,
    cancel,
    flush,
  } = useDebouncedCallback((next: Fields) => {
    let event: SectionSettingsChangeEvent = undefined;
    const minInstances = stringToNumber(next.minInstances);
    const maxInstances = stringToNumber(next.maxInstances);

    if (minInstances && minInstances >= 0) {
      event = { minInstances };
    }

    if (event && maxInstances && maxInstances >= 0) {
      event.maxInstances = maxInstances;
    }

    if (event && next.addLabel && next.addLabel.trim() !== "") {
      event.addLabel = addLabel;
    }

    onChange(event);
  });

  useEffect(() => {
    cancel();
    setMinInstances(section.repeat?.minInstances.toString() ?? "");
    setMaxInstances(section.repeat?.maxInstances?.toString() ?? "");
    setAddLabel(section.repeat?.addInstanceLabel ?? "");
  }, [section]);

  const handleChange = (field: keyof Fields) => {
    return (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      const value = event.target.value;
      const next: Fields = { minInstances, maxInstances, addLabel };

      switch (field) {
        case "minInstances":
          next.minInstances = value;
          setMinInstances(value);
          break;
        case "maxInstances":
          next.maxInstances = value;
          setMaxInstances(value);
          break;
        case "addLabel":
          next.addLabel = value;
          setAddLabel(value);
          break;
      }

      debounceEmit(next);
    };
  };

  const handleBlur = () => flush({ minInstances, maxInstances, addLabel });

  return (
    <Box sx={settingsStyle.container}>
      <Box sx={{ display: "flex", gap: 1 }}>
        <TextField
          type="number"
          label="Min. Instances"
          value={minInstances}
          onChange={handleChange("minInstances")}
          onBlur={handleBlur}
          slotProps={{
            htmlInput: { min: 0 },
          }}
        />
        <TextField
          type="number"
          label="Max Instances"
          value={maxInstances}
          onChange={handleChange("maxInstances")}
          onBlur={handleBlur}
          slotProps={{
            htmlInput: { min: 0 },
          }}
        />
      </Box>
      <TextField
        label="Add Instance Label"
        value={addLabel}
        onChange={handleChange("addLabel")}
        onBlur={handleBlur}
      />
    </Box>
  );
};
