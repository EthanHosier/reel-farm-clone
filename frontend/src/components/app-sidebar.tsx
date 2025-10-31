"use client";

import * as React from "react";
import { WandSparkles } from "lucide-react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar";
import { NavUser } from "@/components/nav-user";
import { useAuth } from "@/contexts/AuthContext";
import { useUser } from "@/features/dashboard/queries/useUser";
import { ROUTES } from "@/types/routes";

interface NavigationItem {
  title: string;
  url: string;
  icon?: React.ComponentType<{ className?: string }>;
  items?: {
    title: string;
    url: string;
    isActive?: boolean;
  }[];
}

const navigationData: NavigationItem[] = [
  {
    title: "Dashboard",
    url: "#",
    items: [
      {
        title: "Your Videos",
        url: ROUTES.userGeneratedVideos,
      },
      {
        title: "Generate AI UGC",
        url: ROUTES.generateAiAvatarVideo,
      },
    ],
  },
  {
    title: "Hooks",
    url: "#",
    items: [
      {
        title: "Generate Hooks",
        url: ROUTES.generateHooks,
      },
    ],
  },
];

interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {}

export function AppSidebar({ ...props }: AppSidebarProps) {
  const { user, signOut } = useAuth();

  const { data: userAccount } = useUser();

  const userData = {
    email: user?.email || "",
    avatar: user?.user_metadata.avatar_url || "",
    plan: userAccount?.plan || "free",
  };

  return (
    <Sidebar {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <a href="#">
                <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                  <WandSparkles className="size-4" />
                </div>
                <div className="flex flex-col gap-0.5 leading-none">
                  <span className="font-medium">Reel Farm</span>
                  <span className="text-xs">Dashboard</span>
                </div>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarMenu className="gap-2">
            {navigationData.map((item) => {
              const Icon = item.icon;
              return (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <a
                      href={item.url}
                      className="font-medium flex items-center gap-2"
                    >
                      {Icon && <Icon className="size-4" />}
                      <span>{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                  {item.items?.length ? (
                    <SidebarMenuSub className="ml-0 border-l-0 px-1.5">
                      {item.items.map((subItem) => (
                        <SidebarMenuSubItem key={subItem.title}>
                          <SidebarMenuSubButton
                            asChild
                            isActive={subItem.isActive}
                          >
                            <a href={subItem.url}>{subItem.title}</a>
                          </SidebarMenuSubButton>
                        </SidebarMenuSubItem>
                      ))}
                    </SidebarMenuSub>
                  ) : null}
                </SidebarMenuItem>
              );
            })}
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      {user && (
        <SidebarFooter>
          <NavUser user={userData} signOut={signOut} />
        </SidebarFooter>
      )}
    </Sidebar>
  );
}
