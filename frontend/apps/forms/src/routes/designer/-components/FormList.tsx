import type { IForm } from "@/types/form";
import Box from "@mui/material/Box";
import { formListStyles } from "./FormList.style";
import { FormCard } from "./FormCard";
import { motion } from "motion/react";

export interface FormListProps {
  forms: IForm[] | null;
  onClick: (form: IForm) => void;
}

export const FormList: React.FC<FormListProps> = function ({ forms, onClick }) {
  return (
    <Box component="ul" sx={formListStyles["list"]}>
      {forms?.map((form) => (
        <Box
          component={motion.li}
          sx={formListStyles["item"]}
          onClick={() => onClick(form)}
          key={form.id}
        >
          <FormCard form={form} />
        </Box>
      ))}
    </Box>
  );
};
