import { Drawer, type DrawerHandle } from "@/components/Drawer";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { useActionState, useImperativeHandle, useRef, useState } from "react";
import { createDataSourceDrawerStyles as styles } from "./CreateDataSourceDrawer.style";
import type { DataSourceType } from "@/types/dataSource";
import { DataSourceCardList } from "./DataSourceCardList";

export interface CreateDataSourceDrawerProps {
  ref: React.Ref<DrawerHandle>;
  onCreate: (event: CreateDataSourceEvent) => void;
}

export interface CreateDataSourceEvent {
  name: string;
  description: string;
  type: DataSourceType;
}

interface CreateFormState {
  error?: string;
}


export const CreateDataSourceDrawer: React.FC<CreateDataSourceDrawerProps> =
  function ({ ref, onCreate }) {
    const drawerRef = useRef<DrawerHandle>(null);
    const [selectedType, setSelectedType] =
      useState<DataSourceType>("scheduled");

    useImperativeHandle(ref, () => ({
      open: () => {
        setSelectedType("scheduled");
        drawerRef.current?.open();
      },
      close: () => drawerRef.current?.close(),
    }));

    const handleCancel = () => drawerRef.current?.close();

    const [_state, formAction, isPending] = useActionState<
      CreateFormState,
      FormData
    >(async (_prev, formData) => {
      const name = String(formData.get("name") ?? "").trim();
      const description = String(formData.get("description") ?? "").trim();
      const type = formData.get("type") as DataSourceType;

      if (!name) {
        return { error: "Name is required" };
      }

      if (!description) {
        return { error: "Description is required" };
      }

      onCreate({ name, description, type });
      drawerRef.current?.close();
      return {};
    }, {});

    return (
      <Drawer title="New Data Source" ref={drawerRef}>
        <Box component="hgroup" sx={styles.title}>
          <Typography sx={styles.h3}>Create a Brand New Data Source</Typography>
          <Typography variant="body2">
            Provide a name and description to register this data source. You
            will select its type (static, scheduled, webhook, or data lake) and
            configure connection attributes next.
          </Typography>
        </Box>
        <Box
          component="form"
          action={formAction}
          sx={styles.form}
          data-testid="create-data-source-drawer-form"
        >
          <input type="hidden" name="type" value={selectedType} />
          <Box sx={styles.fields}>
            <Typography>1. Tell us about this data source</Typography>
            <TextField
              name="name"
              label="Name"
              required
              data-testid="create-data-source-name-input"
            />
            <TextField
              name="description"
              label="Description"
              required
              multiline
              rows={3}
              data-testid="create-data-source-description-input"
            />
            <Typography>2. Where does the data come from?</Typography>
            <DataSourceCardList selection={selectedType} onSelection={setSelectedType} />
          </Box>
          <Box sx={styles.buttons}>
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
              data-testid="submit-create-data-source-btn"
            >
              Create
            </Button>
          </Box>
        </Box>
      </Drawer>
    );
  };
