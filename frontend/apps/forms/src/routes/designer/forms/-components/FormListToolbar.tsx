import SearchIcon from "@mui/icons-material/Search";
import type { Styles } from "@/types/styles";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import InputAdornment from "@mui/material/InputAdornment";
import TextField from "@mui/material/TextField";

const styles: Styles = {
  toolbar: {
    mb: 3,
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
  },
  searchInput: {
    width: "20rem",
    maxWidth: "100%",
  },
};

export interface FormListToolbarProps {
  search: string;
  onSearch: (value: string) => void;
  onNew: () => void;
}

export const FormListToolbar: React.FC<FormListToolbarProps> = function ({
  search,
  onSearch,
  onNew,
}) {
  return (
    <Box sx={styles.toolbar}>
      <Box>
        <Button variant="outlined" onClick={onNew} data-testid="new-form-btn">
          New Form
        </Button>
      </Box>
      <Box>
        <TextField
          variant="outlined"
          placeholder="Search Forms"
          aria-label="search-forms"
          value={search}
          onChange={(event) => onSearch(event.target.value)}
          data-testid="form-list-toolbar-search-input"
          slotProps={{
            input: {
              endAdornment: (
                <InputAdornment position="end">
                  <SearchIcon sx={{ color: "#505050" }} />
                </InputAdornment>
              ),
            },
          }}
        />
      </Box>
    </Box>
  );
};
