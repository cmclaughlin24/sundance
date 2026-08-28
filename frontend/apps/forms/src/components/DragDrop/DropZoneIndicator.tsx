import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";

export const DropZoneIndicator: React.FC<{ text: string }> = function ({
  text,
}) {
  return (
    <Box
      sx={{ display: "flex", justifyContent: "center", alignItems: "center" }}
    >
      <Typography>{text}</Typography>
    </Box>
  );
};
