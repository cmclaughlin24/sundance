import { openDB, type IDBPDatabase } from "idb";
import type { IDBServiceOptions } from "./idbService.type";

export class IDBService {
  private _database: Promise<IDBPDatabase> | null = null;

  constructor(private readonly options: IDBServiceOptions) {}

  async database(): Promise<IDBPDatabase> {
    if (!this._database) {
      this._database = this._open().catch((err) => {
        this._database = null;
        throw err;
      });
    }

    return this._database;
  }

  private async _open(): Promise<IDBPDatabase> {
    const migrations = this.options.migrations;
    const targetVersion = Math.max(...migrations.map((m) => m.version), 1);

    return openDB(this.options.databaseName, targetVersion, {
      upgrade: async (database, oldVersion, newVersion, transaction) => {
        const pendingMigrations = migrations
          .sort((a, b) => a.version - b.version)
          .filter(
            (m) =>
              m.version > oldVersion &&
              (newVersion === null || m.version <= newVersion),
          );

        for (const step of pendingMigrations) {
          await step.migrate(database, transaction);
        }
      },
    });
  }
}
