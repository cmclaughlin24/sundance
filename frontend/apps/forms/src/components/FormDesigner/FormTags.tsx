import Box from "@mui/material/Box";
import { TagsPanel } from "./panels/TagsPanel";
import Typography from "@mui/material/Typography";
import { formTagsStyles as styles } from "./FormTags.styles";
import type { IElement } from "@/types/element";
import type { ITag } from "@/types/tag";
import { useFormPagesSnapshot } from "@/store/formDesigner";
import { getFlattenedElements } from "@/utils/form";
import { useAsyncData } from "@/hooks/useAsyncData";
import { useTagsService } from "@/hooks/useHttpService";
import { TENANT_ID } from "@/constants/tenant";

export const FormTags: React.FC = function () {
  const pages = useFormPagesSnapshot();
  const tagsService = useTagsService();

  const {
    data: tags,
    isLoading: _isLoadingTags,
    error: _tagError,
  } = useAsyncData(async (token) => {
    return await tagsService.getTags({ token, tenantId: TENANT_ID });
  }, []);

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
        >
          {(element) => element.id}
        </TagsPanel.Content>
      </TagsPanel>
      <Box></Box>
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
        >
          {(tag) => tag.id}
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
