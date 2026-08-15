import { checkElementType } from "@/utils/error";
import { FieldElementContainer } from "../../Layout/FieldElementContainer";
import type { ElementComponent } from "../../Renderer/ElementRenderer";
import { useElementValue } from "@/store/submission/useSubmissionContext";
import type { IUserOption } from "@/types/userOption";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { userFieldElementStyles } from "./UserFieldElement.style";
import FormControlLabel from "@mui/material/FormControlLabel";
import Checkbox from "@mui/material/Checkbox";
import { SelectedUserList } from "./SelectedUserList";

export const UserFieldElement: ElementComponent = function ({
  element,
  onChange,
}) {
  checkElementType(element.type, "user");

  const value = useElementValue<IUserOption[]>(element.id, []);

  const handleRemoveAll = () => {
    onChange([]);
  };

  const handleRemove = (option: IUserOption) => {
    if (!value || !value.length) {
      console.warn(
        "UserFieldElement attempted to remove user but element does not have a value",
      );
      return;
    }

    onChange(value.filter((v) => v.value !== option.value));
  };

  let list: React.ReactNode;

  if (value && value.length) {
    list = (
      <Box>
        <Typography variant="body2" sx={userFieldElementStyles["selectedList"]}>
          {value.length} employee(s) selected
        </Typography>
        <Box
          component="button"
          onClick={handleRemoveAll}
          disabled={!value || !value.length}
          sx={userFieldElementStyles["removeAllBtn"]}
          data-testid="user-field-element-remove-all"
        >
          Remove All
        </Box>
        <SelectedUserList selections={value} onRemove={handleRemove} />
      </Box>
    );
  }

  return (
    <FieldElementContainer element={element}>
      <Box sx={userFieldElementStyles["inputContainer"]}>
        <FormControlLabel
          label="For myself"
          control={<Checkbox data-testid="user-field-element-for-myself" />}
        />
        <Box sx={userFieldElementStyles["search"]}>
          <Typography variant="body1">
            Search employee(s) by name, email, or ID:
          </Typography>
          {/* TODO: Add the missing SearchBar component */}
        </Box>
        {list}
      </Box>
    </FieldElementContainer>
  );
};
