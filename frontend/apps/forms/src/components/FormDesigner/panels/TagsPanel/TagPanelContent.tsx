import type { Styles } from "@/types/styles";
import Box from "@mui/material/Box";
import TextField from "@mui/material/TextField";
import * as ArrayUtils from "@/utils/array";
import Typography from "@mui/material/Typography";
import { useMemo, type ChangeEventHandler } from "react";
import { useDebounce } from "@/hooks/useDebounce";

const styles: Styles = {
  list: {
    p: 0,
    m: 0,
  },
  listItem: {
    listStyle: "none",
  },
  noResults: { textAlign: "center" },
};

export interface TagsPanelContentProps<IType> {
  data: IType[] | undefined | null;
  children: (item: IType) => React.ReactNode | undefined;
  filterFn: (search: string, item: IType) => boolean;
  placeholder?: string;
}

export function TagsPanelContent<IType>({
  data = [],
  placeholder = "",
  children,
  filterFn,
}: TagsPanelContentProps<IType>) {
  const {
    value: searchTerm,
    debounceValue: debounceSearchTerm,
    setValue: setSearchTerm,
  } = useDebounce("", []);

  const filtered = useMemo(() => {
    if (!ArrayUtils.hasLengthGreaterThan(data, 0)) {
      return [];
    }

    if (debounceSearchTerm.trim() === "") {
      return data;
    }

    return data!.filter((i) => filterFn(debounceSearchTerm, i));
  }, [data, debounceSearchTerm, filterFn]);

  const handleChange: ChangeEventHandler<
    HTMLInputElement | HTMLTextAreaElement
  > = (event) => {
    setSearchTerm(event.target.value);
  };

  const hasItems = ArrayUtils.hasLengthGreaterThan(filtered, 0);

  return (
    <>
      <TextField
        placeholder={placeholder}
        value={searchTerm}
        onChange={handleChange}
      />
      {hasItems && (
        <Box component="ul" sx={styles.list}>
          {filtered!.map((item) => (
            <Box component="li" sx={styles.listItem}>
              {children(item)}
            </Box>
          ))}
        </Box>
      )}
      {!hasItems && <Typography sx={styles.noResults}>No Results</Typography>}
    </>
  );
}
