import { IDBService } from "./idbService";
import type { IDBMigrationStep } from "./idbService.type";

export const FORMS_DB_STORES = {
  designerDraft: "form_designer_drafts",
} as const;

export const migrationV1: IDBMigrationStep = {
  version: 1,
  description: "Create form designer drafts store",
  migrate: (db) => {
    if (!db.objectStoreNames.contains(FORMS_DB_STORES.designerDraft)) {
      db.createObjectStore(FORMS_DB_STORES.designerDraft, {
        keyPath: "versionId",
      });
    }
  },
};

export const formsIDBService = new IDBService({
  databaseName: "sundance_forms",
  migrations: [migrationV1],
});
