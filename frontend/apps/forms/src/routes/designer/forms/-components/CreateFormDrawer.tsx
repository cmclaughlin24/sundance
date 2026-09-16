import { Drawer, type DrawerHandle } from "@/components/Drawer";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { useActionState, useImperativeHandle, useRef } from "react";
import { createFormDrawerStyles } from "./CreateFormDrawer.style";

export interface CreateFormDrawerProps {
  ref: React.Ref<DrawerHandle>;
  onCreate: (event: CreateFormEvent) => void;
}

export interface CreateFormEvent {
  name: string;
  description: string;
}

interface CreateFormState {
  error?: string;
}

export const CreateFormDrawer: React.FC<CreateFormDrawerProps> = function ({
  ref,
  onCreate,
}) {
  const drawerRef = useRef<DrawerHandle>(null);

  useImperativeHandle(ref, () => ({
    open: () => drawerRef.current?.open(),
    close: () => drawerRef.current?.close(),
  }));

  const handleCancel = () => drawerRef.current?.close();

  const [_state, formAction, isPending] = useActionState<
    CreateFormState,
    FormData
  >(async (_prev, formData) => {
    const name = String(formData.get("name") ?? "").trim();
    const description = String(formData.get("description") ?? "").trim();

    if (!name) {
      return { error: "Name is required" };
    }

    if (!description) {
      return { error: "Description is required" };
    }

    onCreate({ name, description });
    drawerRef.current?.close();
    return {};
  }, {});

  return (
    <Drawer title="New Form" ref={drawerRef}>
      <Box component="hgroup" sx={createFormDrawerStyles.title}>
        <Typography sx={createFormDrawerStyles.h3}>
          Create a new form.
        </Typography>
        <Typography variant="body2">
          Start by giving your form a name and short description. You'll build
          out it's fields next.
        </Typography>
      </Box>
      <Box
        component="form"
        action={formAction}
        sx={createFormDrawerStyles.form}
        data-testid="create-form-drawer-form"
      >
        <Box sx={createFormDrawerStyles.fields}>
          <TextField
            name="name"
            label="Name"
            required
            data-testid="create-form-name-input"
          />
          <TextField
            name="description"
            label="Description"
            required
            multiline
            rows={3}
            data-testid="create-form-description-input"
          />
        </Box>
        <Box sx={createFormDrawerStyles.buttons}>
          <Button
            onClick={handleCancel}
            variant="text"
            data-testid="cancel-create-form-btn"
          >
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={isPending}
            data-testid="cancel-create-form-btn"
          >
            Create
          </Button>
        </Box>
      </Box>
    </Drawer>
  );
};
