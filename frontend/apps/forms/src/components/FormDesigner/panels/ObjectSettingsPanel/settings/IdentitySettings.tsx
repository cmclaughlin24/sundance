import Box from "@mui/material/Box";
import { settingsStyle } from "./Settings.style";
import TextField from "@mui/material/TextField";
import { useFormDesignerDispatch } from "@/store/formDesigner";
import { useDebounce } from "@/hooks/useDebounce";
import { useEffect, type ChangeEvent } from "react";
import type {
  FormDesignerEvent,
  UpdateElementEvent,
  UpdateSectionEvent,
} from "@/store/formDesigner/events";
import type { ISection } from "@/types/section";
import type { IPage } from "@/types/page";
import type { IElement } from "@/types/element";

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

export const IdentitySettings: React.FC<IdentitySettingsProps> = function ({
  type,
  object,
}) {
  const { dispatch } = useFormDesignerDispatch();

  const {
    value: name,
    debounceValue: debounceName,
    setValue: setName,
  } = useDebounce(object.name ?? "", [object]);

  const {
    value: description,
    debounceValue: debounceDescription,
    setValue: setDescription,
  } = useDebounce(object && type === "element" ? object.description : "", [
    object,
  ]);

  const {
    value: technicalKey,
    debounceValue: debounceTechnicalKey,
    setValue: setTechnicalKey,
  } = useDebounce(object.key ?? "", [object]);

  useEffect(() => {
    let event: FormDesignerEvent;

    switch (type) {
      case "page":
        // FIXME: Implement the UpdatePage event.
        break;
      case "section":
        event = {
          type: "UpdateSection",
          sectionId: object.id,
          changes: { name: debounceName, key: debounceTechnicalKey },
        } satisfies UpdateSectionEvent;
        break;
      case "element":
        event = {
          type: "UpdateElement",
          elementId: object.id,
          changes: {
            name: debounceName,
            description: debounceDescription,
            key: debounceTechnicalKey,
          },
        } satisfies UpdateElementEvent;
        break;
    }

    dispatch(event!);
  }, [debounceName, debounceDescription, debounceTechnicalKey]);

  const handleChange = (field: "name" | "description" | "key") => {
    return (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      const value = event.target.value;

      switch (field) {
        case "name":
          setName(value);
          break;
        case "description":
          setDescription(value);
          break;
        case "key":
          setTechnicalKey(value);
          break;
      }
    };
  };

  return (
    <Box sx={settingsStyle.container}>
      <TextField label="Title" value={name} onChange={handleChange("name")} />
      {type === "element" && (
        <TextField
          multiline
          label="Description"
          value={description}
          rows={3}
          onChange={handleChange("description")}
        />
      )}
      <TextField
        label="Technical Key"
        value={technicalKey}
        onChange={handleChange("key")}
      />
    </Box>
  );
};
