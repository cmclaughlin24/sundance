export type TabPanelProps<T> = React.PropsWithChildren<{
  value: T;
}>;

/**
 */
export function TabPanel<T>({ children }: TabPanelProps<T>) {
  return children;
}
