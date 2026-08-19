import { checkElementType } from "@/utils/error";
import { FieldElementContainer } from "../../Layout/FieldElementContainer";
import type { ElementComponent } from "../../Renderer/ElementRenderer";
import {
  useElementErrors,
  useElementValue,
  useSubmissionDispatch,
} from "@/store/submission/useSubmissionContext";
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
import z from "zod";
import { ErrorMessages } from "@/constants/errorMessages";

const minimumSearchCharacters = 6;

export const UserFieldElement: ElementComponent = function ({
  element,
  ruleState,
  onChange,
}) {
  checkElementType(element.type, "user");

  const usersService = useUsersService();
  const tenantId = useTenantId();
  const { setError } = useSubmissionDispatch();
  const value = useElementValue<IUserLookup[]>(element.id, []);
  const errors = useElementErrors(element.id);
  const validationSchema = buildUserValidationSchema({
    required: ruleState.required,
  });

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
    handleBlur([]);
  };

  const handleRemove = (option: IUserLookup) => {
    if (!value || !value.length) {
      console.warn(
        "UserFieldElement attempted to remove user but element does not have a value",
      );
      return;
    }

    const updated = value.filter((v) => v.value !== option.value);
    onChange(updated);
    handleBlur(updated);
  };

  const handleBlur = (value: IUserLookup[]) => {
    const result = validationSchema.safeParse(value);
    setError(
      element.id,
      result.success ? [] : result.error.issues.map((e) => e.message),
    );
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

  const errorMesssage: string | null =
    errors && errors.length > 0 ? errors[0] : null;

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
            onBlur={() => handleBlur(value)}
            error={errorMesssage !== null}
            helperText={
              errorMesssage ??
              `Must enter ${minimumSearchCharacters} characters to search by name or ID`
            }
          />
        </Box>
        {list}
      </Box>
    </FieldElementContainer>
  );
};

function buildUserValidationSchema(options: {
  required: boolean;
}): z.ZodTypeAny {
  const item = z.object({
    key: z.string(),
    label: z.string(),
    value: z.union([z.string(), z.number()]),
  });

  const schema = z.array(item);

  if (!options.required) {
    return schema.optional();
  }

  return schema.min(1, ErrorMessages.required);
}
