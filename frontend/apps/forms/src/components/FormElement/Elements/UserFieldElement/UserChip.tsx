import type { IUserOption } from "@/types/userOption";
import Box from "@mui/material/Box";
import { userChipStyles } from "./UserChip.style";
import Typography from "@mui/material/Typography";

export interface UserChipProps {
  option: IUserOption;
  onRemove: (value: IUserOption) => void;
}

export const UserChip: React.FC<UserChipProps> = function ({
  option,
  onRemove,
}) {
  return (
    <Box sx={userChipStyles["chip"]} data-testid={`user-chip-${option.value}`}>
      <Typography sx={userChipStyles["label"]}>{option.label}</Typography>
      <Box
        component="button"
        onClick={() => onRemove(option)}
        sx={userChipStyles["removeBtn"]}
        data-testid={`remove-user-chip-${option.value}`}
      >
        X
      </Box>
    </Box>
  );
};
