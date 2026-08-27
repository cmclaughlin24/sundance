import { Panel } from "@/components/layout/Panel";
import Typography from "@mui/material/Typography";
import { filterPalette } from "./palette";
import { PaletteCategory } from "./PaletteCategory";
import { useDebounce } from "@/hooks/useDebounce";
import type { ChangeEvent } from "react";
import TextField from "@mui/material/TextField";

export const ToolboxPanel: React.FC = function () {
  const {
    value: searchTerm,
    debounceValue: debounceSearchTerm,
    setValue: setSearchTerm,
  } = useDebounce("");

  const handleChange = (
    event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    setSearchTerm(event.target.value);
  };

  const pallette = filterPalette(debounceSearchTerm);

  return (
    <Panel sx={{ display: "flex", flexDirection: "column", gap: 2.5 }}>
      <Panel.Title title="Library" />
      <TextField
        placeholder="Search"
        value={searchTerm}
        onChange={handleChange}
      />
      <Typography>
        Drag the form elements into the preferred section on the canvas.
      </Typography>
      {pallette.map((category) => (
        <PaletteCategory category={category} key={category.label} />
      ))}
    </Panel>
  );
};
