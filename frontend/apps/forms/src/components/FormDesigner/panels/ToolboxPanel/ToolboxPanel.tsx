import { Panel } from "@/components/layout/Panel";
import Typography from "@mui/material/Typography";
import { filterPalette, type IPaletteCategory } from "./palette";
import { PaletteCategory } from "./PaletteCategory";
import { useDebounce } from "@/hooks/useDebounce";
import type { ChangeEvent } from "react";
import TextField from "@mui/material/TextField";

export interface ToolboxPanelProps<T> {
  palette: IPaletteCategory<T>[];
  helpText: string;
}

export function ToolboxPanel<T>({ palette, helpText }: ToolboxPanelProps<T>) {
  const {
    value: searchTerm,
    debounceValue: debounceSearchTerm,
    setValue: setSearchTerm,
  } = useDebounce("", []);

  const handleChange = (
    event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    setSearchTerm(event.target.value);
  };

  const pallette = filterPalette(debounceSearchTerm, palette);

  return (
    <Panel sx={{ display: "flex", flexDirection: "column", gap: 2.5 }}>
      <Panel.Title title="Library" />
      <TextField
        placeholder="Search"
        value={searchTerm}
        onChange={handleChange}
      />
      <Typography> {helpText} </Typography>
      {pallette.map((category) => (
        <PaletteCategory category={category} key={category.label} />
      ))}
    </Panel>
  );
}
