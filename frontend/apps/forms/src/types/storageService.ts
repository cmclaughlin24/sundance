/**
 * A generic interface for a storage service that provides basic CRUD operations.
 * @template K The type of the key used to identify records in the storage service.
 * @template V The type of entity to be managed by the storage service.
 */
export interface StorageService<K, V> {
  /**
   * Finds an entity by its unique key.
   * @param key The unique identifier of the entity to find.
   * @returns The entity if found, or a promise that resolves to the entity.
   */
  findByKey(key: K): V | Promise<V>;

  /**
   * Adds or updates an entity in the storage.
   * @param value The entity to add.
   * @returns Either void or a promise that resolves when the entity has been added.
   */
  upsert(value: V): void | Promise<void>;

  /**
   * Deletes an entity from the storage by its unique identifier.
   * @param id The unique identifier of the entity to delete.
   * @returns Either void or a promise that resolves when the entity has been deleted.
   */
  delete(key: K): void | Promise<void>;
}
