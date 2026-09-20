import { useFormDesignerSelect } from "@/store/formDesigner";
import type { IElement } from "@/types/element";
import { TagsPanel } from "../panels/TagsPanel";
import { Tag } from "@/components/Tag";
import { findFormObjectPaletteItem } from "../panels/ToolboxPanel/palette";
import Typography from "@mui/material/Typography";
import type { Styles } from "@/types/styles";

const styles: Styles = {
  required: {
    "::after": {
      content: '"*"',
      color: "#971E28",
      marginLeft: 0.25,
    },
  },
  content: {
    display: "flex",
    justifyContent: "end",
    alignItems: "center",
    gap: 1,
  },
};

export interface ElementTagsPanelCardProps {
  element: IElement;
}

export const ElementTagsPanelCard: React.FC<ElementTagsPanelCardProps> =
  function ({ element }) {
    const { selected, isSelected, select } = useFormDesignerSelect(element.id);

    const handleElement = (element: IElement) => {
      const isSelected = selected?.item.id === element.id;
      select(!isSelected ? { type: "element", item: element } : null);
    };

    const paletteItem = findFormObjectPaletteItem(element.type);

    return (
      <TagsPanel.Card
        title={element.name}
        description={element.key}
        isSelected={isSelected}
        onClick={() => handleElement(element)}
        onKeyDown={() => handleElement(element)}
        slotProps={{
          title: {
            sx: element.attributes.isRequired ? styles.required : {},
          },
          content: { sx: styles.content },
        }}
      >
        {paletteItem && <Tag>{paletteItem.label}</Tag>}
        <Typography component="span" sx={{ fontSize: "0.75rem" }}>
          · {element.tags?.length ?? 0} Tag(s)
        </Typography>
      </TagsPanel.Card>
    );
  };
