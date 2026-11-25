import { useDebouncedCallback } from "@/hooks/useDebouncedCallback";
import type { ISection, ISectionRepeat } from "@/types/section";
import { stringToNumber } from "@/utils/utils";
import Box from "@mui/material/Box";
import { useEffect, useState, type ChangeEvent } from "react";
import { settingsStyle } from "./Settings.style";
import TextField from "@mui/material/TextField";
import * as Utils from "@/utils/utils";

export type SectionSettingsChangeEvent = ISectionRepeat | undefined;

interface Fields {
  minInstances: string;
  maxInstances?: string;
  addLabel?: string;
  instanceLabel?: string;
}

export const SectionSettings: React.FC<{
  section: ISection;
  onChange: (event: SectionSettingsChangeEvent) => void;
}> = function ({ section, onChange }) {
  const [minInstances, setMinInstances] = useState(
    Utils.numberToString(section.repeat?.minInstances),
  );
  const [maxInstances, setMaxInstances] = useState(
    Utils.numberToString(section.repeat?.maxInstances),
  );
  const [addLabel, setAddLabel] = useState(section.repeat?.addLabel ?? "");
  const [instanceLabel, setInstanceLabel] = useState(
    section.repeat?.instanceLabel ?? "",
  );

  const {
    debounced: debounceEmit,
    cancel,
    flush,
  } = useDebouncedCallback((next: Fields) => {
    let event: SectionSettingsChangeEvent = {};
    const minInstances = stringToNumber(next.minInstances);
    const maxInstances = stringToNumber(next.maxInstances);

    if (minInstances != null && minInstances >= 0) {
      event = { minInstances };
    }

    if (maxInstances != null && maxInstances > 0) {
      event.maxInstances = maxInstances;
    }

    if (next.addLabel && next.addLabel.trim() !== "") {
      event.addLabel = next.addLabel;
    }

    if (next.instanceLabel && next.instanceLabel.trim() !== "") {
      event.instanceLabel = next.instanceLabel;
    }

    onChange(Object.keys(event).length > 0 ? event : undefined);
  });

  useEffect(() => {
    cancel();
    setMinInstances(Utils.numberToString(section.repeat?.minInstances));
    setMaxInstances(Utils.numberToString(section.repeat?.maxInstances));
    setAddLabel(section.repeat?.addLabel ?? "");
    setInstanceLabel(section.repeat?.instanceLabel ?? "");
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
        case "instanceLabel":
          next.instanceLabel = value;
          setInstanceLabel(value);
          break;
      }

      debounceEmit(next);
    };
  };

  const handleBlur = () =>
    flush({ minInstances, maxInstances, addLabel, instanceLabel });

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
      <TextField
        label="Instance Label"
        value={instanceLabel}
        onChange={handleChange("instanceLabel")}
        onBlur={handleBlur}
      />
    </Box>
  );
};
