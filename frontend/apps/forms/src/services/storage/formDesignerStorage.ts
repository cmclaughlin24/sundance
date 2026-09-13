import type { FormDesignerDraft } from "@/store/formDesigner";
import { IDBRepository } from "../idb";
import { FORMS_DB_STORES, formsIDBService } from "../idb/formsIDBService";


export const formsDraftStorage = new IDBRepository<string, FormDesignerDraft>(
  formsIDBService,
  FORMS_DB_STORES.designerDraft,
);
