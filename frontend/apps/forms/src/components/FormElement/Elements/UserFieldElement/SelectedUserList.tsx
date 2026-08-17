import type { IUserLookup } from "@/types/userLookup";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { UserChip } from "./UserChip";

export interface SelectedUserListProps {
  selections: IUserLookup[];
  onRemove: (option: IUserLookup) => void;
}

export const SelectedUserList: React.FC<SelectedUserListProps> = function ({
  selections,
  onRemove,
}) {
  if (!selections || !selections.length) {
    return (
      <Box sx={{}}>
        <Typography sx={{}}>Select employee</Typography>
      </Box>
    );
  }

  return (
    <Box component="ul">
      {selections.map((option) => (
        <Box component="li" key={option.value}>
          <UserChip option={option} onRemove={onRemove} />
        </Box>
      ))}
    </Box>
  );
};
