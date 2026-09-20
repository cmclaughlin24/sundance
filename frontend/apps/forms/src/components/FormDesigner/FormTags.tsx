import Box from "@mui/material/Box";
import { TagsPanel } from "./panels/TagsPanel";
import Typography from "@mui/material/Typography";
import { formTagsStyles as styles } from "./FormTags.styles";
import type { IElement } from "@/types/element";
import type { ITag } from "@/types/tag";
import {
  useFormDesignerSelect,
  useFormPagesSnapshot,
} from "@/store/formDesigner";
import { getFlattenedElements } from "@/utils/form";
import { useAsyncData } from "@/hooks/useAsyncData";
import { useTagsService } from "@/hooks/useHttpService";
import { TENANT_ID } from "@/constants/tenant";
import { useCallback } from "react";

export const FormTags: React.FC = function () {
  const pages = useFormPagesSnapshot();
  const { selected, select } = useFormDesignerSelect();
  const tagsService = useTagsService();

  const {
    data: tags,
    isLoading: _isLoadingTags,
    error: _tagError,
  } = useAsyncData(async (token) => {
    return await tagsService.getTags({ token, tenantId: TENANT_ID });
  }, []);

  const handleElement = (element: IElement) => {
    const isSelected = selected?.item.id === element.id;
    select(!isSelected ? { type: "element", item: element } : null);
  };

  const elementToTagsPanelCard = useCallback(
    (element: IElement) => {
      const isSelected = selected?.item.id === element.id;
      const isRequired = element.attributes.isRequired;

      return (
        <TagsPanel.Card
          title={element.name}
          description={element.key}
          isSelected={isSelected}
          onClick={() => handleElement(element)}
          onKeyDown={() => handleElement(element)}
          slotProps={{
            title: {
              sx: isRequired
                ? {
                    "::after": {
                      content: '"*"',
                      color: "#971E28",
                      marginLeft: 0.25,
                    },
                  }
                : {},
            },
          }}
        ></TagsPanel.Card>
      );
    },
    [selected],
  );

  const elements = getFlattenedElements(pages);

  return (
    <Box sx={styles.workspace}>
      <TagsPanel sx={styles.fields}>
        <TagsPanel.Header
          title="Form Fields (User Input)"
          slot={
            <Typography component="span" sx={styles.badge}>
              {elements?.length ?? 0} Field(s)
            </Typography>
          }
        />
        <TagsPanel.Content<IElement>
          data={elements}
          placeholder="Filter fields..."
          filterFn={filterElements}
          keyFn={(e) => e.id}
        >
          {elementToTagsPanelCard}
        </TagsPanel.Content>
      </TagsPanel>
      <Box sx={styles.contract}>
        <Typography sx={{ fontWeight: 600 }}>Contract Flow</Typography>
      </Box>
      <TagsPanel sx={styles.tags}>
        <TagsPanel.Header
          title="Canonical IGA Schema (Target)"
          slot={
            <Typography component="span" sx={styles.badge}>
              SCIM 2.0 / Saiyant EARS
            </Typography>
          }
        />
        <TagsPanel.Content<ITag>
          data={tags}
          placeholder="Filter canonical tags..."
          filterFn={filterTags}
          keyFn={(t) => t.id}
        >
          {(t) => (
            <TagsPanel.Card title={t.keyPath} description={t.displayName} />
          )}
        </TagsPanel.Content>
      </TagsPanel>
    </Box>
  );
};

function filterElements(searchTerm: string, element: IElement): boolean {
  const s = searchTerm.toLowerCase();

  return (
    element.name.toLowerCase().includes(s) ||
    element.key.toLowerCase().includes(s)
  );
}

function filterTags(searchTerm: string, tag: ITag): boolean {
  const s = searchTerm.toLowerCase();

  return (
    tag.displayName.toLowerCase().includes(s) ||
    tag.keyPath.toLowerCase().includes(s)
  );
}
