import type { StorageService } from "@/types/storageService";
import type { IDBService } from "./idbService";

export class IDBRepository<K extends IDBValidKey, V> implements StorageService<
  K,
  V
> {
  constructor(
    private readonly idService: IDBService,
    private readonly storeName: string,
  ) {}

  async findByKey(key: K): Promise<V> {
    const db = await this.idService.database();
    return db.get(this.storeName, key);
  }

  async upsert(value: V): Promise<void> {
    const db = await this.idService.database();
    await db.put(this.storeName, value);
  }

  async delete(key: K): Promise<void> {
    const db = await this.idService.database();
    await db.delete(this.storeName, key);
  }
}
