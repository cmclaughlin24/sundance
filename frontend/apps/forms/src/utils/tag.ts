import type { ITag } from "@/types/tag";

export interface TagGroup {
  id: string;
  title: string;
  keyPath: string;
  items: ITag[];
}

export function groupTags(tags: ITag[], depth: number = 2): TagGroup[] {
  const groups = new Map<string, TagGroup>();

  groups.set("", {
    id: "root",
    title: "General Attributes",
    keyPath: "",
    items: [],
  });

  for (const tag of tags) {
    if (tag.nodeType !== "object" || tag.keyPath.split(".").length > depth) {
      continue;
    }

    groups.set(tag.keyPath, {
      id: tag.id,
      title: tag.displayName,
      keyPath: tag.keyPath,
      items: [],
    });
  }

  for (const tag of tags) {
    if (tag.nodeType !== "primitive") {
      continue;
    }

    const parts = tag.keyPath.split(".");

    const groupKey =
      parts.length > 1
        ? parts.slice(0, Math.min(depth, parts.length - 1)).join(".")
        : "";

    if (!groups.has(groupKey)) {
      groups.set(groupKey, {
        id: groupKey,
        title: `${groupKey}.*`,
        keyPath: groupKey,
        items: [],
      });
    }

    groups.get(groupKey)!.items.push(tag);
  }

  return Array.from(groups.values()).filter((g) => g.items.length > 0);
}
