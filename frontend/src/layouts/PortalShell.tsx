import { useState } from "react";
import { Link, Outlet, useLocation } from "react-router";
import { useQueryClient } from "@tanstack/react-query";
import { Header, HeaderMenuButton, SideNav, SideNavItems, SideNavLink, SideNavDivider } from "@carbon/react";
import { AppHeaderBrand, AppHeaderActions } from "@/layouts/AppHeaderChrome";
import RouteErrorBoundary from "@/shared/ui/RouteErrorBoundary";
import { prefetchForPath } from "@/layouts/prefetchOnHover";
import type { NavGroup } from "@/layouts/nav/types";

interface Props {
  navGroups: NavGroup[];
  showSearch?: boolean;
  collapsible?: boolean;
}

function isActivePath(pathname: string, path: string, exact?: boolean) {
  return exact || path === "/" ? pathname === path : pathname.startsWith(path);
}

// Header, sidebar and content area shared by every portal.
export default function PortalShell({ navGroups, showSearch = false, collapsible = false }: Props) {
  const { pathname } = useLocation();
  const queryClient = useQueryClient();
  const [expanded, setExpanded] = useState(true);
  const collapsed = collapsible && !expanded;

  return (
    <>
      <Header aria-label="OpenSchool">
        {collapsible && (
          <HeaderMenuButton
            aria-label={expanded ? "Close menu" : "Open menu"}
            onClick={() => setExpanded((v) => !v)}
            isActive={expanded}
            aria-expanded={expanded}
          />
        )}
        <AppHeaderBrand />
        <AppHeaderActions showSearch={showSearch} />
      </Header>

      <div className="os-layout">
        <aside className={`os-layout__sidebar${collapsed ? " is-collapsed" : ""}`}>
          <SideNav aria-label="Side navigation" isFixedNav expanded={!collapsed} isPersistent>
            <SideNavItems>
              {navGroups.map((group, i) => (
                <div key={group.label}>
                  {i > 0 && <SideNavDivider />}
                  <p className="os-nav-group__label">{group.label}</p>
                  {group.items.map(({ path, label, Icon, exact }) => (
                    <SideNavLink
                      key={path}
                      as={Link}
                      to={path}
                      renderIcon={Icon}
                      isActive={isActivePath(pathname, path, exact)}
                      onMouseEnter={() => prefetchForPath(queryClient, path)}
                    >
                      {label}
                    </SideNavLink>
                  ))}
                </div>
              ))}
            </SideNavItems>
          </SideNav>
        </aside>

        <main className={`os-layout__content${collapsed ? " is-expanded" : ""}`}>
          <RouteErrorBoundary>
            <Outlet />
          </RouteErrorBoundary>
        </main>
      </div>
    </>
  );
}
