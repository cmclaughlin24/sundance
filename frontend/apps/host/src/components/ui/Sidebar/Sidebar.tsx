import * as React from "react";

export type SidebarProps = {
  className?: string;
};

const Sidebar: React.FC<SidebarProps> = function ({ className }) {
  return (
    <aside
      className={className}
    >
      Sidebar
    </aside>
  );
};

export default Sidebar;
