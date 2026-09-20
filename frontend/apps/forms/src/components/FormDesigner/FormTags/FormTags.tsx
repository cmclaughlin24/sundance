import Box from "@mui/material/Box";
import { TagsPanel } from "../panels/TagsPanel";
import Typography from "@mui/material/Typography";
import { formTagsStyles as styles } from "./FormTags.styles";
import type { IElement } from "@/types/element";
import { useFormPagesSnapshot } from "@/store/formDesigner";
import { getFlattenedElements } from "@/utils/form";
import { useAsyncData } from "@/hooks/useAsyncData";
import { useTagsService } from "@/hooks/useHttpService";
import { TENANT_ID } from "@/constants/tenant";
import { ElementTagsPanelCard } from "./ElementTagsPanelCard";
import { groupTags, type TagGroup } from "@/utils/tag";
import { TagGroupList } from "./TagsList";

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
  const tagGroups = groupTags(tags || []);

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
          {(e) => <ElementTagsPanelCard element={e} />}
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
        <TagsPanel.Content<TagGroup>
          data={tagGroups}
          placeholder="Filter canonical tags..."
          filterFn={filterTags}
          keyFn={(group) => group.id}
        >
          {(group) => <TagGroupList group={group} />}
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

function filterTags(searchTerm: string, group: TagGroup): boolean {
  const s = searchTerm.toLowerCase();

  return group.items.some(
    (i) =>
      i.displayName.toLowerCase().includes(s) ||
      i.keyPath.toLowerCase().includes(s),
  );
}
