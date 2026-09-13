import type { IDBPDatabase, IDBPTransaction } from "idb";

/**
 * Interface defining the options for configuring the `IDBStorageService`.
 */
export interface IDBServiceOptions {
  /**
   * The name of the IndexedDB database.
   */
  databaseName: string;

  /**
   */
  migrations: IDBMigrationStep[];
}

export interface IDBMigrationStep {
  version: number;

  description: string;

  migrate: (
    db: IDBPDatabase,
    transaction: IDBPTransaction<any, any, "versionchange">,
  ) => void | Promise<void>;
}
