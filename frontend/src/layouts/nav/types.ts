import type { ComponentType } from "react";

export interface NavItem {
  path: string;
  label: string;
  Icon: ComponentType<{ size?: number }>;
  // Match the path exactly instead of by prefix ("/" always matches exactly).
  exact?: boolean;
}

export interface NavGroup {
  label: string;
  items: NavItem[];
}
