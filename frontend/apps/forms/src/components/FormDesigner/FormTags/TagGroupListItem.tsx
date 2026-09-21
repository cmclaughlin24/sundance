import type { ITag, TagVersionStatus } from "@/types/tag";
import { TagsPanel } from "../panels/TagsPanel";
import { useFormDesignerSelect } from "@/store/formDesigner";
import * as ArrayUtils from "@/utils/array";
import { getLatestTagVersionByStatus } from "@/utils/tag";
import { TagVersionStatusTag } from "@/components/TagVersionStatusTag";
import { useContextMenuDispatch } from "@/components/ContextMenu";
import type { MouseEventHandler } from "react";
import type { TagVersionContextData } from "../menus/TagContextMenu";

export const TagGroupListItem: React.FC<{ tag: ITag }> = function ({ tag }) {
  const { selected } = useFormDesignerSelect();
  const { open } = useContextMenuDispatch();

  const handleContextMenu: MouseEventHandler<HTMLDivElement> = (event) => {
    event.preventDefault();
    event.stopPropagation();

    if (!activeVersion) {
      return;
    }

    open({
      position: { x: event.clientX, y: event.clientY },
      data: {
        type: "version",
        version: activeVersion,
      } satisfies TagVersionContextData,
    });
  };

  const isTagSelected = (tag: ITag): boolean => {
    if (
      !selected ||
      selected.type !== "element" ||
      !ArrayUtils.hasLengthGreaterThan(selected.item.tags, 0) ||
      !ArrayUtils.hasLengthGreaterThan(tag.versions ?? [], 0)
    ) {
      return false;
    }

    return tag.versions!.some((tv) =>
      selected.item.tags.some((t) => (t.tagVersionId === tv.id)),
    );
  };

  const activeVersion = getLatestTagVersionByStatus(
    tag.versions ?? [],
    "active",
  );

  return (
    <TagsPanel.Card
      title={tag.keyPath}
      description={tag.displayName}
      onContextMenu={handleContextMenu}
      isSelected={isTagSelected(tag)}
      slotProps={{
        content: {
          sx: { display: "flex", justifyContent: "end", alignItems: "center" },
        },
      }}
    >
      {activeVersion && (
        <TagVersionStatusTag
          version={activeVersion.version}
          status={activeVersion.status}
        />
      )}
    </TagsPanel.Card>
  );
};
