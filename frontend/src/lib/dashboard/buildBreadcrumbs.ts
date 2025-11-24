import { ROUTES } from "@/types/routes";

function formatBreadcrumbLabel(path: string): string {
  return path
    .split("-")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}

export function buildBreadcrumbs(
  pathname: string
): Array<{ label: string; href?: string }> {
  if (pathname === ROUTES.dashboard) {
    return [];
  }

  const segments = pathname.split("/").filter(Boolean);
  const breadcrumbs: Array<{ label: string; href?: string }> = [];

  segments.forEach((segment, _index) => {
    // const href = "/" + segments.slice(0, index + 1).join("/");
    const label = formatBreadcrumbLabel(segment);

    breadcrumbs.push({
      label,
      // href: index < segments.length - 1 ? href : undefined,
      href: "#",
    });
  });

  return breadcrumbs;
}
