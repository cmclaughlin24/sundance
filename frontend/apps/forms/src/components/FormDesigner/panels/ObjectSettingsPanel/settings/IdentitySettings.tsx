import Box from "@mui/material/Box";
import { settingsStyle } from "./Settings.style";
import TextField from "@mui/material/TextField";
import { useEffect, useState, type ChangeEvent } from "react";
import type {
  UpdateElementEvent,
  UpdateSectionEvent,
} from "@/store/formDesigner/events";
import type { ISection } from "@/types/section";
import type { IPage } from "@/types/page";
import type { IElement } from "@/types/element";
import { useDebouncedCallback } from "@/hooks/useDebouncedCallback";

interface SectionIdentitySettingsProps {
  type: "section";
  object: ISection;
}

interface PageIdentitySettingsProps {
  type: "page";
  object: IPage;
}

interface ElementIdentitySettingsProps {
  type: "element";
  object: IElement;
}

export type IdentitySettingsProps =
  | SectionIdentitySettingsProps
  | PageIdentitySettingsProps
  | ElementIdentitySettingsProps;

interface SectionIdentitySettingsEvent {
  type: "section";
  changes: UpdateSectionEvent["changes"];
}

interface ElementIdentitySettingsEvent {
  type: "element";
  changes: UpdateElementEvent["changes"];
}

export type IdentitySettingsEvent =
  SectionIdentitySettingsEvent | ElementIdentitySettingsEvent;

type Fields = { name: string; description: string; key: string };

export const IdentitySettings: React.FC<
  IdentitySettingsProps & { onChange: (event: IdentitySettingsEvent) => void }
> = function ({ type, object, onChange }) {
  const [name, setName] = useState(object.name ?? "");
  const [description, setDescription] = useState(
    object && type === "element" ? object.description : "",
  );
  const [technicalKey, setTechnicalKey] = useState(object.key ?? "");

  const {
    debounced: debounceEmit,
    cancel,
    flush,
  } = useDebouncedCallback((next: Fields) => {
    let event: IdentitySettingsEvent;

    switch (type) {
      case "page":
        // FIXME: Implement the UpdatePage event.
        break;
      case "section":
        event = {
          type: "section",
          changes: { name: next.name, key: next.key },
        } satisfies SectionIdentitySettingsEvent;
        break;
      case "element":
        event = {
          type: "element",
          changes: {
            name: next.name,
            description: next.description,
            key: next.key,
          },
        } satisfies ElementIdentitySettingsEvent;
        break;
    }

    onChange(event!);
  });

  useEffect(() => {
    cancel();
    setName(object.name ?? "");
    setDescription(object && type === "element" ? object.description : "");
    setTechnicalKey(object.key ?? "");
  }, [object]);

  const handleChange = (field: "name" | "description" | "key") => {
    return (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      const value = event.target.value;
      const next: Fields = { name, description, key: technicalKey };

      switch (field) {
        case "name":
          next.name = value;
          setName(value);
          break;
        case "description":
          next.description = value;
          setDescription(value);
          break;
        case "key":
          next.key = value;
          setTechnicalKey(value);
          break;
      }

      debounceEmit(next);
    };
  };

  const handleBlur = () => flush({ name, description, key: technicalKey });

  return (
    <Box sx={settingsStyle.container}>
      <TextField
        label="Title"
        value={name}
        onChange={handleChange("name")}
        onBlur={handleBlur}
      />
      {type === "element" && (
        <TextField
          multiline
          label="Description"
          value={description}
          rows={3}
          onChange={handleChange("description")}
          onBlur={handleBlur}
        />
      )}
      <TextField
        label="Technical Key"
        value={technicalKey}
        onChange={handleChange("key")}
        onBlur={handleBlur}
      />
    </Box>
  );
};
