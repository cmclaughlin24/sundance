import { checkElementType } from "@/utils/error";
import { FieldElementContainer } from "../../Layout/FieldElementContainer";
import type { ElementComponent } from "../../Renderer/ElementRenderer";
import { useElementValue } from "@/store/submission/useSubmissionContext";
import type { IUserLookup } from "@/types/userLookup";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { userFieldElementStyles } from "./UserFieldElement.style";
import FormControlLabel from "@mui/material/FormControlLabel";
import Checkbox from "@mui/material/Checkbox";
import { SelectedUserList } from "./SelectedUserList";
import { useCallback, type ChangeEvent } from "react";
import { SearchBar } from "@/components/SearchBar/Searchbar";
import { useUsersService } from "@/hooks/useHttpService";
import { useTenantId } from "@/store/formDefinition";

const minimumSearchCharacters = 6;

export const UserFieldElement: ElementComponent = function ({
  element,
  onChange,
}) {
  checkElementType(element.type, "user");

  const usersService = useUsersService();
  const tenantId = useTenantId();
  const value = useElementValue<IUserLookup[]>(element.id, []);

  const findUsers = useCallback(
    async (token: string, searchTerm?: string) => {
      if (
        !searchTerm ||
        searchTerm.length < minimumSearchCharacters ||
        !token ||
        !tenantId
      ) {
        return [];
      }

      return await usersService.getUserLookups(searchTerm, { tenantId, token });
    },
    [tenantId, usersService],
  );

  const handleOnMyself = (event: ChangeEvent<HTMLInputElement>) => {
    // TODO: Implement logic to lookup the current user & add them to the list.
    console.log(event.target.checked);
  };

  const handleSelection = (option: IUserLookup) => {
    onChange([...value, option]);
  };

  const handleRemoveAll = () => {
    onChange([]);
  };

  const handleRemove = (option: IUserLookup) => {
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
          control={
            <Checkbox
              onChange={handleOnMyself}
              data-testid="user-field-element-for-myself"
            />
          }
        />
        <Box sx={userFieldElementStyles["search"]}>
          <Typography variant="body1">
            Search employee(s) by name, email, or ID:
          </Typography>
          <SearchBar
            value={null}
            queryFn={findUsers}
            onSelection={handleSelection}
            helperText={`Must enter ${minimumSearchCharacters} characters to search by name or ID`}
          />
        </Box>
        {list}
      </Box>
    </FieldElementContainer>
  );
};
