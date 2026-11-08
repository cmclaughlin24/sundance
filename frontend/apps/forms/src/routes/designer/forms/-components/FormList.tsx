import type { IForm } from "@/types/form";
import { FormCard } from "./FormCard";
import { CardList } from "@/components/CardList";

export interface FormListProps {
  forms: IForm[] | null;
  onClick: (form: IForm) => void;
}

export const FormList: React.FC<FormListProps> = function ({ forms, onClick }) {
  return (
    <CardList data={forms} onClick={onClick} keyFn={(form) => form.id}>
      {(form) => <FormCard form={form} />}
    </CardList>
  );
};
