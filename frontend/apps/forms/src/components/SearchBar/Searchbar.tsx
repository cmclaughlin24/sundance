import { useAsyncData } from "@/hooks/useAsyncData";
import type { ILookup } from "@/types/data";
import type {
  AutocompleteOwnerState,
  AutocompleteRenderInputParams,
  AutocompleteRenderOptionState,
} from "@mui/material/Autocomplete";
import Box from "@mui/material/Box";
import FormControl from "@mui/material/FormControl";
import TextField from "@mui/material/TextField";
import type { AutocompleteValue } from "@mui/material/useAutocomplete";
import { type FocusEventHandler, type SyntheticEvent } from "react";
import { DefaultSearchItem } from "./DefaultSearchItem";
import Autocomplete from "@mui/material/Autocomplete";
import CircularProgress from "@mui/material/CircularProgress";
import Typography from "@mui/material/Typography";
import FormHelperText from "@mui/material/FormHelperText";
import { useDebounce } from "@/hooks/useDebounce";

export interface SearchBarProps<T> {
  value?: T | null;
  queryFn: (accessToken: string, searchTerm?: string) => Promise<T[]>;
  onSelection: (option: T) => void;
  children?: (option: T) => React.ReactNode;
  helperText?: string;
  onBlur?: FocusEventHandler<HTMLDivElement> | undefined;
  error?: boolean;
}

type RenderFn<T> = (
  props: React.HTMLAttributes<HTMLLIElement> & { key: any },
  option: T,
  state: AutocompleteRenderOptionState,
  ownerState: AutocompleteOwnerState<T, false, false, false, "div">,
) => React.ReactNode;

export function SearchBar<T extends ILookup>({
  queryFn,
  onSelection,
  children,
  helperText,
  error,
  ...props
}: SearchBarProps<T>) {
  const {
    value: searchTerm,
    debounceValue: debounceSearchTerm,
    setValue: setSearchTerm,
  } = useDebounce("", []);

  const { data, isLoading } = useAsyncData(
    (accessToken) => queryFn(accessToken, debounceSearchTerm),
    [queryFn, debounceSearchTerm],
  );

  const handleChange = (
    _event: SyntheticEvent,
    value: AutocompleteValue<T, false, false, false>,
  ) => {
    if (!value) {
      return;
    }

    setSearchTerm("");
    onSelection(value);
  };

  const handleInputChange = (_event: SyntheticEvent, value: string) => {
    setSearchTerm(value);
  };

  const handleRenderInput = (props: AutocompleteRenderInputParams) => {
    return <TextField {...props} error={error} />;
  };

  const handleRenderOptions: RenderFn<T> = (props, option) => {
    return (
      <Box component="li" {...props} key={option.value}>
        {children ? children(option) : <DefaultSearchItem option={option} />}
      </Box>
    );
  };

  return (
    <FormControl error={error}>
      <Autocomplete
        options={(data as Readonly<T[]>) ?? []}
        renderInput={handleRenderInput}
        renderOption={handleRenderOptions}
        onChange={handleChange}
        inputValue={searchTerm}
        onInputChange={handleInputChange}
        loading={isLoading}
        loadingText={
          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
            <CircularProgress size={16} />
            <Typography variant="body2">Searching...</Typography>
          </Box>
        }
        {...props}
      />
      {helperText && (
        <FormHelperText data-testid="search-bar-helper-text">
          {helperText}
        </FormHelperText>
      )}
    </FormControl>
  );
}
