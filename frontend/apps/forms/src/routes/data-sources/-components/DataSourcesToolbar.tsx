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

export interface DataSourcesToolbarFilters {
  search: string;
}

export interface DataSourcesToolbarProps {
  filters: DataSourcesToolbarFilters;
  onFilter: (fitler: DataSourcesToolbarFilters) => void;
  onNew: () => void;
}

export const DataSourcesToolbar: React.FC<DataSourcesToolbarProps> = function ({
  filters = { search: "" },
  onFilter,
  onNew,
}) {
  return (
    <Box sx={styles.toolbar}>
      <Box>
        <Button variant="outlined" onClick={onNew} data-testid="new-form-btn">
          New Data Source
        </Button>
      </Box>
      <Box>
        <TextField
          variant="outlined"
          placeholder="Search Data Sources"
          aria-label="search-data-sources"
          value={filters?.search}
          onChange={(event) =>
            onFilter({ ...filters, search: event.target.value })
          }
          data-testid="data-sources-toolbar-search-input"
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
